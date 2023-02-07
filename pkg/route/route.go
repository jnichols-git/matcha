package route

import (
	"net/http"
)

type Route interface {
	Hash() string
	MatchAndUpdateContext(*http.Request) *http.Request
}

// Create a new Route based on a string expression.
func New(expr string) (Route, error) {
	return build_defaultRoute(expr)
}

// Create a new Route based on a string expression, and panic if this fails.
func ForceNew(expr string) Route {
	r, err := build_defaultRoute(expr)
	if err != nil {
		panic(err)
	}
	return r
}
