package step

import (
	"context"
	"dubhe-ci/core"
	"dubhe-ci/errors"
	"dubhe-ci/logger"
	"dubhe-ci/store/shared/db"
)

type stepStore struct {
	db *db.DB
}

func New(db *db.DB) core.StepStore {
	return &stepStore{db: db}
}

func (s *stepStore) Query() *db.Session {
	return db.WrapSession(s.db.Table(core.Step{}))
}

func (s *stepStore) Session() *db.Session {
	return db.WrapSession(s.db.NewSession())
}

func (s *stepStore) Find(ctx context.Context, id string) (*core.Step, error) {
	log := logger.WithAction("获取存储库")

	var item core.Step
	has, err := s.Query().ID(id).Get(&item)
	if err != nil {
		log.WithError(err).Error("查询失败")
		return nil, err
	}

	if !has {
		return nil, errors.New(900000)
	}

	return &item, nil
}

func (s *stepStore) FindNumber(ctx context.Context, stageId string, number uint32) (*core.Step, error) {
	log := logger.WithAction("根据阶段版本获取步骤信息")

	var item core.Step
	has, err := s.Query().Eq("stage_id", stageId).Eq("number", number).SelectOne(&item)
	if err != nil {
		log.WithError(err).Error("查询失败")
		return nil, err
	}

	if !has {
		return nil, errors.New(900000)
	}

	return &item, nil
}

func (s *stepStore) Create(ctx context.Context, steps []*core.Step) error {
	log := logger.WithAction("创建构建阶段步骤信息")

	if len(steps) == 0 {
		return nil
	}
	session := s.Session()
	for _, step := range steps {
		if len(step.Error) > 1000 {
			step.Error = step.Error[:1000]
		}
	}
	_, err := session.InsertTran(steps)
	if err != nil {
		log.WithError(err).Error("插入构建阶段步骤失败")
		return err
	}

	return nil
}

func (s *stepStore) Update(ctx context.Context, step *core.Step) error {
	log := logger.WithAction("更新构建阶段信息")

	if len(step.Error) > 1000 {
		step.Error = step.Error[:1000]
	}

	session := s.Session()
	_ = session.UseBool("active").ID(step.Id)
	_, err := session.UpdateTran(step)
	if err != nil {
		log.WithError(err).Error("更新失败")
		return err
	}

	return nil
}

func (s *stepStore) List(ctx context.Context, buildId string) ([]*core.Step, error) {
	log := logger.WithAction("查询构建阶段步骤信息")

	stages := make([]*core.Step, 0)
	err := s.Query().Eq("build_id", buildId).SelectList(&stages)
	if err != nil {
		log.WithError(err).Error("构建阶段步骤信息查询错误")
		return nil, err
	}

	return stages, nil
}
