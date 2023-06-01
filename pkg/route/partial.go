package route

import (
	"context"
	"net/http"
	"strings"

	"github.com/cloudretic/matcha/pkg/middleware"
	"github.com/cloudretic/matcha/pkg/path"
	"github.com/cloudretic/matcha/pkg/rctx"
	"github.com/cloudretic/matcha/pkg/route/require"
)

// =====PARTS=====

// partialEndPart implements Part to match against a specific subPart repeatedly, with a given optional route parameter.
type partialEndPart struct {
	param   string
	subPart Part
}

// parse a partialEndPart from a token.
func parse_partialEndPart(token string) (*partialEndPart, error) {
	result := &partialEndPart{}
	// get subToken from token (exclude +)
	subToken := token[:len(token)-1]
	// if subToken is empty, use an unqualified anyWord
	if subToken == "/" {
		result.subPart = &regexPart{"", regexp_anyWord_compiled}
		return result, nil
	}
	// otherwise, parse out subToken
	subPart, err := parse(subToken)
	if err != nil {
		return nil, err
	}
	// If the subpart has a parameter, move it to the result.
	// This has no real effect if the subPart has an empty parameter (intended behavior)
	if subPartWithParam, ok := subPart.(paramPart); ok {
		result.param = subPartWithParam.ParameterName()
		subPartWithParam.SetParameterName("")
	}
	result.subPart = subPart

	return result, nil
}

func IsPartialEndPart(p Part) bool {
	_, ok := p.(*partialEndPart)
	return ok
}

// partialEndPart assumes that it's starting at the first partial token.
// For example, in route /file/[filename]{.+}, partialEndPart will start on any token after file
func (part *partialEndPart) Match(ctx context.Context, token string) bool {
	ok := part.subPart.Match(ctx, token)
	if !ok {
		return false
	}
	// If there's a match, get the current path from params and append the token
	if ctx != nil && part.param != "" {
		rctx.SetParam(ctx, part.param, rctx.GetParam(ctx, part.param)+token)
	}
	return true
}

func (part *partialEndPart) Eq(other Part) bool {
	if otherPep, ok := other.(*partialEndPart); ok {
		if otherPep.param != part.param {
			return false
		}
		return part.subPart.Eq(otherPep.subPart)
	}
	return false
}

func (part *partialEndPart) Expr() string {
	return "*"
}

func (part *partialEndPart) ParameterName() string {
	return part.param
}

func (part *partialEndPart) SetParameterName(s string) {
	part.param = s
}

// =====ROUTE=====

// Convenience function to determine if a route expression is partial.
func isPartialRouteExpr(s string) bool {
	return len(s) > 0 && s[len(s)-1] == '+'
}

// partialRoute is specialized to allow routes that may match on extensions, rather than on
// an exact match
type partialRoute struct {
	origExpr   string
	method     string
	parts      []Part
	middleware []middleware.Middleware
	required   []require.Required
}

// Tokenize and parse a route expression into a partialRoute.
//
// See interface Route.
func build_partialRoute(method, expr string) (*partialRoute, error) {
	route := &partialRoute{
		origExpr: "",
		method:   method,
		parts:    make([]Part, 0),
	}

	tokenCt := strings.Count(expr, "/")
	var token string
	var partIdx int
	for next := 0; next < len(expr); {
		token, next = path.Next(expr, next)
		route.origExpr += token
		var part Part
		var err error
		if partIdx < tokenCt-1 {
			part, err = parse(token)
		} else {
			part, err = parse_partialEndPart(token)
		}
		if err != nil {
			return nil, err
		}
		route.parts = append(route.parts, part)
		partIdx++
		if next == -1 {
			break
		}
	}
	return route, nil
}

// Get the route prefix.
//
// See interface Route.
func (route *partialRoute) Prefix() string {
	switch r := route.parts[0].(type) {
	case *stringPart:
		return r.val
	default:
		return "*"
	}
}

// Get a string value unique to the route.
//
// See interface Route.
func (route *partialRoute) Hash() string {
	return route.method + " " + route.origExpr
}

// Get the length of the route.
// For partialRoutes, this is the number of *absolute* parts; the adaptive part at the end is excluded.
// This ensures that when matching for longest route, the more specialized route is always picked.
//
// See interface Route.
func (route *partialRoute) Length() int {
	return len(route.parts) - 1
}

// Get the parts of the route.
//
// See interface Route.
func (route *partialRoute) Parts() []Part {
	return route.parts
}

// Return the route method.
//
// See interface Route.
func (route *partialRoute) Method() string {
	return route.method
}

// Match a request and update its context.
// If the request path is longer than the route, partialRoute will do two things:
//   - Check each token beyond the last against the last Part
//   - If the last part is a Wildcard, stores the leftover route as the parameter
//
// See interface Route.
func (route *partialRoute) MatchAndUpdateContext(req *http.Request) *http.Request {
	if req.Method != route.method {
		return nil
	}
	//route.ctx.ResetOnto(req.Context())
	expr := req.URL.Path
	if strings.Count(expr, "/") < len(route.parts)-1 {
		return nil
	}

	rctx.ResetRequestContext(req)

	var token string
	var partIdx int
	for next := 0; next < len(expr); {
		part := route.parts[partIdx]
		token, next = path.Next(expr, next)
		if ok := part.Match(req.Context(), token); !ok {
			return nil
		}
		if partIdx+1 < len(route.parts) {
			partIdx++
		}
		if next == -1 {
			break
		}
	}
	// If there were no empty tokens to begin with, run the last rou
	return req
}

func (route *partialRoute) Attach(m middleware.Middleware) {
	route.middleware = append(route.middleware, m)
}

func (route *partialRoute) Require(v require.Required) {
	route.required = append(route.required, v)
}

func (route *partialRoute) Middleware() []middleware.Middleware {
	return route.middleware
}

func (route *partialRoute) Required() []require.Required {
	return route.required
}
