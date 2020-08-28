package main

import (
	"dubhe-ci/config"
	"dubhe-ci/core"
	pluginConfig "dubhe-ci/plugins/config"
	"dubhe-ci/plugins/converter"
	"dubhe-ci/runner"
	"dubhe-ci/runner/compiler"
	"dubhe-ci/runner/streamer/livelog"
	"dubhe-ci/service/content/workspace"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var runnerSet = wire.NewSet(
	livelog.New,
	provideFileService,
	provideConverter,
	pluginConfig.Repository,
	compiler.New,
	runner.NewManager,
	provideRunner,
)

func provideRunner(
	manager core.Manager,
	compiler *compiler.Compiler,
	buildStore core.BuildStore,
	config *config.Config,
) *runner.Runner {
	engine, err := runner.NewEnvEngine(false)
	if err != nil {
		logrus.WithError(err).Fatalln("cannot load the docker engine")
	}
	return &runner.Runner{
		Engine:     engine,
		Manager:    manager,
		Compiler:   compiler,
		BuildStore: buildStore,
		Workspace:  config.Workspace,
	}
}

func provideFileService(config *config.Config) core.FileService {
	return workspace.New(config.Workspace)
}

func provideConverter() core.ConvertService {
	return converter.Combined()
}
