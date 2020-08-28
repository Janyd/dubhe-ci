package main

import (
	"dubhe-ci/config"
	"dubhe-ci/scheduler/queue"
	"dubhe-ci/server"
	"dubhe-ci/service/repos"
	"dubhe-ci/service/rpc"
	"dubhe-ci/service/trigger"
	"github.com/google/wire"
)

var serverSet = wire.NewSet(
	repos.NewRepoService,
	repos.NewCredentialService,
	repos.NewBuildService,
	repos.NewBranchService,
	provideServer,
	provideRegisters,
	queue.New,
	trigger.New,
)

func provideServer(registers *rpc.Registers, config *config.Config) *server.GrpcServer {
	return &server.GrpcServer{
		Addr:      config.Server.Address,
		Network:   config.Server.Network,
		Registers: registers,
	}
}

func provideRegisters(
	repo *repos.RepositoryService,
	cred *repos.CredentialService,
	branch *repos.BranchService,
) *rpc.Registers {
	registers := rpc.New()

	registers.Add(repo)
	registers.Add(cred)
	registers.Add(branch)

	return registers
}
