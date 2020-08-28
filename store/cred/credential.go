package cred

import (
	"context"
	"dubhe-ci/core"
	"dubhe-ci/errors"
	"dubhe-ci/logger"
	"dubhe-ci/store/shared/db"
)

func New(db *db.DB) core.CredentialStore {
	return &credStore{db: db}
}

type credStore struct {
	db *db.DB
}

func (c *credStore) Session() *db.Session {
	return db.WrapSession(c.db.NewSession())
}

func (c *credStore) Query() *db.Session {
	return db.WrapSession(c.db.Table(core.Credential{}))
}

func (c *credStore) List(ctx context.Context) ([]*core.Credential, error) {
	log := logger.WithAction("查询凭据")

	list := make([]*core.Credential, 0)
	session := c.Query()
	session.Select("id, name, credential_type, username, public_key, created, updated")
	err := session.SelectList(&list)
	if err != nil {
		log.WithError(err).Error("查询错误")
		return nil, err
	}

	return list, err
}

//查询名称是否重复
func (c *credStore) CheckName(ctx context.Context, id string, name string) (bool, error) {
	count, err := c.Query().NeIf(id != "", "id", id).Eq("name", name).SelectCount(core.Credential{})
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (c *credStore) Create(ctx context.Context, credential *core.Credential) (*core.Credential, error) {
	log := logger.WithAction("创建凭据")
	ok, err := c.CheckName(ctx, "", credential.Name)
	if err != nil {
		return nil, err
	}

	if ok {
		return nil, errors.New(500001)
	}

	if credential.CredentialType == 1 && (credential.Username == "" || credential.Password == "") {
		return nil, errors.New(500002)
	} else if credential.CredentialType == 2 && (credential.PublicKey == "" || credential.PrivateKey == "") {
		return nil, errors.New(500003)
	} else if credential.CredentialType == 3 && (credential.Username == "" || credential.Password == "" || credential.Address == "") {
		return nil, errors.New(500002)
	}

	_, err = c.Session().InsertTran(credential)
	if err != nil {
		log.WithError(err).Error("创建失败")
		return nil, err
	}
	return credential, nil
}

func (c *credStore) Update(ctx context.Context, credential *core.Credential) error {
	log := logger.WithAction("创建凭据")
	ok, err := c.CheckName(ctx, "", credential.Name)
	if err != nil {
		return err
	}

	if ok {
		return errors.New(500001)
	}

	if credential.CredentialType == 1 && (credential.Username == "" || credential.Password == "") {
		return errors.New(500002)
	} else if credential.CredentialType == 2 && (credential.PublicKey == "" || credential.PrivateKey == "") {
		return errors.New(500003)
	}

	dbCred, err := c.Find(ctx, credential.Id)
	if err != nil {
		return err
	}

	dbCred.Name = credential.Name
	dbCred.CredentialType = credential.CredentialType
	if dbCred.CredentialType == 1 {
		dbCred.Username = credential.Username
		dbCred.Password = credential.Password
		dbCred.PublicKey = ""
		dbCred.PrivateKey = ""
	} else if dbCred.CredentialType == 2 {
		dbCred.Username = ""
		dbCred.Password = ""
		dbCred.PublicKey = credential.PublicKey
		dbCred.PrivateKey = credential.PrivateKey
	}

	session := c.Session()
	_, err = session.UpdateTran(dbCred)
	if err != nil {
		log.WithError(err).Error("更新失败")
		return err
	}

	return nil
}

func (c *credStore) Find(ctx context.Context, credId string) (*core.Credential, error) {
	log := logger.WithAction("获取凭据")

	var item core.Credential
	has, err := c.Query().ID(credId).Get(&item)
	if err != nil {
		log.WithError(err).Error("查询失败")
		return nil, err
	}

	if !has {
		return nil, errors.New(900000)
	}

	return &item, nil
}

func (c *credStore) Delete(ctx context.Context, credential *core.Credential) error {
	log := logger.WithAction("删除凭据")

	session := c.Session()

	_, err := session.DeleteTran(credential)
	if err != nil {
		log.WithError(err).Error("删除失败")
		return err
	}
	return nil
}

func (c *credStore) ListRegistryCred(ctx context.Context) ([]*core.Credential, error) {
	log := logger.WithAction("查询所有docker凭据")

	list := make([]*core.Credential, 0)
	session := c.Query()
	err := session.Eq("credential_type", 3).SelectList(&list)
	if err != nil {
		log.WithError(err).Error("查询错误")
		return nil, err
	}

	return list, err
}
