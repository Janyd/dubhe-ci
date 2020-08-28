package db

import (
	"dubhe-ci/config"
	c "dubhe-ci/core"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
)

func Connect(database config.Database) (*DB, error) {
	engine, err := xorm.NewEngine(database.DBType, database.DSN)
	if err != nil {
		return nil, err
	}

	if database.Debug {
		engine.ShowSQL(true)
	}

	err = engine.Ping()
	if err != nil {
		return nil, err
	}

	engine.SetMaxIdleConns(database.MaxIdleConns)
	engine.SetMaxOpenConns(database.MaxOpenConns)
	engine.SetConnMaxLifetime(time.Duration(database.MaxLifetime))
	engine.SetTableMapper(core.SnakeMapper{})

	db := Wrap(engine)

	if err := setupDatabase(db); err != nil {
		return nil, err
	}

	return db, nil
}

func setupDatabase(db *DB) error {
	//同步表结构
	return db.Sync2(
		new(c.Repository),
		new(c.Branch),
		new(c.Build),
		new(c.Step),
		new(c.Credential),
		new(c.Logs),
	)
}
