package db

import (
	"errors"
)

var (
	ErrInvalidType = errors.New("soulgost: invalid type")
)

type IField interface {
	Sqlexpr() string
	Val() interface{}
}

type Field struct {
	Name  string
	Value interface{}
}

func (self *Field) Sqlexpr() string {
	return "`" + self.Name + "`=?"
}

func (self *Field) Val() interface{} {
	return self.Value
}

type Cond struct {
	Name  string
	Op    string
	Value interface{}
}

func (self *Cond) Sqlexpr() string {
	return "`" + self.Name + "`" + self.Op + "?"
}

func (self *Cond) Val() interface{} {
	return self.Value
}

type Expr struct {
	Exprstr  string
}

func (self *Expr) Sqlexpr() string {
	return self.Exprstr
}

func (self *Expr) Val() interface{} {
	return nil
}