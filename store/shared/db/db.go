package db

import "xorm.io/xorm"

type (
	DB struct {
		*xorm.Engine
	}
)

func Wrap(e *xorm.Engine) *DB {
	return &DB{e}
}
