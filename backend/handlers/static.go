package handlers

import (
	"net/http"
)

type Static struct {
	fs http.Handler
}

func NewStatic(root, prefix string) *Static {
	s := &Static{
		fs: http.StripPrefix(prefix, http.FileServer(http.Dir(root))),
	}
	return s
}

func (s *Static) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		req.URL.Path = "/static/wasm/wasm_exec.html"
	}
	if req.URL.Path == "/favicon.ico" {
		req.URL.Path = "/static/image/isomorphic_go_icon.png"
	}
	s.fs.ServeHTTP(w, req)
}

