package route

import (
	"net/http"

	"github.com/cloudretic/router/pkg/middleware"
)

type Route interface {
	// Get a unique hash value for the route.
	//
	// Route implementations must ensure Hash is always unique for two different Routes.
	Hash() string
	// Attach middleware to the route.
	//
	// Route implementations may define the order that middleware is handled.
	Attach(middleware.Middleware)
	// Match a request and update its context.
	//
	// Route implementations must return nil if a request does not match the Route, but may otherwise define any return behavior.
	MatchAndUpdateContext(*http.Request) *http.Request
}

// Create a new Route based on a string expression.
func New(expr string, confs ...ConfigFunc) (Route, error) {
	r, err := build_defaultRoute(expr)
	if err != nil {
		return nil, err
	}
	for _, conf := range confs {
		err = conf(r)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

// Create a new Route based on a string expression, and panic if this fails.
// You should not use this unless you are creating a route on program start and do not intend to modify the route after the fact.
func NewDecl(expr string, confs ...ConfigFunc) Route {
	r, err := New(expr, confs...)
	if err != nil {
		panic(err)
	}
	for _, conf := range confs {
		err = conf(r)
		if err != nil {
			panic(err)
		}
	}
	return r
}
