package api

import (
	"io"
	"errors"
	"strings"
	"context"
	"net/http"
	"encoding/json"
	"github.com/s0ulw1sh/soulgost/proto/ws"
	"github.com/s0ulw1sh/soulgost/hash"
)

const (
	rpcErrParse   = `{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`
)

var (
	ErrRpcParams = errors.New("invalid params")
)

type RpcRequest struct {
	ctx     context.Context    `json:"-"`

	Id      int64              `json:"id"`
	JsonRpc string             `json:"jsonrpc"`
	Method  string             `json:"method"`
	Params  []*json.RawMessage `json:"params"`
}

func (self *RpcRequest) Ctx(c context.Context) context.Context {
	if c != nil {
		self.ctx = c
	}
	return self.ctx
}

func (self *RpcRequest) GetParams(v interface{}) error {
	if len(self.Params) != 1 {
		return ErrRpcParams
	}

	return json.Unmarshal(*(self.Params[0]), v)
}

type RpcResponse struct {
	w       io.Writer        `json:"-"`

	Id      int64            `json:"id"`
	JsonRpc string           `json:"jsonrpc"`
	Result  interface{}      `json:"result,omitempty"`
	Error   interface{}      `json:"error,omitempty"`
}

func (self *RpcResponse) WriteResult(v interface{}, err error) error {

	self.JsonRpc = "2.0"

	if err != nil {
		self.Error = err.Error()
	} else {
		self.Result = v
	}

	return json.NewEncoder(self.w).Encode(*self)
}

type RpcOnBefore       = func(*http.Request, Request) error
type RpcOnWsConnect    = func(Request) error
type RpcOnWsDisconnect = func(Request) error

type Rpc struct {
	BeforeCb    RpcOnBefore
	ConnWsCb    RpcOnWsConnect
	DisconnWsCb RpcOnWsDisconnect
}

func (self *Rpc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var rpcreq RpcRequest

	rpcreq.Ctx(r.Context())

	if self.BeforeCb != nil {
		if err := self.BeforeCb(r, &rpcreq); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
		}
	}

	if strings.HasSuffix(r.URL.Path, "/ws") {
		self.ServeWs(w, r, &rpcreq)
	} else {
		self.ServeRpc(&rpcreq, w, r.Body)
		r.Body.Close()
	}
}

func (self *Rpc) ServeWs(w http.ResponseWriter, r *http.Request, apireq *RpcRequest) {
	var (
		err error
		reader ws.Reader
		writer ws.Writer
	)

	conn, err := ws.WsUpgrader(w, r)

	if err != nil {
		return
	}

	defer conn.Close()

	if (self.ConnWsCb != nil) {
		if err = self.ConnWsCb(apireq); err != nil {
			// TODO send error
			return
		}
	}

	for {
		err = conn.NextReader(&reader)

		if err != nil {
			break
		}

		switch reader.Type {
		case ws.TextMessage:
			conn.NextWriter(ws.TextMessage, &writer)
			self.ServeRpc(apireq, &writer, &reader)
			writer.Close()
		}

	}

	if (self.DisconnWsCb != nil) {
		self.DisconnWsCb(apireq)
	}
}

func (self *Rpc) ServeRpc(req *RpcRequest, w io.Writer, r io.Reader) {
	var res RpcResponse

	err := json.NewDecoder(r).Decode(req)

	if err != nil && err != io.EOF {
		w.Write([]byte(rpcErrParse))
		return
	}

	res.Id = req.Id
	res.w  = w

	servmet := strings.Split(req.Method, ".")

	if len(servmet) != 2 {
		res.WriteResult(nil, ErrMethodNotFound)
		return
	}

	h := hash.MurMur2([]byte(strings.ToLower(servmet[0])))

	if apir, ok := apis[h]; ok {
		apir.CallApi(strings.ToLower(servmet[1]), req, &res)
	} else {
		res.WriteResult(nil, ErrMethodNotFound)
	}
}