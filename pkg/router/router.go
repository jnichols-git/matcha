package router

import (
	"net/http"

	"github.com/cloudretic/router/pkg/middleware"
	"github.com/cloudretic/router/pkg/route"
)

type Router interface {
	// Attach middleware to the router.
	//
	// Router implementations must run all attached middleware on all incoming requests.
	Attach(mw middleware.Middleware)
	// Add a route to the router.
	//
	// Router implementations must use h for any request that matches r, in order of addition to the router.
	AddRoute(r route.Route, h http.Handler)
	// Add a handler for any request that is not matched.
	//
	// Router implementations should define default behavior, and must allow user assignment of behavior.
	AddNotFound(h http.Handler)
	// Implements http.Handler
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// Create a new Router.
// Returns an error if creation fails.
func New(with Router, cfs ...ConfigFunc) (Router, error) {
	for _, cf := range cfs {
		err := cf(with)
		if err != nil {
			return nil, err
		}
	}
	return with, nil
}

// Declare a new Router.
// Panics if creation fails.
func Declare(with Router, cfs ...ConfigFunc) Router {
	for _, cf := range cfs {
		err := cf(with)
		if err != nil {
			panic(err)
		}
	}
	return with
}
