package main

import (
	"dubhe-ci/config"
	"dubhe-ci/store/build"
	"dubhe-ci/store/cred"
	"dubhe-ci/store/logs"
	"dubhe-ci/store/repo"
	"dubhe-ci/store/shared/db"
	"dubhe-ci/store/step"
	"github.com/google/wire"
)

var storeSet = wire.NewSet(
	provideDatabase,
	repo.New,
	repo.NewBranchStore,
	cred.New,
	build.New,
	step.New,
	logs.New,
)

//创建数据库连接
func provideDatabase(config *config.Config) (*db.DB, error) {
	return db.Connect(config.Database)
}
