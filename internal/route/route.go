// Package route defines the API for creating and interacting with a route, and its internal behavior.
//
// See [https://github.com/jnichols-git/matcha/v2/blob/main/docs/routes.md].
package route

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2/internal/route/require"
	"github.com/jnichols-git/matcha/v2/pkg/middleware"
	"github.com/jnichols-git/matcha/v2/pkg/path"
	"github.com/jnichols-git/matcha/v2/pkg/rctx"
)

// Route is the default behavior for router, which is to match requests exactly.
type Route struct {
	origExpr   string
	method     string
	parts      []Part
	middleware []middleware.Middleware
	required   []require.Required
}

// Tokenize and parse a route expression into a Route.
//
// See interface Route.
func ParseRoute(method, expr string) (*Route, error) {
	route := &Route{
		origExpr:   "",
		method:     method,
		parts:      make([]Part, 0),
		middleware: make([]middleware.Middleware, 0),
		required:   make([]require.Required, 0),
	}
	var token string
	for next := 0; next < len(expr); {
		token, next = path.Next(expr, next)
		route.origExpr += token
		part, err := ParsePart(token)
		if err != nil {
			return nil, err
		}
		route.parts = append(route.parts, part)
		if next == -1 {
			break
		}
	}
	return route, nil
}

// Get a string value unique to the route.
//
// See interface Route.
func (route *Route) String() string {
	return route.method + " " + route.origExpr
}

// Get the length of the route.
// For Routes, this is the total number of Parts it contains.
//
// See interface Route.
func (route *Route) Length() int {
	return len(route.parts)
}

// Get the parts of the route.
//
// See interface Route.
func (route *Route) Parts() []Part {
	return route.parts
}

// Return the route method.
//
// See interface Route.
func (route *Route) Method() string {
	return route.method
}

// Match a request and update its context.
//
// See interface Route.
func (route *Route) MatchAndUpdateContext(req *http.Request) *http.Request {
	if req.Method != route.method {
		return nil
	}
	// route.ctx.ResetOnto(req.Context())
	// Check for path length
	expr := req.URL.Path
	rctx.ResetRequestContext(req)

	var token string
	var partIdx int
	for next := 0; next < len(expr); {
		part := route.parts[partIdx]
		token, next = path.Next(expr, next)
		if ok := part.Match(token); !ok {
			return nil
		}
		if param := part.Parameter(); param != "" {
			val := token[1:]
			if part.Multi() {
				val = rctx.GetParam(req.Context(), param) + "/" + val
			}
			rctx.SetParam(req.Context(), param, val)
		}
		if !part.Multi() {
			partIdx++
		}
		if next == -1 || partIdx >= route.Length() {
			break
		}
	}
	return req
}

func (route *Route) Use(mws ...middleware.Middleware) *Route {
	route.middleware = append(route.middleware, mws...)
	return route
}

func (route *Route) Require(rs ...require.Required) *Route {
	route.required = append(route.required, rs...)
	return route
}

func (route *Route) Middleware() []middleware.Middleware {
	return route.middleware
}

func (route *Route) Required() []require.Required {
	return route.required
}

// Create a new Route based on a string expression.
func New(method, expr string) (*Route, error) {
	// Determine route type
	var r *Route
	var err error
	r, err = ParseRoute(method, expr)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Create a new Route based on a string expression, and panic if this fails.
// You should not use this unless you are creating a route on program start and do not intend to modify the route after the fact.
func Declare(method, expr string) *Route {
	r, err := New(method, expr)
	if err != nil {
		panic(err)
	}
	return r
}
