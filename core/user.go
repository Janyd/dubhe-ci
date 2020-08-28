package core

import (
	"context"
	"dubhe-ci/common"
	"dubhe-ci/utils"
	"time"
)

type (
	User struct {
		Id       string    `json:"id" xorm:"notnull pk comment('ID') varchar(20)"`
		Name     string    `json:"name" xorm:"notnull default '' comment('昵称') varchar(50)"`
		Username string    `json:"username" xorm:"notnull default '' comment('用户名') varchar(50)"`
		Password string    `json:"password" xorm:"notnull default '' comment('密码') varchar(100)"`
		Salt     string    `json:"-" xorm:"notnull default '' comment('加密盐值') varchar(50)"`
		Active   bool      `json:"active" xorm:"notnull default true comment('是否激活状态') BOOL"`
		Created  time.Time `json:"created" xorm:"notnull default CURRENT_TIMESTAMP created comment('创建时间') TIMESTAMP"`
		Updated  time.Time `json:"updated" xorm:"notnull default CURRENT_TIMESTAMP updated comment('更新时间') TIMESTAMP"`
	}

	UserStore interface {
		//分页获取用户
		List(ctx context.Context, page *common.Page) (*common.Page, error)

		//查询名称是否重复
		CheckName(ctx context.Context, id string, username string) (bool, error)

		//创建用户
		Create(ctx context.Context, user *User) (*User, error)

		//更新用户
		Update(ctx context.Context, user *User) error

		//获取用户
		Find(ctx context.Context, userId string) (*User, error)

		//删除用户
		Delete(ctx context.Context, user *User) error
	}
)

func (u *User) CreatedAt() int64 {
	return u.Created.Unix()
}

func (u *User) UpdatedAt() int64 {
	return u.Updated.Unix()
}

func (u *User) BeforeInsert() {
	u.Id = utils.GetId()
}
