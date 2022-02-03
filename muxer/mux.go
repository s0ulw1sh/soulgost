package muxer

import (
	"net/http"

	"github.com/s0ulw1sh/soulgost/hash"
)

type Mux struct {
	isinit bool
	items map[uint32]http.Handler
	fallback http.Handler
}

func (self *Mux) Handle(pattern string, handler http.Handler) {
	var h uint32 = 0

	if !self.isinit {
		self.isinit  = true
		self.items = make(map[uint32]http.Handler)
	}

	if len(pattern) == 0 || pattern[0] != '/' { return }

	p := pattern[1:]

	if p == "*" {
		self.fallback = handler
		return
	}

	for i, c := range p {
		if c == '/' {
			p = p[:i]
			break
		}
	}

	if len(p) > 0 {
		h = hash.MurMur2([]byte(p))
	}

	self.items[h] = handler
}

func (self *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h uint32 = 0
	if len(r.URL.Path) == 0 || r.URL.Path[0] != '/' { return }

	p := r.URL.Path[1:]

	for i, c := range p {
		if c == '/' && i > 0 {
			p = p[:i]
			break
		}
	}

	if len(p) != 0 {
		h = hash.MurMur2([]byte(p))
	}

	if s, ok := self.items[h]; ok {
		s.ServeHTTP(w, r)
		return
	}

	if self.fallback != nil {
		self.fallback.ServeHTTP(w, r)
		return
	}

	http.NotFound(w, r)
}