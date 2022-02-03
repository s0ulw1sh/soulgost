package muxer

import "net/http"
import "testing"
import "github.com/s0ulw1sh/soulgost/hash"

func TestMux(t *testing.T) {
	var mux Mux
	var tmh TestMuxHandler

	if mux.isinit != false {
		t.Error("Mux must be not init")
	}

	mux.Handle("/test/hello", &tmh)

	if len(mux.items) != 1 {
		t.Error("mux.items invalid")
	}

	h := hash.MurMur2([]byte("test"))

	if _, ok := mux.items[h]; !ok {
		t.Error("mux.items not found /test")
	}

}

type TestMuxHandler struct {}

func (self *TestMuxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}