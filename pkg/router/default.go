package router

import (
	"net/http"

	"github.com/cloudretic/router/pkg/middleware"
	"github.com/cloudretic/router/pkg/path"
	"github.com/cloudretic/router/pkg/rctx"
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
	prefix := r.Prefix()
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

// Implements http.Handler.
//
// Serve request using the registered middleware, routes, and handlers.
// Default Router organizes routes by their 'prefixes' (first path elements) and serves based on the first
// path element of the request. Since wildcard and regex parts do not statically evaluate, they are stored as "*".
func (rt *defaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req = executeMiddleware(rt.mws, w, req)
	if req == nil {
		return
	}
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	reqPrefix, _ := path.Next(req.URL.Path, 0)
	var routes []route.Route
	if rts, ok := rt.routes[reqPrefix]; ok {
		routes = rts
	} else if rts, ok := rt.routes["*"]; ok {
		routes = rts
	}
	for _, r := range routes {
		reqWithCtx := r.MatchAndUpdateContext(req)
		if reqWithCtx != nil {
			reqWithCtx = executeMiddleware(r.Middleware(), w, reqWithCtx)
			if reqWithCtx == nil {
				return
			}
			rt.handlers[r.Hash()].ServeHTTP(w, reqWithCtx)
			return
		}
	}
	rt.notfound.ServeHTTP(w, req)
	return
}

func executeMiddleware(mw []middleware.Middleware, w http.ResponseWriter, req *http.Request) *http.Request {
	for _, m := range mw {
		if req = m(w, req); req == nil {
			return nil
		}
	}
	return req
}