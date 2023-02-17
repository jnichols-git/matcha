package router

import (
	"net/http"

	"github.com/cloudretic/router/pkg/middleware"
	"github.com/cloudretic/router/pkg/path"
	"github.com/cloudretic/router/pkg/route"
)

type defaultRouter struct {
	mws      []middleware.Middleware
	routes   map[string][]route.Route
	handlers map[string]http.Handler
	notfound http.Handler
}

func Default() *defaultRouter {
	return &defaultRouter{
		mws:      make([]middleware.Middleware, 0),
		routes:   make(map[string][]route.Route),
		handlers: make(map[string]http.Handler),
		notfound: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) }),
	}
}

func (rt *defaultRouter) Attach(mw middleware.Middleware) {
	rt.mws = append(rt.mws, mw)
}

func (rt *defaultRouter) AddRoute(r route.Route, h http.Handler) {
	prefix := r.Part(0).Expr()
	if rt.routes[prefix] != nil {
		rt.routes[prefix] = append(rt.routes[prefix], r)
	} else {
		rt.routes[prefix] = make([]route.Route, 1)
		rt.routes[prefix][0] = r
	}
	rt.handlers[r.Hash()] = h
}

func (rt *defaultRouter) AddNotFound(h http.Handler) {
	rt.notfound = h
}

func (rt *defaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, mw := range rt.mws {
		if req = mw(w, req); req == nil {
			return
		}
	}
	reqPrefix, _ := path.Next(req.URL.Path, 0)
	if routes, ok := rt.routes[reqPrefix]; ok {
		for _, r := range routes {
			reqWithContext := r.MatchAndUpdateContext(req)
			if reqWithContext != nil {
				rt.handlers[r.Hash()].ServeHTTP(w, reqWithContext)
				return
			}
		}
	} else if routes, ok := rt.routes["*"]; ok {
		for _, r := range routes {
			reqWithContext := r.MatchAndUpdateContext(req)
			if reqWithContext != nil {
				rt.handlers[r.Hash()].ServeHTTP(w, reqWithContext)
				return
			}
		}
	}
	rt.notfound.ServeHTTP(w, req)
	return
}
