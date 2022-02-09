package api

import (
	"io"
	"context"
	"strings"
	"github.com/s0ulw1sh/soulgost/hash"
)

var (
	apis = make(map[uint32]ApiType)
)

func Register(name string, apist ApiType) {
	h := hash.MurMur2([]byte(strings.ToLower(name)))
	apis[h] = apist
}

type ApiType interface {
	CallApi(string, Request, Response)
}

type Request struct {
	Ctx(context.Context) context.Context
	Params(interface{}) error
}

type Response struct {
	WriteResult(interface{}, error) error
}