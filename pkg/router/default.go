package router

import (
	"net/http"

	"github.com/cloudretic/router/pkg/middleware"
	"github.com/cloudretic/router/pkg/route"
)

type defaultRouter struct {
	mws      []middleware.Middleware
	routes   []route.Route
	handlers map[string]http.Handler
	notfound http.Handler
}

func Default() *defaultRouter {
	return &defaultRouter{
		mws:      make([]middleware.Middleware, 0),
		routes:   make([]route.Route, 0),
		handlers: make(map[string]http.Handler),
		notfound: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) }),
	}
}

func (rt *defaultRouter) Attach(mw middleware.Middleware) {
	rt.mws = append(rt.mws, mw)
}

func (rt *defaultRouter) AddRoute(r route.Route, h http.Handler) {
	rt.routes = append(rt.routes, r)
	rt.handlers[r.Hash()] = h
}

func (rt *defaultRouter) AddNotFound(h http.Handler) {

}

func (rt *defaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, mw := range rt.mws {
		if req = mw(req); req == nil {
			return
		}
	}
	for _, r := range rt.routes {
		reqWithContext := r.MatchAndUpdateContext(req)
		if reqWithContext != nil {
			rt.handlers[r.Hash()].ServeHTTP(w, reqWithContext)
			return
		}
	}
	rt.notfound.ServeHTTP(w, req)
	return
}
