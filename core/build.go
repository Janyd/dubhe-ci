package core

import (
	"context"
	"dubhe-ci/common"
	"dubhe-ci/utils"
	"time"
)

type (
	Build struct {
		Id          string `json:"id" xorm:"notnull pk comment('ID') varchar(20)"`
		RepoId      string `json:"repoId" xorm:"notnull default 0 comment('存储库ID') varchar(20)"`
		Branch      string `json:"branch" xorm:"notnull default '' comment('分支') varchar(50)"`
		Trigger     string `json:"trigger" xorm:"notnull default '' comment('触发器') varchar(50)"`
		Number      uint32 `json:"number" xorm:"notnull default 0 comment('构建数字') int"`
		Status      string `json:"status" xorm:"notnull default 0 comment('状态') varchar(20)"`
		ExitCode    int    `json:"exitCode" xorm:"notnull default 0 comment('退出码') int"`
		Error       string `json:"error" xorm:"notnull default 0 comment('错误信息') text"`
		Event       string `json:"event" xorm:"notnull default 0 comment('事件') text"`
		Title       string `json:"title" xorm:"notnull default 0 comment('构建标题') varchar(100)"`
		Message     string `json:"message" xorm:"notnull default 0 comment('构建信息') text"`
		Before      string `json:"before" xorm:"notnull default 0 comment('build before ref hash') varchar(100)"`
		After       string `json:"after" xorm:"notnull default 0 comment('build after ref hash') varchar(100)"`
		Ref         string `json:"ref" xorm:"notnull default 0 comment('git ref hash') varchar(100)"`
		Author      string `json:"author" xorm:"notnull default 0 comment('git ref hash') varchar(100)"`
		AuthorEmail string `json:"author_email" xorm:"notnull default 0 comment('git ref hash') varchar(100)"`
		Cron        string `json:"cron" xorm:"notnull default '' comment('cron') varchar(20)"`

		Started  time.Time `json:"started" xorm:" comment('开始时间') TIMESTAMP"`
		Finished time.Time `json:"finished" xorm:" comment('结束时间') TIMESTAMP"`
		Created  time.Time `json:"created" xorm:"notnull default CURRENT_TIMESTAMP created comment('创建时间') TIMESTAMP"`
		Updated  time.Time `json:"updated" xorm:"notnull default CURRENT_TIMESTAMP updated comment('更新时间') TIMESTAMP"`
		Changes  []string  `json:"changes" xorm:"notnull default '[]' comment('更变文件') json"`
		Steps    []*Step   `json:"steps" xorm:"-"`
	}

	BuildStore interface {

		//查询构建信息
		Find(ctx context.Context, id string) (*Build, error)

		//通过构建版本查询
		FindNumber(ctx context.Context, repoId string, number uint32) (*Build, error)

		//创建构建信息
		Create(ctx context.Context, build *Build) error

		//更新构建信息
		Update(ctx context.Context, build *Build) error

		//删除构建信息
		Delete(ctx context.Context, repository *Build) error

		//构建列表
		List(ctx context.Context, repoId string, branchId string, page *common.Page) (*common.Page, error)

		//查询未完成构建
		ListIncomplete(ctx context.Context) ([]*Build, error)
	}
)

func (b *Build) BeforeInsert() {
	b.Id = utils.GetId()
}

func (b *Build) CreatedAt() int64 {
	return b.Created.Unix()
}

func (b *Build) UpdatedAt() int64 {
	return b.Updated.Unix()
}

func (b *Build) StartedAt() int64 {
	return b.Started.Unix()
}

func (b *Build) FinishedAt() int64 {
	return b.Finished.Unix()
}

func (b *Build) IsDone() bool {
	switch b.Status {
	case StatusWaiting,
		StatusPending,
		StatusRunning,
		StatusBlocked:
		return false
	default:
		return true
	}
}

func (b *Build) IsFailed() bool {
	switch b.Status {
	case StatusFailing,
		StatusKilled,
		StatusError:
		return true
	default:
		return false
	}
}
