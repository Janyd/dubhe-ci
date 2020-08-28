package runner

import (
	"bytes"
	"context"
	"dubhe-ci/core"
	"dubhe-ci/errors"
	"dubhe-ci/logger"
	"dubhe-ci/server"
	"dubhe-ci/socket"
	"dubhe-ci/utils"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

func NewManager(
	repos core.RepositoryStore,
	branchs core.BranchStore,
	builds core.BuildStore,
	configs core.ConfigService,
	scheduler core.Scheduler,
	converter core.ConvertService,
	steps core.StepStore,
	logz core.LogStream,
	logs core.LogStore,
	server *server.SocketIOServer,
) core.Manager {
	return &manager{
		Jobs:           make(map[string]*core.Job),
		Repos:          repos,
		Branchs:        branchs,
		Builds:         builds,
		Config:         configs,
		Converter:      converter,
		Scheduler:      scheduler,
		Steps:          steps,
		Logz:           logz,
		Logs:           logs,
		SocketIOServer: server,
	}
}

type (
	manager struct {
		sync.Mutex

		Jobs      map[string]*core.Job
		Repos     core.RepositoryStore
		Branchs   core.BranchStore
		Builds    core.BuildStore
		Config    core.ConfigService
		Converter core.ConvertService
		Scheduler core.Scheduler
		Steps     core.StepStore
		Logz      core.LogStream
		Logs      core.LogStore

		SocketIOServer *server.SocketIOServer
	}
)

func (m *manager) Request(ctx context.Context) (*core.Build, error) {
	log := logger.WithRole("manager")

	log.Debugln("request queue item")

	build, err := m.Scheduler.Request(ctx)
	if err != nil && ctx.Err() != nil {
		log.Debugln("manager: context canceled")
		return nil, err
	}

	if err != nil {
		logrus.WithError(err).Warnln("manager: request queue item error")
		return nil, err
	}

	return build, err
}

func (m *manager) Accept(ctx context.Context, buildId string) (*core.Build, error) {

	log := logrus.WithFields(
		logrus.Fields{
			"build-id": buildId,
		},
	)

	log.Debugln("manager: accept build")

	build, err := m.Builds.Find(ctx, buildId)
	if err != nil {
		log = log.WithError(err)
		log.Warnln("manager: cannot find build")
		return nil, err
	}

	build.Status = core.StatusRunning

	err = m.Builds.Update(ctx, build)
	if err != nil {
		log = log.WithError(err)
		log.Debugln("manager: cannot update stage")
	}

	return build, err
}

func (m *manager) Details(ctx context.Context, buildId string) (*core.Context, error) {
	log := logrus.WithField("build-id", buildId)
	log.Debugln("manager: fetching build details")

	build, err := m.Builds.Find(ctx, buildId)
	if err != nil {
		log = log.WithError(err)
		log.Warnln("manager: cannot find build")
		return nil, err
	}

	repo, err := m.Repos.Find(ctx, build.RepoId)
	if err != nil {
		log = log.WithError(err)
		log.Warnln("manager: cannot find repository")
		return nil, err
	}

	branch, err := m.Branchs.FindByName(ctx, repo.Id, build.Branch)
	if err != nil {
		log = log.WithError(err)
		log.Warnln("manager: cannot find branch")
		return nil, err
	}

	log = log.WithFields(
		logrus.Fields{
			"build": build.Number,
			"repo":  repo.Name,
		},
	)
	config, err := m.Config.Find(ctx, &core.ConfigArgs{
		Repo:   repo,
		Branch: branch,
		Build:  build,
	})
	if err != nil {
		log = log.WithError(err)
		log.Warnln("manager: cannot find config")
		return nil, err
	}
	config, _ = m.Converter.Convert(ctx, &core.ConvertArgs{
		Repo:   repo,
		Build:  build,
		Config: config,
	})

	return &core.Context{
		Repo:   repo,
		Branch: branch,
		Build:  build,
		Config: &core.File{Data: []byte(config.Data)},
	}, nil
}

func (m *manager) Watch(ctx context.Context, buildId string) (bool, error) {
	ok, err := m.Scheduler.Cancelled(ctx, buildId)
	if err != nil {
		return ok, err
	}

	build, err := m.Builds.Find(ctx, buildId)
	if err != nil {
		logrus.WithError(err).WithField("build-id", buildId).Warnln("manager: cannot find build")
		return ok, err
	}

	return build.IsDone(), err
}

func (m *manager) BeforeAll(ctx context.Context, build *core.Build) error {
	log := logger.WithRole("manager").WithFields(logrus.Fields{
		"build.number": build.Number,
		"build.id":     build.Id,
		"build.branch": build.Branch,
		"repo.id":      build.RepoId,
	})

	if len(build.Error) > 1000 {
		build.Error = build.Error[:1000]
	}
	err := m.Builds.Update(ctx, build)
	if err != nil {
		log.WithError(err).WithField("build.status", build.Status).
			Warnln("cannot update the build")
		return err
	}
	err = m.Steps.Create(ctx, build.Steps)
	if err != nil {
		log.WithError(err).Warnln("cannot insert th step")
		return err
	}

	build.Steps, err = m.Steps.List(ctx, build.Id)
	if err != nil {
		log.WithError(err).Warnln("cannot list the steps")
		return err
	}

	err = m.Builds.Update(ctx, build)
	if err != nil {
		log.WithError(err).Warnln("cannot update th build")
		return err
	}

	return m.AddJob(ctx, build)
}

func (m *manager) AfterAll(ctx context.Context, build *core.Build) error {
	log := logger.WithRole("manager").WithFields(logrus.Fields{
		"build.number": build.Number,
		"build.id":     build.Id,
		"build.branch": build.Branch,
		"repo.id":      build.RepoId,
	})

	for _, step := range build.Steps {
		err := m.Steps.Update(ctx, step)
		if err != nil {
			log.WithError(err).
				WithField("step.name", step.Name).
				WithField("step.id", step.Id).
				WithField("build.status", build.Status).
				Warnln("cannot update the step")
			return err
		}
	}

	if len(build.Error) > 1000 {
		build.Error = build.Error[:1000]
	}

	branch, err := m.Branchs.FindByName(ctx, build.RepoId, build.Branch)
	if err != nil {
		log.WithError(err).
			WithField("build.status", build.Status).
			Warnln("cannot find the branch")
		return err
	}

	switch build.Status {
	case core.StatusFailing,
		core.StatusError,
		core.StatusKilled:
		branch.LastFailure = time.Now()
	case core.StatusPassing:
		branch.LastSuccess = time.Now()
	}

	err = m.Branchs.Update(ctx, branch)
	if err != nil {
		log.WithError(err).
			WithField("build.branch", build.Branch).
			WithField("build.number", build.Number).
			WithField("build.status", build.Status).
			Warnln("cannot update the branch")
		return err
	}

	build.Finished = time.Now()
	err = m.Builds.Update(ctx, build)
	if err != nil {
		log.WithError(err).Warnln("cannot update the build")
		return err
	}

	return m.FinishJob(ctx, build)
}

func (m *manager) Write(ctx context.Context, stepId string, line *core.Line) error {
	return m.Logz.Write(ctx, stepId, line)
}

func (m *manager) UploadBytes(ctx context.Context, stepId string, b []byte) error {
	buf := bytes.NewBuffer(b)
	step, err := m.Steps.Find(ctx, stepId)
	if err != nil {
		log := logrus.WithError(err)
		log = log.WithField("step-id", step)
		log.Warnln("manager: cannot find step")
		return err
	}
	err = m.Logs.Create(ctx, step.BuildId, stepId, buf)
	if err != nil {
		log := logrus.WithError(err)
		log = log.WithField("step-id", stepId)
		log = log.WithField("build-id", step.BuildId)
		log.Warnln("manager: cannot upload complete logs")
		return err
	}

	return err
}

func (m *manager) Before(ctx context.Context, step *core.Step) error {
	log := logrus.WithFields(
		logrus.Fields{
			"step.status": step.Status,
			"step.name":   step.Name,
			"step.id":     step.Id,
		},
	)
	log.Debugln("manager: updating step status")

	err := m.Logz.Create(ctx, step.Id)
	if err != nil {
		log.WithError(err).Warnln("manager: cannot create log stream")
		return err
	}

	updater := &updater{Steps: m.Steps}
	err = updater.do(ctx, step)
	if err != nil {
		return err
	}

	build, err := m.Builds.Find(ctx, step.BuildId)
	if err != nil {
		log.WithError(err).Warnln("manager: cannot find build")
		return err
	}

	return m.UpdateJob(ctx, build, step)
}

func (m *manager) After(ctx context.Context, step *core.Step) error {
	log := logrus.WithFields(
		logrus.Fields{
			"step.status": step.Status,
			"step.name":   step.Name,
			"step.id":     step.Id,
		},
	)
	log.Debugln("manager: updating step status")

	var errs error
	updater := &updater{Steps: m.Steps}

	if err := updater.do(ctx, step); err != nil {
		errs = multierror.Append(errs, err)
		log.WithError(err).Warnln("manager: cannot update step")
	}
	if err := m.Logz.Delete(ctx, step.Id); err != nil {
		log.WithError(err).Warnln("manager: cannot teardown log stream")
	}

	return errs
}

func (m *manager) AddJob(ctx context.Context, build *core.Build) error {
	m.Lock()
	defer m.Unlock()
	log := logger.WithAction("add job")

	repo, err := m.Repos.Find(ctx, build.RepoId)
	if err != nil {
		log.WithError(err).Warnln("manager: cannot find repo")
		return err
	}

	branch, err := m.Branchs.FindByName(ctx, build.RepoId, build.Branch)
	if err != nil {
		log.WithError(err).Warnln("manager: cannot find branch")
		return err
	}

	job := &core.Job{
		RepoId:      build.RepoId,
		Repo:        repo.Name,
		BranchId:    branch.Id,
		Branch:      branch.Name,
		BuildId:     build.Id,
		BuildNumber: build.Number,
	}

	m.Jobs[build.Id] = job

	go func() {
		m.SocketIOServer.Server.BroadcastToRoom("", socket.Broadcast, socket.LoadJob, utils.RSuccess(nil))
	}()

	return nil
}

func (m *manager) UpdateJob(ctx context.Context, build *core.Build, step *core.Step) error {
	m.Lock()
	defer m.Unlock()

	log := logger.WithAction("update job")

	job, ok := m.Jobs[build.Id]
	if !ok {
		log.Warnln("manager: cannot find job")
		return errors.New(900000)
	}

	job.StepId = step.Id
	job.StepName = step.Name

	return nil
}

func (m *manager) FinishJob(ctx context.Context, build *core.Build) error {
	m.Lock()
	defer m.Unlock()

	log := logger.WithAction("finish job")

	job, ok := m.Jobs[build.Id]
	if !ok {
		log.Warnln("manager: cannot find job")
		return errors.New(900000)
	}

	delete(m.Jobs, job.BuildId)

	return nil
}

func (m *manager) AllJob() []*core.Job {
	m.Lock()
	defer m.Unlock()

	jobs := make([]*core.Job, len(m.Jobs))
	for id := range m.Jobs {
		jobs = append(jobs, m.Jobs[id])
	}

	return jobs
}
