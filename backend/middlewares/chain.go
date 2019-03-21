package middlewares

import "net/http"

type Interface interface {
	http.Handler
	SetNext(http.Handler)
}

func NewChain(root http.Handler, mws []Interface) *Chain {
	if len(mws) == 0 {
		panic("empty middlewares")
	}
	for i, mw := range mws {
		if i == len(mws)-1 {
			mw.SetNext(root)
			break
		}
		mw.SetNext(mws[i+1])
	}
	return &Chain{mws: mws}
}

type Chain struct {
	mws []Interface
}

func (chain *Chain) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	chain.mws[0].ServeHTTP(w, req)
}

