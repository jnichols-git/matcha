// Package route defines the API for creating and interacting with a route, and its internal behavior.
//
// See [https://github.com/cloudretic/matcha/blob/main/docs/routes.md].
package route

import (
	"net/http"

	"github.com/cloudretic/matcha/pkg/middleware"
	"github.com/cloudretic/matcha/pkg/route/require"
)

type Route interface {
	// Get a prefix for the route.
	//
	// Route implementations should return a literal string for the first Part in the route, or "*" if not applicable.
	Prefix() string
	// Get a unique hash value for the route.
	//
	// Route implementations must ensure Hash is always unique for two different Routes.
	Hash() string
	// Get the length of the route.
	//
	// Route implementations may determine how to represent their own length.
	Length() int
	// Get the parts of the route.
	//
	// Route implementations must return an exact slice of their parts.
	Parts() []Part
	// Get the method of the route.
	//
	// Route implementations must return a nonempty string containing exactly one method, compliant with http.MethodX
	Method() string
	// Match a request and update its context.
	//
	// Route implementations must return nil if a request does not match the Route, but may otherwise define any return behavior.
	MatchAndUpdateContext(*http.Request) *http.Request
	// Attach middleware to the route.
	//
	// Middleware cannot be removed from a router once it is added.
	Attach(mw middleware.Middleware)
	// Attach a validator to the route.
	//
	// Validators cannot be removed from a router once they are added.
	Require(v require.Required)
	// Get the middleware attached to the route.
	Middleware() []middleware.Middleware
	// Get the validators attached to the route.
	Required() []require.Required
}

// Create a new Route based on a string expression.
func New(method, expr string, confs ...ConfigFunc) (Route, error) {
	// Determine route type
	var r Route
	var err error
	if isPartialRouteExpr(expr) {
		r, err = build_partialRoute(method, expr)
	} else {
		r, err = build_defaultRoute(method, expr)
	}
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
func Declare(method, expr string, confs ...ConfigFunc) Route {
	r, err := New(method, expr, confs...)
	if err != nil {
		panic(err)
	}
	return r
}
