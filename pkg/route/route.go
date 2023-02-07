package route

import (
	"net/http"

	"github.com/cloudretic/router/pkg/middleware"
)

type Route interface {
	Hash() string
	Attach(middleware.Middleware)
	MatchAndUpdateContext(*http.Request) *http.Request
}

// Create a new Route based on a string expression.
func New(expr string, confs ...RouteConfigFunc) (Route, error) {
	r, err := build_defaultRoute(expr)
	if err != nil {
		return nil, err
	}
	for _, conf := range confs {
		conf(r)
	}
	return r, nil
}

// Create a new Route based on a string expression, and panic if this fails.
func ForceNew(expr string, confs ...RouteConfigFunc) Route {
	r, err := New(expr, confs...)
	if err != nil {
		panic(err)
	}
	return r
}
