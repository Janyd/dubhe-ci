package core

import (
	"context"
	"dubhe-ci/common"
	"dubhe-ci/utils"
	"time"
)

type (
	Repository struct {
		Id           string    `json:"id" xorm:"notnull pk comment('ID') varchar(20)"`
		Name         string    `json:"name" xorm:"notnull default '' comment('库名称') varchar(50)"`
		Description  string    `json:"description" xorm:"notnull default '' comment('项目描述') varchar(100)"`
		Url          string    `json:"url" xorm:"notnull default '' comment('git代码库') varchar(300)"`
		CredentialId string    `json:"credentialId" xorm:"notnull default 0 comment('凭据ID') varchar(20)"`
		Active       bool      `json:"active" xorm:"notnull default true comment('是否激活状态') BOOL"`
		Timeout      int32     `json:"timeout" xorm:"notnull default 0 comment('执行超时时间(分钟)') int"`
		Config       string    `json:"config" xorm:"notnull default '' comment('配置文件路径') varchar(300)"`
		Created      time.Time `json:"created" xorm:"notnull default CURRENT_TIMESTAMP created comment('创建时间') TIMESTAMP"`
		Updated      time.Time `json:"updated" xorm:"notnull default CURRENT_TIMESTAMP updated comment('更新时间') TIMESTAMP"`
	}

	//存储库业务
	RepositoryStore interface {

		//分页获取存储库
		List(ctx context.Context, page *common.Page) (*common.Page, error)

		//查询名称是否重复
		CheckName(ctx context.Context, id string, name string) (bool, error)

		//创建存储库
		Create(ctx context.Context, repository *Repository) (*Repository, error)

		//更新存储库
		Update(ctx context.Context, repository *Repository) error

		//获取存储库
		Find(ctx context.Context, repoId string) (*Repository, error)

		//删除存储库
		Delete(ctx context.Context, repository *Repository) error
	}
)

func (r *Repository) CreatedAt() int64 {
	return r.Created.Unix()
}

func (r *Repository) UpdatedAt() int64 {
	return r.Updated.Unix()
}

func (r *Repository) BeforeInsert() {
	r.Id = utils.GetId()
}
