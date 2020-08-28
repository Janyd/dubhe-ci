package runner

import (
	"context"
	"dubhe-ci/core"
	"dubhe-ci/logger"
	"dubhe-ci/runner/compiler"
	"github.com/hashicorp/go-multierror"
)

type Execer struct {
	engine *DockerEngine
	hook   *Hook
	spec   *compiler.Spec
	state  *core.State
}

func New(engine *DockerEngine, spec *compiler.Spec, hook *Hook, state *core.State) *Execer {
	return &Execer{
		engine: engine,
		hook:   hook,
		spec:   spec,
		state:  state,
	}
}

func (e *Execer) Exec(ctx context.Context) error {
	defer e.engine.Destroy(ctx, e.spec)

	if err := e.engine.Setup(ctx, e.spec); err != nil {
		e.state.FailAll(err)
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, step := range e.spec.Steps {
		if err := e.exec(ctx, step); err != nil {
			return err
		}
	}

	e.state.FinishAll()
	return nil

}

func (e *Execer) exec(ctx context.Context, step *compiler.Step) error {
	var result error

	select {
	case <-ctx.Done():
		e.state.Cancel()
		return nil
	default:
	}

	log := logger.WithRole("execer").WithField("step.name", step.Name)
	log.Infoln("starting exec")

	switch {
	case e.state.Skipped():
		return nil
	case e.state.Cancelled():
		return nil
	case step.RunPolicy == compiler.RunNever:
		return nil
	case step.RunPolicy == compiler.RunAlways:
		break
	case step.RunPolicy == compiler.RunOnFailure && e.state.Failed() == false:
		e.state.Skip(step.Name)
		return nil
	case step.RunPolicy == compiler.RunOnSuccess && e.state.Failed():
		e.state.Skip(step.Name)
		return nil
	}

	e.state.Start(step.Name)

	//wc := e.streamer.Stream(ctx, state, step.Name)

	state := snapshot(e, step, nil)
	//output := outFn(state, e.hook)

	if step.Detach {
		go func() {
			_, _ = e.engine.Run(ctx, state)
		}()
		return nil
	}

	exited, err := e.engine.Run(ctx, state)
	if err != nil {
		result = multierror.Append(result, err)
	}

	if exited != nil {
		e.state.Finish(step.Name, exited.ExitCode)

		if exited.ExitCode == 78 {
			e.state.SkipAll()
		}

		return result
	}

	switch err {
	case context.Canceled, context.DeadlineExceeded:
		e.state.Cancel()
		return nil
	}

	return err
}
