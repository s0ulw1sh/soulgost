package api

import (
	"context"
	"errors"
	"strings"
	"github.com/s0ulw1sh/soulgost/hash"
)

var (
	apis = make(map[uint32]ApiType)

	ErrParams = errors.New("method not found")
)

func Register(name string, apist ApiType) {
	h := hash.MurMur2([]byte(strings.ToLower(name)))
	apis[h] = apist
}

type ApiType interface {
	CallApi(string, Request, Response)
}

type Request interface {
	Ctx(context.Context) context.Context
	GetParams(interface{}) error
}

type Response interface {
	WriteResult(interface{}, error) error
}