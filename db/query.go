package db

import (
	"database/sql"
)

type Query struct {
	Db     *sql.DB
}

type Table struct {
	Name  string
	Alias string
}

type Field struct {
	Name  string
	Value interface{}
}

type Condition struct {
	Name  string
	Cnd   string
	Value interface{}
}

type Expr struct {
	Value string
}

func (self *Query) Table(tables ...string) *Query {
	return self
}

func (self *Query) First(fields []Field, conds ...[]Condition) *Query {
	return self
}

func (self *Query) Count(field string, conds ...[]Condition) *Query {
	return self
}

func (self *Query) Select(fields []Field, conds ...[]Condition) *Query {
	return self
}

func (self *Query) Update(fields []Field, conds ...[]Condition) *Query {
	return self
}

func (self *Query) Insert(fields []Field) *Query {
	return self
}

func (self *Query) Delete(conds ...[]Condition) *Query {
	return self
}

func (self *Query) Limit(lim interface{}) *Query {
	return self
}

func (self *Query) Offset(off interface{}) *Query {
	return self
}

func Q(db *sql.DB) *Query {
	return &Query{db}
}
