package runner

import (
	"context"
	"dubhe-ci/core"
	"dubhe-ci/logger"
	"dubhe-ci/runner/compiler"
	"dubhe-ci/yaml"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"
)

type Runner struct {
	sync.Mutex

	Engine   *DockerEngine
	Manager  core.Manager
	Compiler *compiler.Compiler

	BuildStore core.BuildStore
	Workspace  string
}

//启动运行者，n为能够多少个任务在运行
func (r *Runner) Start(ctx context.Context, n int) error {
	var g errgroup.Group
	for i := 0; i < n; i++ {
		g.Go(func() error {
			return r.start(ctx)
		})
	}
	return g.Wait()
}

func (r *Runner) start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			r.poll(ctx)
		}
	}
}

func (r *Runner) poll(ctx context.Context) error {

	log := logger.WithRole("runner")

	build, err := r.Manager.Request(ctx)
	if err != nil {
		log.WithError(err).Warnln("runner: cannot get queue item")
		return err
	}

	if build == nil || build.Id == "" {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = r.Manager.Accept(ctx, build.Id)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{
				"build-id": build.Id,
				"repo-id":  build.RepoId,
				"branch":   build.Branch,
			},
		).Warnln("cannot ack build")
		return err
	}

	go func() {
		log.Debugln("runner: watch for cancel signal")
		done, _ := r.Manager.Watch(ctx, build.Id)
		if done {
			cancel()
			log.Debugln("runner: received cancel signal")
		} else {
			log.Debugln("runner: done listening for cancel signals")
		}
	}()

	return r.Run(ctx, build.Id)
}

func (r *Runner) Run(ctx context.Context, buildId string) error {
	log := logrus.WithFields(
		logrus.Fields{
			"build-id": buildId,
		},
	)

	defer func() {
		// taking the paranoid approach to recover from
		// a panic that should absolutely never happen.
		if r := recover(); r != nil {
			log.Errorf("runner: unexpected panic: %s", r)
			debug.PrintStack()
		}
	}()

	m, err := r.Manager.Details(ctx, buildId)
	if err != nil {
		log.WithError(err).Warnln("runner: cannot get build details")
		return err
	}

	log = log.WithFields(
		logrus.Fields{
			"repo":  m.Repo.Name,
			"build": m.Build.Number,
		},
	)

	if m.Build.Status == core.StatusKilled || m.Build.Status == core.StatusSkipped {
		log = log.WithError(err)
		log.Infoln("runner: cannot run a canceled build")
		return nil
	}

	resource, err := yaml.ParseByte(m.Config.Data)
	pipeline, ok := resource.(*yaml.Pipeline)
	if !ok {
		return r.handleError(ctx, m.Build, errors.New("cannot find named pipeline"))
	}

	log = log.WithField("pipeline", pipeline.Name)

	compilerArgs := compiler.Args{
		Pipeline:  pipeline,
		Repo:      m.Repo,
		Build:     m.Build,
		Workspace: filepath.Join(r.Workspace, m.Repo.Name, m.Build.Branch),
	}

	spec := r.Compiler.Compile(ctx, compilerArgs)

	build := m.Build
	steps := map[string]*core.Step{}

	for _, src := range spec.Steps {
		if src.RunPolicy == compiler.RunNever {
			continue
		}
		dst := &core.Step{
			RepoId:  m.Repo.Id,
			BuildId: build.Id,
			Number:  uint32(len(build.Steps) + 1),
			Name:    src.Name,
			Status:  core.StatusPending,
		}

		build.Steps = append(build.Steps, dst)
		steps[src.Name] = dst
	}

	build.Started = time.Now()
	build.Status = core.StatusRunning

	err = r.Manager.BeforeAll(ctx, build)
	if err != nil {
		log.WithError(err).Warnln("cannot initialize pipeline")
		return r.handleError(ctx, build, err)
	}
	log.Debug("updated build to running")

	state := &core.State{
		Build: build,
		Repo:  m.Repo,
	}

	hook := &Hook{
		BeforeEach: func(state *State) error {
			r.Lock()
			step, ok := steps[state.step.Name]
			if ok {
				step.Status = core.StatusRunning
				step.Started = time.Now()
				state.step.Envs["DUBHE_STEP_NAME"] = step.Name
				state.step.Envs["DUBHE_STEP_NUMBER"] = fmt.Sprint(step.Number)
			}

			stepClone := new(core.Step)
			*stepClone = *step
			r.Unlock()
			return r.Manager.Before(ctx, stepClone)
		},
		AfterEach: func(state *State) error {
			r.Lock()
			step, ok := steps[state.step.Name]
			if ok {
				step.Status = core.StatusPassing
				step.Stopped = time.Now()
				step.ExitCode = state.state.ExitCode
				if state.state.ExitCode != 0 && state.state.ExitCode != 78 {
					step.Status = core.StatusFailing
				}
			}
			stepClone := new(core.Step)
			*stepClone = *step
			r.Unlock()
			return r.Manager.After(ctx, stepClone)
		},
		GotLine: func(state *State, line *core.Line) error {
			r.Lock()
			step, ok := steps[state.step.Name]
			r.Unlock()
			if !ok {
				return nil
			}

			return r.Manager.Write(ctx, step.Id, line)
		},
		GotLogs: func(state *State, lines []*core.Line) error {
			r.Lock()
			step, ok := steps[state.step.Name]
			r.Unlock()
			if !ok {
				// TODO log error
				return nil
			}
			raw, _ := json.Marshal(lines)
			return r.Manager.UploadBytes(ctx, step.Id, raw)
		},
	}

	execer := New(r.Engine, spec, hook, state)

	err = execer.Exec(ctx)
	if err != nil {
		log.WithError(err).Debug("build failed")
		return r.handleError(ctx, build, err)
	}

	return r.Manager.AfterAll(ctx, build)
}

func (r *Runner) handleError(ctx context.Context, build *core.Build, err error) error {
	switch build.Status {
	case core.StatusPending,
		core.StatusRunning:
	default:
	}

	for _, step := range build.Steps {
		if step.Status == core.StatusPending {
			step.Status = core.StatusSkipped
		}
		if step.Status == core.StatusRunning {
			step.Status = core.StatusPassing
			step.Stopped = time.Now()
		}
	}

	build.Status = core.StatusError
	build.Error = err.Error()

	switch v := err.(type) {
	case *ExitError:
		build.Error = ""
		build.Status = core.StatusFailing
		build.ExitCode = v.Code
	case *OomError:
		build.Error = "OOM kill signaled by host operating system"
	}
	return r.Manager.AfterAll(ctx, build)
}
