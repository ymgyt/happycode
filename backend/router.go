package backend

import (
	"net/http"

	"github.com/benbjohnson/immutable"
)

type routingEntry struct {
	path    string
	handler http.Handler
}

func (re *routingEntry) score(req *http.Request) int {
	// TODO: refactor
	rp := req.URL.Path
	var score int
	for i := 0; i < len(re.path) && i < len(rp); i++ {
		if rp[i] != re.path[i] {
			if re.path[i] == '*' {
				score = i
			}
			break
		}
		if i == len(re.path)-1 && rp[i] == re.path[i] {
			score = i + 1
		}
	}
	return score
}

type Router struct {
	entries map[string][]*routingEntry
	cache   *immutable.Map
}

func NewRouter() *Router {
	r := &Router{
		entries: make(map[string][]*routingEntry),
		cache:   immutable.NewMap(nil),
	}
	return r
}

func (r *Router) GET(path string, h http.HandlerFunc) {
	r.addEntry("GET", path, h)
}

func (r *Router) addEntry(method, path string, h http.Handler) {
	r.entries[method] = append(r.entries[method], &routingEntry{path: path, handler: h})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h := r.lookup(req)
	h.ServeHTTP(w, req)
}

func (r *Router) lookup(req *http.Request) http.Handler {
	h, found := r.lookupFromCache(req)
	if found {
		return h
	}
	entries := r.entries[req.Method]
	if entries == nil {
		return notFoundHandler
	}

	var maxScore = -1
	for _, entry := range entries {
		score := entry.score(req)
		if score > maxScore {
			maxScore = score
			h = entry.handler
		}
	}
	if maxScore > -1 {
		r.setCache(req, h)
		return h
	}
	return notFoundHandler
}

func (r *Router) lookupFromCache(req *http.Request) (http.Handler, bool) {
	h, found := r.cache.Get(r.cacheKey(req))
	if found {
		return h.(http.Handler), found
	}
	return nil, found
}

func (r *Router) setCache(req *http.Request, h http.Handler) {
	r.cache = r.cache.Set(r.cacheKey(req), h)
}

func (r *Router) cacheKey(req *http.Request) string {
	return req.Method + "_" + req.URL.Path
}

var notFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("not found :("))
})
