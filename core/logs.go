package core

import (
	"context"
	"dubhe-ci/utils"
	"io"
	"time"
)

type Line struct {
	Number    int    `json:"pos"`
	Message   string `json:"out"`
	Timestamp int64  `json:"time"`
}

type (
	Logs struct {
		Id      string    `json:"id" xorm:"notnull pk comment('ID') varchar(20)"`
		BuildId string    `json:"buildId" xorm:"notnull default 0 comment('构建ID') varchar(20)"`
		StepId  string    `json:"stepId" xorm:"notnull default 0 comment('步骤ID') varchar(20)"`
		Data    []byte    `json:"data" xorm:"comment('日志') text"`
		Created time.Time `json:"created" xorm:"notnull default CURRENT_TIMESTAMP created comment('创建时间') TIMESTAMP"`
	}

	LogStore interface {
		Find(ctx context.Context, stepId string) (io.ReadCloser, error)

		List(ctx context.Context, buildId string) (io.ReadCloser, error)

		Create(ctx context.Context, buildId string, stepId string, r io.Reader) error

		Update(ctx context.Context, stepId string, r io.Reader) error

		Delete(ctx context.Context, stepId string) error
	}

	LogStream interface {
		Create(ctx context.Context, stepId string) error

		Delete(ctx context.Context, stepId string) error

		Write(ctx context.Context, stepId string, line *Line) error

		Tail(ctx context.Context, stepId string, sub Subscriber) bool
	}

	Subscriber interface {
		Publish(line *Line)

		Close()
	}
)

func (s *Logs) BeforeInsert() {
	s.Id = utils.GetId()
}

func (s *Logs) CreatedAt() int64 {
	return s.Created.Unix()
}
