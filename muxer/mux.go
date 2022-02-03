package muxer

import (
	"net/http"

	"github.com/s0ulw1sh/soulgost/hash"
)

type Mux struct {
	isinit bool
	items map[uint32]http.Handler
}

func (self *Mux) Handle(pattern string, handler http.Handler) {

	if !self.isinit {
		self.isinit  = true
		self.items = make(map[uint32]http.Handler)
	}

	if len(pattern) == 0 || pattern[0] != '/' { return }

	p := pattern[1:]

	if len(p) == 0 {
		self.items[0] = handler
		return
	}

	for i, c := range p {
		if c == '/' {
			p = p[:i]
			break
		}
	}

	h := hash.MurMur2([]byte(p))

	self.items[h] = handler
}

func (self *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) == 0 || r.URL.Path[0] != '/' { return }

	p := r.URL.Path[1:]

	if len(p) == 0 {
		if s, ok := self.items[0]; ok {
			s.ServeHTTP(w, r)
			return
		} else {
			http.NotFound(w, r)
			return
		}
	}

	for i, c := range p {
		if c == '/' && i > 0 {
			p = p[:i]
			break
		}
	}

	h := hash.MurMur2([]byte(p))

	if s, ok := self.items[h]; ok {
		s.ServeHTTP(w, r)
		return
	}

	http.NotFound(w, r)
}