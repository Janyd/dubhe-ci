package core

import (
	"context"
)

type (
	Context struct {
		Repo   *Repository `json:"repository"`
		Branch *Branch     `json:"branch"`
		Build  *Build      `json:"build"`
		Config *File       `json:"config"`
	}

	Manager interface {
		//获取一个构建任务
		Request(ctx context.Context) (*Build, error)

		//接受此构建任务
		Accept(ctx context.Context, buildId string) (*Build, error)

		//获取构建详细信息
		Details(ctx context.Context, buildId string) (*Context, error)

		//关注构建信息动向
		Watch(ctx context.Context, buildId string) (bool, error)

		//执行步骤前操作
		Before(ctx context.Context, step *Step) error

		//执行步骤后操作
		After(ctx context.Context, step *Step) error

		//构建前操作
		BeforeAll(ctx context.Context, build *Build) error

		//构建后操作
		AfterAll(ctx context.Context, build *Build) error

		//写入日志
		Write(ctx context.Context, stepId string, line *Line) error

		//上传完整日志
		UploadBytes(ctx context.Context, stepId string, b []byte) error

		//添加任务
		AddJob(ctx context.Context, build *Build) error

		//更新任务
		UpdateJob(ctx context.Context, build *Build, step *Step) error

		//完成任务
		FinishJob(ctx context.Context, build *Build) error

		//获取所有任务
		AllJob() []*Job
	}

	Job struct {
		RepoId      string `json:"repoId"`
		Repo        string `json:"repo"`
		BranchId    string `json:"branchId"`
		Branch      string `json:"branch"`
		BuildId     string `json:"buildId"`
		BuildNumber uint32 `json:"buildNumber"`
		StepId      string `json:"stepId"`
		StepName    string `json:"stepName"`
	}
)
