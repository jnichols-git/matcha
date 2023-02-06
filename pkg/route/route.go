package route

import (
	"net/http"
)

// The Route interface defines
type Route interface {
	MatchAndUpdateContext(*http.Request) *http.Request
	Handler() http.Handler
}

type RouteSet []Route

// Create a new Route based on a string expression.
func New(expr string, handler http.Handler) (Route, error) {
	return build_defaultRoute(expr, handler)
}
