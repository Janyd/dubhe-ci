package core

import (
	"context"
	"dubhe-ci/utils"
	"time"
)

type (
	Branch struct {
		Id           string    `json:"id" xorm:"notnull pk comment('ID') varchar(20)"`
		Name         string    `json:"name" xorm:"notnull unique(branch_unq) default '' comment('分支名称') varchar(50)"`
		RepoId       string    `json:"repoId"  xorm:"notnull unique(branch_unq) default 0 comment('存储库ID') varchar(20)"`
		Active       bool      `json:"active" xorm:"notnull default true comment('是否激活状态') BOOL"`
		Counter      uint32    `json:"counter" xorm:"notnull default 0 comment('计数器') int"`
		LastDuration uint32    `json:"lastDuration" xorm:"notnull default 0 comment('上次构建耗时') int"`
		LastSuccess  time.Time `json:"lastSuccess" xorm:"comment('上次成功构建时间') TIMESTAMP"`
		LastFailure  time.Time `json:"lastFailure" xorm:"comment('上次失败构建时间') TIMESTAMP"`
		Created      time.Time `json:"created" xorm:"notnull default CURRENT_TIMESTAMP created comment('创建时间') TIMESTAMP"`
		Updated      time.Time `json:"updated" xorm:"notnull default CURRENT_TIMESTAMP updated comment('更新时间') TIMESTAMP"`
	}

	BranchStore interface {
		//获取所有分支
		List(ctx context.Context, repoId string) ([]*Branch, error)

		//创建分支
		Create(ctx context.Context, repoId string, name string) (*Branch, error)

		//激活或停用
		InActivate(ctx context.Context, repoId string, name string) error

		FindByName(ctx context.Context, repoId string, name string) (*Branch, error)

		//获取分支
		Find(ctx context.Context, branchId string) (*Branch, error)

		//更新
		Update(ctx context.Context, branch *Branch) error

		//删除分支信息
		Delete(ctx context.Context, branch *Branch) error

		//增加分支构建次数
		Increment(ctx context.Context, branch *Branch) (*Branch, error)
	}
)

func (b *Branch) BeforeInsert() {
	b.Id = utils.GetId()
}

func (b *Branch) CreatedAt() int64 {
	return b.Created.Unix()
}

func (b *Branch) UpdatedAt() int64 {
	return b.Updated.Unix()
}

func (b *Branch) LastSuccessAt() int64 {
	return b.LastSuccess.Unix()
}

func (b *Branch) LastFailureAt() int64 {
	return b.LastFailure.Unix()
}
