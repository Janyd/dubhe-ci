package core

import (
	"context"
	"dubhe-ci/utils"
	"time"
)

type (
	Credential struct {
		Id             string    `json:"id" xorm:"notnull pk comment('ID') varchar(20)"`
		Name           string    `json:"name" xorm:"notnull default '' comment('凭据名称') varchar(50)"`
		CredentialType uint      `json:"credentialType" xorm:"not null default 1 comment('凭证类型 1-用户名密码 2-密钥 3-docker认证信息') int(2)"`
		Username       string    `json:"username" xorm:"not null default '' comment('用户名') varchar(50)" valid:"required~请输入用户名"`
		Password       string    `json:"password" xorm:"not null default '' comment('密码或Passphrase') varchar(200)" `
		PublicKey      string    `json:"publicKey" xorm:"comment('公钥') varchar(500)"`
		PrivateKey     string    `json:"privateKey" xorm:"comment('私钥') varchar(3000)"`
		Address        string    `json:"address" xorm:"comment('地址') varchar(200)"`
		Created        time.Time `json:"created" xorm:"notnull default CURRENT_TIMESTAMP created comment('创建时间') TIMESTAMP"`
		Updated        time.Time `json:"updated" xorm:"notnull default CURRENT_TIMESTAMP updated comment('更新时间') TIMESTAMP"`
	}

	CredentialStore interface {
		//凭据列表
		List(ctx context.Context) ([]*Credential, error)

		//查询所有Docker认证信息
		ListRegistryCred(ctx context.Context) ([]*Credential, error)

		//查询名称是否重复
		CheckName(ctx context.Context, id string, name string) (bool, error)

		//创建凭据
		Create(ctx context.Context, credential *Credential) (*Credential, error)

		//更新凭据
		Update(ctx context.Context, credential *Credential) error

		//更新凭据
		Find(ctx context.Context, credId string) (*Credential, error)

		//删除凭据
		Delete(ctx context.Context, credential *Credential) error
	}
)

func (c *Credential) BeforeInsert() {
	c.Id = utils.GetId()
}

func (c *Credential) CreatedAt() int64 {
	return c.Created.Unix()
}

func (c *Credential) UpdatedAt() int64 {
	return c.Updated.Unix()
}
