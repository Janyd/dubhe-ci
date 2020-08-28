package db

import (
	"dubhe-ci/common"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type (
	Session struct {
		*xorm.Session
		params builder.Cond
	}

	SessionWrapper interface {
		//查询
		Query() *Session

		//操作
		Session() *Session
	}
)

func WrapSession(session *xorm.Session) *Session {
	return &Session{
		Session: session,
		params:  builder.NewCond(),
	}
}

func (s *Session) StartTran() error {
	defer s.Close()
	err := s.Begin()
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) CommitTran() error {
	err := s.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) AndCond(cond builder.Cond) *Session {
	s.params = s.params.And(cond)
	return s
}

func (s *Session) OrCond(cond builder.Cond) *Session {
	s.params = s.params.Or(cond)
	return s
}

func (s *Session) Eq(column string, val interface{}) *Session {
	return s.EqIf(true, column, val)
}

func (s *Session) EqIf(condition bool, column string, val interface{}) *Session {
	if condition {
		s.AndCond(builder.Eq{column: val})
	}
	return s
}

func (s *Session) Ne(column string, val interface{}) *Session {
	return s.NeIf(true, column, val)
}

func (s *Session) NeIf(condition bool, column string, val interface{}) *Session {
	if condition {
		s.AndCond(builder.Neq{column: val})
	}
	return s
}

func (s *Session) Gt(column string, val interface{}) *Session {
	return s.GtIf(true, column, val)
}

func (s *Session) GtIf(condition bool, column string, val interface{}) *Session {
	if condition {
		s.AndCond(builder.Gt{column: val})
	}
	return s
}

func (s *Session) Ge(column string, val interface{}) *Session {
	return s.GeIf(true, column, val)
}

func (s *Session) GeIf(condition bool, column string, val interface{}) *Session {
	if condition {
		s.AndCond(builder.Gte{column: val})
	}
	return s
}

func (s *Session) Lt(column string, val interface{}) *Session {
	return s.LtIf(true, column, val)
}

func (s *Session) LtIf(condition bool, column string, val interface{}) *Session {
	if condition {
		s.AndCond(builder.Lt{column: val})
	}
	return s
}

func (s *Session) Le(column string, val interface{}) *Session {
	return s.LeIf(true, column, val)
}

func (s *Session) LeIf(condition bool, column string, val interface{}) *Session {
	if condition {
		s.AndCond(builder.Lte{column: val})
	}
	return s
}

func (s *Session) Between(column string, val1 interface{}, val2 interface{}) *Session {
	return s.BetweenIf(true, column, val1, val2)
}

func (s *Session) BetweenIf(condition bool, column string, val1 interface{}, val2 interface{}) *Session {
	if condition {
		s.AndCond(builder.Between{Col: column, LessVal: val1, MoreVal: val2})
	}
	return s
}

func (s *Session) IN(column string, vals ...interface{}) *Session {
	return s.INIf(true, column, vals...)
}

func (s *Session) INIf(condition bool, column string, vals ...interface{}) *Session {
	if condition {
		s.AndCond(builder.In(column, vals...))
	}
	return s
}

func (s *Session) NotIN(column string, vals ...interface{}) *Session {
	return s.NotINIf(true, column, vals...)
}

func (s *Session) NotINIf(condition bool, column string, vals ...interface{}) *Session {
	if condition {
		s.AndCond(builder.NotIn(column, vals...))
	}
	return s
}

func (s *Session) OrEq(column string, val interface{}) *Session {
	return s.EqIf(true, column, val)
}

func (s *Session) OrEqIf(condition bool, column string, val interface{}) *Session {
	if condition {
		s.OrCond(builder.Eq{column: val})
	}
	return s
}

func (s *Session) OrNe(column string, val interface{}) *Session {
	return s.NeIf(true, column, val)
}

func (s *Session) OrNeIf(condition bool, column string, val interface{}) *Session {
	if condition {
		s.OrCond(builder.Neq{column: val})
	}
	return s
}

func (s *Session) OrGt(column string, val interface{}) *Session {
	return s.GtIf(true, column, val)
}

func (s *Session) OrGtIf(condition bool, column string, val interface{}) *Session {
	if condition {
		s.OrCond(builder.Gt{column: val})
	}
	return s
}

func (s *Session) OrGe(column string, val interface{}) *Session {
	return s.GeIf(true, column, val)
}

func (s *Session) OrGeIf(condition bool, column string, val interface{}) *Session {
	if condition {
		s.OrCond(builder.Gte{column: val})
	}
	return s
}

func (s *Session) OrLt(column string, val interface{}) *Session {
	return s.LtIf(true, column, val)
}

func (s *Session) OrLtIf(condition bool, column string, val interface{}) *Session {
	if condition {
		s.OrCond(builder.Lt{column: val})
	}
	return s
}

func (s *Session) OrLe(column string, val interface{}) *Session {
	return s.LeIf(true, column, val)
}

func (s *Session) OrLeIf(condition bool, column string, val interface{}) *Session {
	if condition {
		s.OrCond(builder.Lte{column: val})
	}
	return s
}

func (s *Session) OrBetween(column string, val1 interface{}, val2 interface{}) *Session {
	return s.BetweenIf(true, column, val1, val2)
}

func (s *Session) OrBetweenIf(condition bool, column string, val1 interface{}, val2 interface{}) *Session {
	if condition {
		s.OrCond(builder.Between{Col: column, LessVal: val1, MoreVal: val2})
	}
	return s
}

func (s *Session) SelectOne(bean interface{}) (bool, error) {
	return s.Where(s.params).Get(bean)
}

func (s *Session) SelectCount(bean interface{}) (int64, error) {
	return s.Where(s.params).Count(bean)
}

func (s *Session) SelectList(list interface{}) error {
	return s.Where(s.params).Find(list)
}

func (s *Session) SelectPage(page *common.Page) error {
	total, err := s.Where(s.params).Count(page.GetBean())
	if err != nil {
		return err
	}

	page.SetTotal(total)

	s.Where(s.params).Limit(int(page.GetSize()), int(page.Offset()))

	if len(page.Desc) > 0 {
		s.Desc(page.Desc...)
	}

	if len(page.Asc) > 0 {
		s.Asc(page.Asc...)
	}

	columns := page.GetColumns()
	if columns != "" {
		s.Select(columns)
	}

	err = s.Find(page.Records)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) InsertTran(beans ...interface{}) (int64, error) {
	err := s.StartTran()
	if err != nil {
		_ = s.Rollback()
		return 0, err
	}

	affected, err := s.Insert(beans...)
	if err != nil {
		_ = s.Rollback()
		return 0, err
	}

	err = s.CommitTran()

	if err != nil {
		_ = s.Rollback()
		return 0, err
	}

	return affected, nil
}

func (s *Session) UpdateTran(bean interface{}, condiBean ...interface{}) (int64, error) {
	err := s.StartTran()
	if err != nil {
		return 0, err
	}

	affected, err := s.Update(bean, condiBean...)
	if err != nil {
		_ = s.Rollback()
		return 0, err
	}
	err = s.CommitTran()
	if err != nil {
		_ = s.Rollback()
	}

	return affected, nil
}

func (s *Session) DeleteTran(bean interface{}) (int64, error) {
	err := s.StartTran()
	if err != nil {
		return 0, err
	}

	affected, err := s.Delete(bean)
	if err != nil {
		_ = s.Rollback()
		return 0, err
	}

	err = s.CommitTran()
	if err != nil {
		_ = s.Rollback()
		return 0, err
	}

	return affected, err
}
