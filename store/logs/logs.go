package logs

import (
	"bytes"
	"context"
	"dubhe-ci/core"
	"dubhe-ci/errors"
	"dubhe-ci/logger"
	"dubhe-ci/store/shared/db"
	"io"
	"io/ioutil"
)

func New(db *db.DB) core.LogStore {
	return &logStore{db: db}
}

type logStore struct {
	db *db.DB
}

func (l *logStore) Query() *db.Session {
	return db.WrapSession(l.db.Table(core.Logs{}))
}

func (l *logStore) Session() *db.Session {
	return db.WrapSession(l.db.NewSession())
}

func (l *logStore) Find(ctx context.Context, stepId string) (io.ReadCloser, error) {
	log := logger.WithAction("查询日志")
	var item core.Logs
	has, err := l.Query().Eq("step_id", stepId).SelectOne(&item)
	if err != nil {
		log.WithError(err).WithField("step.id", stepId).Error("查询失败")
		return nil, err
	}

	if !has {
		return nil, errors.New(900000)
	}

	return ioutil.NopCloser(
		bytes.NewBuffer(item.Data),
	), nil
}

func (l *logStore) List(ctx context.Context, buildId string) (io.ReadCloser, error) {
	log := logger.WithAction("查询构建日志")
	list := make([]*core.Logs, 0)
	err := l.Query().Eq("build_id", buildId).SelectList(&list)
	if err != nil {
		log.WithError(err).WithField("build.id", buildId).Error("构建日志查询失败")
		return nil, err
	}

	var buffer bytes.Buffer

	for _, logs := range list {
		buffer.Write(logs.Data)
	}

	return ioutil.NopCloser(
		bytes.NewBuffer(buffer.Bytes()),
	), nil
}

func (l *logStore) Create(ctx context.Context, buildId string, stepId string, r io.Reader) error {
	log := logger.WithAction("添加日志")
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	item := &core.Logs{
		BuildId: buildId,
		StepId:  stepId,
		Data:    data,
	}

	_, err = l.Session().InsertTran(item)
	if err != nil {
		log.WithError(err).
			WithField("build.id", buildId).
			WithField("step.id", stepId).
			Error("添加日志失败")
	}
	return err
}

func (l *logStore) Update(ctx context.Context, stepId string, r io.Reader) error {
	log := logger.WithAction("更新日志")
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	var item core.Logs

	has, err := l.Query().Eq("step_id", stepId).SelectOne(&item)
	if err != nil {
		log.WithError(err).WithField("step.id", stepId).Error("查询失败")
		return err
	}

	if !has {
		return errors.New(900000)
	}
	item.Data = data

	_, err = l.Session().UpdateTran(item)
	if err != nil {
		log.WithError(err).WithField("step.id", stepId).Error("更新日志失败")
	}
	return err
}

func (l *logStore) Delete(ctx context.Context, stepId string) error {
	log := logger.WithAction("删除日志")
	var item core.Logs

	has, err := l.Query().Eq("step_id", stepId).SelectOne(&item)
	if err != nil {
		log.WithError(err).WithField("step.id", stepId).Error("查询失败")
		return err
	}

	if !has {
		return errors.New(900000)
	}
	_, err = l.Session().DeleteTran(item)
	if err != nil {
		log.WithError(err).WithField("step.id", stepId).Error("删除日志失败")
	}
	return err
}
