package build

import (
	"context"
	"dubhe-ci/common"
	"dubhe-ci/core"
	"dubhe-ci/errors"
	"dubhe-ci/logger"
	"dubhe-ci/store/shared/db"
)

func New(db *db.DB, stepStore core.StepStore) core.BuildStore {
	return &buildStore{
		db:        db,
		stepStore: stepStore,
	}
}

type buildStore struct {
	db        *db.DB
	stepStore core.StepStore
}

func (b *buildStore) Query() *db.Session {
	return db.WrapSession(b.db.Table(core.Build{}))
}

func (b *buildStore) Session() *db.Session {
	return db.WrapSession(b.db.NewSession())
}

func (b *buildStore) Find(ctx context.Context, id string) (*core.Build, error) {
	log := logger.WithAction("获取构建信息")
	var item core.Build
	has, err := b.Query().ID(id).Get(&item)
	if err != nil {
		log.WithError(err).Error("查询失败")
		return nil, err
	}

	if !has {
		return nil, errors.New(900000)
	}

	return &item, nil
}

func (b *buildStore) FindNumber(ctx context.Context, repoId string, number uint32) (*core.Build, error) {
	log := logger.WithAction("根据构建版本获取构建信息")

	var item core.Build
	has, err := b.Query().Eq("repo_id", repoId).Eq("number", number).SelectOne(&item)
	if err != nil {
		log.WithError(err).Error("查询失败")
		return nil, err
	}

	if !has {
		return nil, errors.New(900000)
	}

	return &item, nil
}

func (b *buildStore) Create(ctx context.Context, build *core.Build) error {
	log := logger.WithAction("创建构建信息")

	session := b.Session()
	_, err := session.InsertTran(build)
	if err != nil {
		log.WithError(err).Error("插入构建失败")
		return err
	}
	return nil
}

func (b *buildStore) Update(ctx context.Context, build *core.Build) error {
	log := logger.WithAction("更新构建信息")

	session := b.Session()
	_ = session.UseBool("active").ID(build.Id)
	_, err := session.UpdateTran(build)
	if err != nil {
		log.WithError(err).Error("更新失败")
		return err
	}

	return nil
}

func (b *buildStore) Delete(ctx context.Context, build *core.Build) error {
	log := logger.WithAction("删除构建信息")

	session := b.Session()

	_, err := session.DeleteTran(build)
	if err != nil {
		log.WithError(err).Error("删除失败")
		return err
	}
	return nil
}

func (b *buildStore) List(ctx context.Context, repoId string, branch string, page *common.Page) (*common.Page, error) {
	log := logger.WithAction("分页查询构建信息")
	list := make([]*core.Build, 0)
	page.Bind(&list, core.Build{})

	err := b.Query().Eq("repo_id", repoId).Eq("branch", branch).SelectPage(page)
	if err != nil {
		log.WithError(err).Error("构建信息查询错误")
		return nil, err
	}

	return page, nil
}

func (b *buildStore) ListIncomplete(ctx context.Context) ([]*core.Build, error) {
	log := logger.WithAction("查询未完成的构建信息")
	list := make([]*core.Build, 0)

	err := b.Query().IN("status", core.StatusPending, core.StatusRunning).SelectList(&list)
	if err != nil {
		log.WithError(err).Error("构建信息查询错误")
		return nil, err
	}

	return list, nil
}
