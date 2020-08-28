// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"dubhe-ci/config"
	"github.com/google/wire"
)

func InitializeApplication(config *config.Config) (application, error) {
	wire.Build(
		storeSet,
		serverSet,
		runnerSet,
		ginSet,
		gitSet,
		newApplication,
	)

	return application{}, nil
}
