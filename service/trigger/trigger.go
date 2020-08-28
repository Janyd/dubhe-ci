package trigger

import (
	"context"
	"dubhe-ci/core"
	"dubhe-ci/yaml"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"time"
)

type service struct {
	config    core.ConfigService
	convert   core.ConvertService
	repos     core.RepositoryStore
	builds    core.BuildStore
	branchs   core.BranchStore
	scheduler core.Scheduler
}

func New(config core.ConfigService,
	convert core.ConvertService,
	repos core.RepositoryStore,
	builds core.BuildStore,
	branchs core.BranchStore,
	scheduler core.Scheduler) core.TriggerService {
	return &service{
		config:    config,
		convert:   convert,
		repos:     repos,
		builds:    builds,
		branchs:   branchs,
		scheduler: scheduler,
	}
}

func (t *service) Trigger(ctx context.Context, repo *core.Repository, hook *core.Hook) (*core.Build, error) {
	logger := logrus.WithFields(
		logrus.Fields{
			"repo":   repo.Name,
			"branch": hook.Branch,
			"ref":    hook.Ref,
			"event":  hook.Event,
			"commit": hook.After,
		},
	)

	defer func() {
		// taking the paranoid approach to recover from
		// a panic that should absolutely never happen.
		if r := recover(); r != nil {
			logger.Errorf("runner: unexpected panic: %s", r)
			debug.PrintStack()
		}
	}()

	branch, err := t.branchs.FindByName(ctx, repo.Id, hook.Branch)
	if err != nil {
		logger.WithError(err).Error("trigger: cannot find branch")
		return nil, err
	}

	if skipMessage(hook) {
		logger.Infoln("trigger: skipping hook. found skip directive")
		return nil, nil
	}

	//TODO 是否需要忽略pull或forks事件

	build := &core.Build{
		RepoId:      repo.Id,
		Branch:      hook.Branch,
		Trigger:     hook.Trigger,
		Status:      core.StatusPending,
		Event:       hook.Event,
		Title:       hook.Title,
		Message:     hook.Message,
		Before:      hook.Before,
		After:       hook.After,
		Ref:         hook.Ref,
		Author:      hook.Author,
		AuthorEmail: hook.AuthorEmail,
		Created:     time.Now(),
		Updated:     time.Now(),
	}

	raw, err := t.config.Find(ctx, &core.ConfigArgs{
		Repo:   repo,
		Branch: branch,
		Build:  build,
	})

	if err != nil {
		logger.WithError(err).Error("trigger: cannot find config")
		return nil, err
	}

	raw, err = t.convert.Convert(ctx, &core.ConvertArgs{
		Repo:   repo,
		Build:  build,
		Config: raw,
	})

	if err != nil {
		logger = logger.WithError(err)
		logger.Warnln("trigger: cannot convert config")
		return t.createBuildError(ctx, repo, branch, hook, err.Error())
	}

	resource, err := yaml.ParseString(raw.Data)
	if err != nil {
		logger = logger.WithError(err)
		logger.Warnln("trigger: cannot parse yaml")
		return t.createBuildError(ctx, repo, branch, hook, err.Error())
	}
	pipeline, ok := resource.(*yaml.Pipeline)
	if !ok {
		logger.Infoln("trigger: skipping build, no matching pipelines")
	}

	if skipBranch(pipeline, branch.Name) {
		logger.WithField("pipeline", pipeline.Name).
			Infoln("trigger: skipping pipeline, does not match branch")
	} else if skipEvent(pipeline, hook.Event) {
		logger.WithField("pipeline", pipeline.Name).
			Infoln("trigger: skipping pipeline, does not match event")
	} else if skipRef(pipeline, hook.Ref) {
		logger.WithField("pipeline", pipeline.Name).
			Infoln("trigger: skipping pipeline, does not match ref")
	} else if skipRepo(pipeline, repo.Name) {
		logger.WithField("pipeline", pipeline.Name).
			Infoln("trigger: skipping pipeline, does not match repo")
	}

	branch, err = t.branchs.Increment(ctx, branch)
	if err != nil {
		logger.WithError(err).Errorln("trigger: cannot increment build sequence")
		return nil, err
	}

	build.Number = branch.Counter

	//TODO save steps
	err = t.builds.Create(ctx, build)
	if err != nil {
		logger.WithError(err).Errorln("trigger: cannot create build")
		return nil, err
	}

	err = t.scheduler.Schedule(ctx, build)
	if err != nil {
		logger = logger.WithError(err)
		logger.Errorln("trigger: cannot enqueue build")
		return nil, err
	}

	return build, nil
}

func (t *service) createBuildError(ctx context.Context, repo *core.Repository, branch *core.Branch, hook *core.Hook, message string) (*core.Build, error) {
	branch, err := t.branchs.Increment(ctx, branch)
	if err != nil {
		return nil, err
	}

	build := &core.Build{
		RepoId:      repo.Id,
		Branch:      branch.Name,
		Trigger:     hook.Trigger,
		Number:      branch.Counter,
		Status:      core.StatusError,
		Error:       message,
		Event:       hook.Event,
		Title:       hook.Title,
		Message:     hook.Message,
		Before:      hook.Before,
		After:       hook.After,
		Ref:         hook.Ref,
		Author:      hook.Author,
		AuthorEmail: hook.AuthorEmail,
		Started:     time.Now(),
		Finished:    time.Now(),
	}

	err = t.builds.Create(ctx, build)

	branch.LastFailure = time.Now()

	err = t.branchs.Update(ctx, branch)
	if err != nil {
		return nil, err
	}

	return build, err
}
