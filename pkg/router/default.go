package router

import (
	"net/http"

	"github.com/decentplatforms/matcha/pkg/middleware"
	"github.com/decentplatforms/matcha/pkg/path"
	"github.com/decentplatforms/matcha/pkg/rctx"
	"github.com/decentplatforms/matcha/pkg/route"
	"github.com/decentplatforms/matcha/pkg/tree"
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
func (rt *defaultRouter) Attach(mws ...middleware.Middleware) {
	rt.mws = append(rt.mws, mws...)
}

// Add a route to the router.
//
// AddRoute is deprecated; use HandleRoute instead.
//
// See interface Router.
func (rt *defaultRouter) AddRoute(r route.Route, h http.Handler) {
	register(rt, r, h)
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
	if h != nil {
		rt.handlers[r.Method()][id] = h
	} else {
		rt.handlers[r.Method()][id] = nil
	}
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
	if h != nil {
		register(rt, r, h)
	} else {
		register(rt, r, nil)
	}
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
	if h != nil {
		register(rt, r, h)
	} else {
		register(rt, r, nil)
	}
}

// Mount mounts a handler at path.
//
// See interface Router.
func (rt *defaultRouter) Mount(rpath string, h http.Handler, methods ...string) error {
	if len(methods) == 0 {
		methods = []string{
			http.MethodPut, http.MethodGet, http.MethodPatch, http.MethodDelete, http.MethodPost,
			http.MethodOptions, http.MethodHead, http.MethodTrace, http.MethodConnect,
		}
	}
	trim := func(w http.ResponseWriter, req *http.Request) *http.Request {
		proxyTo := rctx.GetParam(req.Context(), rctx.PARAM_MOUNTPROXYTO)
		req.URL.Path = proxyTo
		return req
	}
	rpath = path.MakePartial(rpath, rctx.PARAM_MOUNTPROXYTO)
	for _, method := range methods {
		r, err := route.New(method, rpath)
		if err != nil {
			return err
		}
		r.Attach(trim)
		rt.HandleRoute(r, h)
	}
	return nil
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
		handler := rt.handlers[req.Method][leaf_id]
		if handler != nil {
			handler.ServeHTTP(w, reqWithCtx)
		} else {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		rctx.ReturnRequestContext(req)
		return
	}
	rt.notfound.ServeHTTP(w, req)
	return
}
