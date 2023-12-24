package router

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2/internal/path"
	"github.com/jnichols-git/matcha/v2/internal/rctx"
	"github.com/jnichols-git/matcha/v2/internal/tree"
	"github.com/jnichols-git/matcha/v2/route"
	"github.com/jnichols-git/matcha/v2/teaware"
)

type Router struct {
	mws       []teaware.Middleware
	rtree     *tree.RouteTree
	compiled  http.Handler
	routes    map[int64]*route.Route
	handlers  map[int64]http.Handler
	notfound  http.Handler
	maxParams int
}

func Default() *Router {
	r := &Router{
		mws:       make([]teaware.Middleware, 0),
		rtree:     tree.New(),
		routes:    make(map[int64]*route.Route),
		handlers:  make(map[int64]http.Handler),
		notfound:  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) }),
		maxParams: rctx.DefaultMaxParams,
	}
	r.reload()
	return r
}

// Attach middleware to the router.
//
// See interface Router.
func (rt *Router) Use(mws ...teaware.Middleware) *Router {
	rt.mws = append(rt.mws, mws...)
	rt.reload()
	return rt
}

func register(rt *Router, r *route.Route, h http.Handler) {
	id := rt.rtree.Add(r)
	rt.routes[id] = r
	h = teaware.Handler(h, r.Middleware()...)
	if h != nil {
		rt.handlers[id] = h
	} else {
		rt.handlers[id] = nil
	}
	rt.reload()
}

// Add a route to the router.
//
// See interface Router.
func (rt *Router) Handle(method, path string, h http.Handler) error {
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
func (rt *Router) HandleFunc(method, path string, h http.HandlerFunc) error {
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
func (rt *Router) HandleRoute(r *route.Route, h http.Handler) {
	register(rt, r, h)
}

// Add a route to the router.
//
// See interface Router.
func (rt *Router) HandleRouteFunc(r *route.Route, h http.HandlerFunc) {
	if h != nil {
		register(rt, r, h)
	} else {
		register(rt, r, nil)
	}
}

// Mount mounts a handler at path.
//
// See interface Router.
func (rt *Router) Mount(rpath string, h http.Handler, methods ...string) error {
	if len(methods) == 0 {
		methods = []string{
			http.MethodPut, http.MethodGet, http.MethodPatch, http.MethodDelete, http.MethodPost,
			http.MethodOptions, http.MethodHead, http.MethodTrace, http.MethodConnect,
		}
	}
	trim := teaware.TrimPrefix(rpath)
	rpath = path.MakePartial(rpath, "")
	for _, method := range methods {
		r, err := route.New(method, rpath)
		if err != nil {
			return err
		}
		r.Use(trim)
		rt.HandleRoute(r, h)
	}
	rt.reload()
	return nil
}

// Set the handler for instances where no route is found.
//
// See interface Router.
func (rt *Router) HandleNotFound(h http.Handler) {
	rt.notfound = h
}

func (rt *Router) reload() {
	var h http.Handler = http.HandlerFunc(rt.matchRoute)
	h = teaware.Handler(h, rt.mws...)
	rt.compiled = h
}

func (rt *Router) matchRoute(w http.ResponseWriter, req *http.Request) {
	leaf_id := rt.rtree.Match(req)
	if leaf_id == tree.NO_LEAF_ID {
		rt.notfound.ServeHTTP(w, req)
		return
	}
	r := rt.routes[leaf_id]
	req = rctx.PrepareRequestContext(req, route.NumParams(r))
	reqWithCtx := r.Execute(req)
	handler := rt.handlers[leaf_id]
	if handler != nil {
		handler.ServeHTTP(w, reqWithCtx)
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	rctx.ReturnRequestContext(req)
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rt.compiled.ServeHTTP(w, req)
}
