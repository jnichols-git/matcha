package route

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2/internal/route/require"
	"github.com/jnichols-git/matcha/v2/pkg/middleware"
	"github.com/jnichols-git/matcha/v2/pkg/path"
	"github.com/jnichols-git/matcha/v2/pkg/rctx"
)

// =====ROUTE=====

// defaultRoute is the default behavior for router, which is to match requests exactly.
type defaultRoute struct {
	origExpr   string
	method     string
	parts      []Part
	middleware []middleware.Middleware
	required   []require.Required
}

// Tokenize and parse a route expression into a defaultRoute.
//
// See interface Route.
func ParseRoute(method, expr string) (*defaultRoute, error) {
	route := &defaultRoute{
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
func (route *defaultRoute) Hash() string {
	return route.method + " " + route.origExpr
}

// Get the length of the route.
// For defaultRoutes, this is the total number of Parts it contains.
//
// See interface Route.
func (route *defaultRoute) Length() int {
	return len(route.parts)
}

// Get the parts of the route.
//
// See interface Route.
func (route *defaultRoute) Parts() []Part {
	return route.parts
}

// Return the route method.
//
// See interface Route.
func (route *defaultRoute) Method() string {
	return route.method
}

// Match a request and update its context.
//
// See interface Route.
func (route *defaultRoute) MatchAndUpdateContext(req *http.Request) *http.Request {
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

func (route *defaultRoute) Attach(mws ...middleware.Middleware) {
	route.middleware = append(route.middleware, mws...)
}

func (route *defaultRoute) Require(rs ...require.Required) {
	route.required = append(route.required, rs...)
}

func (route *defaultRoute) Middleware() []middleware.Middleware {
	return route.middleware
}

func (route *defaultRoute) Required() []require.Required {
	return route.required
}
