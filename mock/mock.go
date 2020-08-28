package mock

//go:generate mockgen -package=mock -destination=mock_gen.go dubhe-ci/core RepositoryStore,BuildStore,BranchStore,CredentialStore,StepStore,Scheduler,ConfigService,ConvertService,Manager
