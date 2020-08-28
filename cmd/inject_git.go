package main

import (
	"dubhe-ci/config"
	"dubhe-ci/scm"
	"github.com/google/wire"
)

var gitSet = wire.NewSet(
	provideGitService,
)

func provideGitService(config *config.Config) scm.GitService {
	return scm.New(config.Workspace)
}
