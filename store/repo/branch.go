package repo

import (
	"context"
	"dubhe-ci/core"
	"dubhe-ci/errors"
	"dubhe-ci/logger"
	"dubhe-ci/store/shared/db"
)

type branchStore struct {
	db *db.DB
}

func NewBranchStore(db *db.DB) core.BranchStore {
	return &branchStore{db: db}
}

func (b *branchStore) Session() *db.Session {
	return db.WrapSession(b.db.NewSession())
}

func (b *branchStore) Query() *db.Session {
	return db.WrapSession(b.db.Table(core.Branch{}))
}

func (b *branchStore) List(ctx context.Context, repoId string) ([]*core.Branch, error) {
	log := logger.WithAction("查询分支")

	branches := make([]*core.Branch, 0)
	err := b.Query().Eq("repo_id", repoId).SelectList(&branches)
	if err != nil {
		log.WithError(err).Error("查询分支错误")
		return nil, err
	}

	return branches, nil
}

func (b *branchStore) Create(ctx context.Context, repoId string, name string) (*core.Branch, error) {
	log := logger.WithAction("创建分支")

	ok, err := b.checkName(ctx, repoId, name)
	if err != nil {
		return nil, err
	}

	if ok {
		return nil, errors.New(500001)
	}

	branch := &core.Branch{
		Name:    name,
		RepoId:  repoId,
		Active:  true,
		Counter: 0,
	}

	session := b.Session()
	_, err = session.InsertTran(branch)
	if err != nil {
		log.WithError(err).Error("创建分支失败")
		return nil, err
	}

	return b.Find(ctx, branch.Id)
}

func (b *branchStore) InActivate(ctx context.Context, repoId string, name string) error {
	log := logger.WithAction("激活或停用分支")
	branch, err := b.FindByName(ctx, repoId, name)
	if err != nil {
		return err
	}

	branch.Active = false
	err = b.Update(ctx, branch)
	if err != nil {
		log.WithError(err).Error("更新失败")
		return err
	}

	return nil

}

func (b *branchStore) FindByName(ctx context.Context, repoId string, name string) (*core.Branch, error) {
	log := logger.WithAction("根据分支名称获取分支")

	var item core.Branch
	has, err := b.Query().Eq("name", name).Eq("repo_id", repoId).SelectOne(&item)
	if err != nil {
		log.WithError(err).Error("查询失败")
		return nil, err
	}

	if !has {
		return nil, errors.New(900000)
	}

	return &item, nil
}

func (b *branchStore) Find(ctx context.Context, branchId string) (*core.Branch, error) {
	log := logger.WithAction("获取分支")
	var item core.Branch
	has, err := b.Query().ID(branchId).Get(&item)
	if err != nil {
		log.WithError(err).Error("查询失败")
		return nil, err
	}

	if !has {
		return nil, errors.New(900000)
	}

	return &item, nil
}

func (b *branchStore) Delete(ctx context.Context, branch *core.Branch) error {
	log := logger.WithAction("删除分支")

	session := b.Session()

	_, err := session.DeleteTran(branch)
	if err != nil {
		log.WithError(err).Error("删除失败")
		return err
	}
	return nil
}

func (b *branchStore) Update(ctx context.Context, branch *core.Branch) error {
	log := logger.WithAction("更新分支")
	session := b.Session()
	_ = session.UseBool("active").ID(branch.Id)
	_, err := session.UpdateTran(branch)
	if err != nil {
		log.WithError(err).Error("更新失败")
		return err
	}

	return nil
}

func (b *branchStore) Increment(ctx context.Context, branch *core.Branch) (*core.Branch, error) {
	log := logger.WithAction("存储库增加构建次数")

	//TODO:考虑幂等性问题
	branch.Counter++
	err := b.Update(ctx, branch)
	if err != nil {
		log.WithError(err)
		return nil, err
	}

	return branch, nil
}

func (b *branchStore) checkName(ctx context.Context, repoId string, name string) (bool, error) {
	count, err := b.Query().Eq("repo_id", repoId).Eq("name", name).SelectCount(core.Branch{})
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}
