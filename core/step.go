package core

import (
	"context"
	"dubhe-ci/utils"
	"time"
)

type (
	Step struct {
		Id        string    `json:"id" xorm:"notnull pk comment('ID') varchar(20)"`
		RepoId    string    `json:"repoId" xorm:"notnull default 0 comment('存储库ID') varchar(20)"`
		BuildId   string    `json:"buildId" xorm:"notnull default 0 comment('构建ID') varchar(20)"`
		Number    uint32    `json:"number" xorm:"notnull default 0 comment('阶段数字') int"`
		Name      string    `json:"name" xorm:"notnull default '' comment('阶段名称') varchar(50)"`
		Status    string    `json:"status" xorm:"notnull default 0 comment('状态') varchar(20)"`
		ErrIgnore bool      `json:"errIgnore" xorm:"notnull default false comment('是否忽略错误')"`
		Error     string    `json:"error" xorm:"notnull default 0 comment('错误信息') text"`
		ExitCode  int       `json:"exitCode" xorm:"notnull default 0 comment('退出码') int"`
		Started   time.Time `json:"started" xorm:"notnull default CURRENT_TIMESTAMP comment('开始时间') TIMESTAMP"`
		Stopped   time.Time `json:"stopped" xorm:"notnull default CURRENT_TIMESTAMP comment('结束时间') TIMESTAMP"`
		Created   time.Time `json:"created" xorm:"notnull default CURRENT_TIMESTAMP created comment('创建时间') TIMESTAMP"`
		Updated   time.Time `json:"updated" xorm:"notnull default CURRENT_TIMESTAMP updated comment('更新时间') TIMESTAMP"`
	}

	StepStore interface {
		//查询阶段信息
		Find(ctx context.Context, id string) (*Step, error)

		//通过步骤版本查询
		FindNumber(ctx context.Context, stageId string, number uint32) (*Step, error)

		//批量创建
		Create(ctx context.Context, step []*Step) error

		//更新步骤信息
		Update(ctx context.Context, step *Step) error

		//步骤列表
		List(ctx context.Context, buildId string) ([]*Step, error)
	}
)

func (s *Step) BeforeInsert() {
	s.Id = utils.GetId()
}

func (s *Step) CreatedAt() int64 {
	return s.Created.Unix()
}

func (s *Step) UpdatedAt() int64 {
	return s.Updated.Unix()
}

func (s *Step) StartedAt() int64 {
	return s.Started.Unix()
}

func (s *Step) StoppedAt() int64 {
	return s.Stopped.Unix()
}
