package repo

import (
	"context"
	"dubhe-ci/common"
	"dubhe-ci/core"
	"dubhe-ci/errors"
	"dubhe-ci/logger"
	"dubhe-ci/store/shared/db"
)

//创建存储库store
func New(db *db.DB) core.RepositoryStore {
	return &repoStore{db: db}
}

type repoStore struct {
	db *db.DB
}

func (r *repoStore) Session() *db.Session {
	return db.WrapSession(r.db.NewSession())
}

func (r *repoStore) Query() *db.Session {
	return db.WrapSession(r.db.Table(core.Repository{}))
}

func (r *repoStore) List(ctx context.Context, page *common.Page) (*common.Page, error) {
	log := logger.WithAction("查询存储库")

	list := make([]*core.Repository, 0)
	page.Bind(&list, core.Repository{})

	err := r.Query().SelectPage(page)
	if err != nil {
		log.WithError(err).Error("存储查询错误")
		return nil, err
	}

	return page, nil
}

func (r *repoStore) CheckName(ctx context.Context, id string, name string) (bool, error) {
	count, err := r.Query().NeIf(id != "", "id", id).Eq("name", name).SelectCount(core.Repository{})
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (r *repoStore) Create(ctx context.Context, repository *core.Repository) (*core.Repository, error) {
	log := logger.WithAction("创建存储库")

	ok, err := r.CheckName(ctx, "", repository.Name)
	if err != nil {
		return nil, err
	}

	if ok {
		return nil, errors.New(500001)
	}

	repository.Active = true

	session := r.Session()
	_, err = session.InsertTran(repository)
	if err != nil {
		log.WithError(err).Error("插入存储失败")
		return nil, err
	}
	return repository, nil
}

func (r *repoStore) Update(ctx context.Context, repository *core.Repository) error {
	log := logger.WithAction("更新存储库")

	ok, err := r.CheckName(ctx, repository.Id, repository.Name)
	if err != nil {
		log.WithError(err).Error("查询异常")
		return err
	}

	if ok {
		return errors.New(500001)
	}

	dbRepo, err := r.Find(ctx, repository.Id)
	if err != nil {
		return err
	}

	dbRepo.Name = repository.Name
	dbRepo.Description = repository.Description
	dbRepo.Url = repository.Url
	dbRepo.CredentialId = repository.CredentialId
	dbRepo.Timeout = repository.Timeout

	session := r.Session()
	_ = session.UseBool("active").ID(dbRepo.Id)
	_, err = session.UpdateTran(repository)
	if err != nil {
		log.WithError(err).Error("更新失败")
		return err
	}

	return nil
}

func (r *repoStore) Find(ctx context.Context, repoId string) (*core.Repository, error) {
	log := logger.WithAction("获取存储库")

	var item core.Repository
	has, err := r.Query().ID(repoId).Get(&item)
	if err != nil {
		log.WithError(err).Error("查询失败")
		return nil, err
	}

	if !has {
		return nil, errors.New(900000)
	}

	return &item, nil
}

func (r *repoStore) Delete(ctx context.Context, repository *core.Repository) error {
	log := logger.WithAction("删除存储库")

	session := r.Session()

	_, err := session.DeleteTran(repository)
	if err != nil {
		log.WithError(err).Error("删除失败")
		return err
	}
	return nil
}
