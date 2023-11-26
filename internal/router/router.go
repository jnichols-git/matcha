package router

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2/internal/rctx"
	"github.com/jnichols-git/matcha/v2/internal/route"
	"github.com/jnichols-git/matcha/v2/pkg/middleware"
	"github.com/jnichols-git/matcha/v2/pkg/path"
	"github.com/jnichols-git/matcha/v2/pkg/tree"
)

type Router struct {
	mws       []middleware.Middleware
	rtree     *tree.RouteTree
	routes    map[int64]*route.Route
	handlers  map[int64]http.Handler
	notfound  http.Handler
	maxParams int
}

func Default() *Router {
	return &Router{
		mws:       make([]middleware.Middleware, 0),
		rtree:     tree.New(),
		routes:    make(map[int64]*route.Route),
		handlers:  make(map[int64]http.Handler),
		notfound:  http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) }),
		maxParams: rctx.DefaultMaxParams,
	}
}

// Attach middleware to the router.
//
// See interface Router.
func (rt *Router) Use(mws ...middleware.Middleware) {
	rt.mws = append(rt.mws, mws...)
}

func register(rt *Router, r *route.Route, h http.Handler) {
	id := rt.rtree.Add(r)
	rt.routes[id] = r
	if h != nil {
		rt.handlers[id] = h
	} else {
		rt.handlers[id] = nil
	}
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
	trim := middleware.TrimPrefix(rpath)
	rpath = path.MakePartial(rpath, "")
	for _, method := range methods {
		r, err := route.New(method, rpath)
		if err != nil {
			return err
		}
		r.Use(trim)
		rt.HandleRoute(r, h)
	}
	return nil
}

// Set the handler for instances where no route is found.
//
// See interface Router.
func (rt *Router) AddNotFound(h http.Handler) {
	rt.notfound = h
}

// Implements http.Handler.
//
// Serve request using the registered middleware, routes, and handlers.
// Tree Router organizes routes by their 'prefixes' (first path elements) and serves based on the first
// path element of the request. Since wildcard and regex parts do not statically evaluate, they are stored as "*".
func (rt *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req = middleware.ExecuteMiddleware(rt.mws, w, req)
	if req == nil {
		return
	}
	leaf_id := rt.rtree.Match(req)
	if leaf_id != tree.NO_LEAF_ID {
		r := rt.routes[leaf_id]
		req = rctx.PrepareRequestContext(req, route.NumParams(r))
		reqWithCtx := r.Execute(req)
		reqWithCtx = middleware.ExecuteMiddleware(r.Middleware(), w, reqWithCtx)
		if reqWithCtx == nil {
			rctx.ReturnRequestContext(req)
			return
		}
		handler := rt.handlers[leaf_id]
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
