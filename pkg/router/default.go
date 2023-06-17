package router

import (
	"net/http"

	"github.com/cloudretic/matcha/pkg/middleware"
	"github.com/cloudretic/matcha/pkg/rctx"
	"github.com/cloudretic/matcha/pkg/route"
	"github.com/cloudretic/matcha/pkg/tree"
)

type defaultRouter struct {
	mws       []middleware.Middleware
	routes    map[string]map[int]route.Route
	rtree     *tree.RouteTree
	handlers  map[string]map[int]http.Handler
	notfound  http.Handler
	maxParams int
}

func Default() *defaultRouter {
	return &defaultRouter{
		mws:       make([]middleware.Middleware, 0),
		routes:    make(map[string]map[int]route.Route),
		rtree:     tree.New(),
		handlers:  make(map[string]map[int]http.Handler),
		notfound:  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) }),
		maxParams: rctx.DefaultMaxParams,
	}
}

// Attach middleware to the router.
//
// See interface Router.
func (rt *defaultRouter) Attach(mw middleware.Middleware) {
	rt.mws = append(rt.mws, mw)
}

// Add a route to the router.
//
// See interface Router.
func (rt *defaultRouter) AddRoute(r route.Route, h http.Handler) {
	id := rt.rtree.Add(r)
	if rt.routes[r.Method()] == nil {
		rt.routes[r.Method()] = make(map[int]route.Route)
	}
	rt.routes[r.Method()][id] = r
	if rt.handlers[r.Method()] == nil {
		rt.handlers[r.Method()] = make(map[int]http.Handler)
	}
	rt.handlers[r.Method()][id] = h
}

func register(rt *defaultRouter, r route.Route, h http.Handler) {
	id := rt.rtree.Add(r)
	if rt.routes[r.Method()] == nil {
		rt.routes[r.Method()] = make(map[int]route.Route)
	}
	rt.routes[r.Method()][id] = r
	if rt.handlers[r.Method()] == nil {
		rt.handlers[r.Method()] = make(map[int]http.Handler)
	}
	rt.handlers[r.Method()][id] = h
}

// Add a route to the router.
//
// See interface Router.
func (rt *defaultRouter) Handle(method, path string, h http.Handler) error {
	r, err := route.New(method, path)
	if err != nil {
		return err
	}
	register(rt, r, h)
	return nil
}

// Add a route to the router.
//
// See interface Router.
func (rt *defaultRouter) HandleFunc(method, path string, h http.HandlerFunc) error {
	r, err := route.New(method, path)
	if err != nil {
		return err
	}
	register(rt, r, h)
	return nil
}

// Add a route to the router.
//
// See interface Router.
func (rt *defaultRouter) HandleRoute(r route.Route, h http.Handler) {
	register(rt, r, h)
}

// Add a route to the router.
//
// See interface Router.
func (rt *defaultRouter) HandleRouteFunc(r route.Route, h http.HandlerFunc) {
	register(rt, r, h)
}

// Set the handler for instances where no route is found.
//
// See interface Router.
func (rt *defaultRouter) AddNotFound(h http.Handler) {
	rt.notfound = h
}

// Implements http.Handler.
//
// Serve request using the registered middleware, routes, and handlers.
// Tree Router organizes routes by their 'prefixes' (first path elements) and serves based on the first
// path element of the request. Since wildcard and regex parts do not statically evaluate, they are stored as "*".
func (rt *defaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req = middleware.ExecuteMiddleware(rt.mws, w, req)
	if req == nil {
		return
	}
	leaf_id := rt.rtree.Match(req)
	if leaf_id != tree.NO_LEAF_ID {
		r := rt.routes[req.Method][leaf_id]
		req = rctx.PrepareRequestContext(req, route.NumParams(r))
		reqWithCtx := r.MatchAndUpdateContext(req)
		reqWithCtx = middleware.ExecuteMiddleware(r.Middleware(), w, reqWithCtx)
		if reqWithCtx == nil {
			rctx.ReturnRequestContext(req)
			return
		}
		rt.handlers[req.Method][leaf_id].ServeHTTP(w, reqWithCtx)
		rctx.ReturnRequestContext(req)
		return
	}
	rt.notfound.ServeHTTP(w, req)
	return
}
