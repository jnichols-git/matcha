package router

import (
	"net/http"

	"github.com/cloudretic/router/pkg/route"
)

type defaultRouter struct {
	routes   []route.Route
	handlers map[string]http.Handler
}

func Default() *defaultRouter {
	return &defaultRouter{
		routes:   make([]route.Route, 0),
		handlers: make(map[string]http.Handler),
	}
}

func (rt *defaultRouter) AddRoute(r route.Route, h http.Handler) {
	rt.routes = append(rt.routes, r)
	rt.handlers[r.Hash()] = h
}

func (rt *defaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, r := range rt.routes {
		reqWithContext := r.MatchAndUpdateContext(req)
		if reqWithContext != nil {
			rt.handlers[r.Hash()].ServeHTTP(w, reqWithContext)
			return
		}
	}
}
