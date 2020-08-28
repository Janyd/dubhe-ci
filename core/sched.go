package core

import (
	"context"
)

//调度器
type Scheduler interface {

	//获取任务
	Request(ctx context.Context) (*Build, error)

	//调度构建
	Schedule(ctx context.Context, build *Build) error

	//取消构建
	Cancel(ctx context.Context, buildId string) error

	//是否已取消构建
	Cancelled(ctx context.Context, buildId string) (bool, error)

	//暂停调度构建
	Pause(ctx context.Context) error

	//恢复构建
	Resume(ctx context.Context) error
}
