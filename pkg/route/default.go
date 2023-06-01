package route

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"regexp"

	"github.com/cloudretic/matcha/pkg/middleware"
	"github.com/cloudretic/matcha/pkg/path"
	"github.com/cloudretic/matcha/pkg/rctx"
	"github.com/cloudretic/matcha/pkg/route/require"
)

//=====PARTS=====

// LITERALS

// Literal route Parts match and pass without additionally transforming
// the request.

// stringPart; literal string part
type stringPart struct {
	val string
}

func build_stringPart(val string) (*stringPart, error) {
	if url.PathEscape(val[1:]) != val[1:] {
		return nil, errors.New("static part " + val + " is not a valid URL part (" + url.PathEscape(val) + ")")
	} else {
		return &stringPart{val}, nil
	}
}

// stringParts match a literal token exactly.
func (part *stringPart) Match(ctx context.Context, token string) bool {
	if part.val == token {
		return true
	} else {
		return false
	}
}

func (part *stringPart) Eq(other Part) bool {
	if otherSp, ok := other.(*stringPart); ok {
		return otherSp.val == part.val
	}
	return false
}

// WILDCARDS

// Wildcard route Parts store parameters for use by the router in handlers.
// They use a syntax of [wildcard] to denote their name, and can additionally
// be qualified by some conditions by splitting with the : character.

// wildcardParts always match, and add the token as a request param.
type wildcardPart struct {
	param string
}

func build_wildcardPart(param string) (*wildcardPart, error) {
	return &wildcardPart{param}, nil
}

// wildcardParts match any value.
func (part *wildcardPart) Match(ctx context.Context, token string) bool {
	token = token[1:]
	if ctx != nil && part.param != "" {
		rctx.SetParam(ctx, part.param, token)
		// setParam(ctx, part.param, token)
	}
	return true
}

func (part *wildcardPart) Eq(other Part) bool {
	if otherWp, ok := other.(*wildcardPart); ok {
		return otherWp.param == part.param
	}
	return false
}

func (part *wildcardPart) ParameterName() string {
	return part.param
}

func (part *wildcardPart) SetParameterName(s string) {
	part.param = s
}

// regexParts match against regular expressions.
// They're created using the syntax [wildcard]:{regex}
type regexPart struct {
	param string
	expr  *regexp.Regexp
}

func build_regexPart(param, expr string) (*regexPart, error) {
	expr_compiled, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	} else {
		return &regexPart{param, expr_compiled}, nil
	}
}

// regexParts match any token that matches the regex value.
func (part *regexPart) Match(ctx context.Context, token string) bool {
	token = token[1:]
	// Match against regex
	matched := part.expr.FindString(token)
	if matched != token {
		return false
	}
	// If a parameter is set, act as a wildcard param.
	if ctx != nil && part.param != "" {
		// If a token matched, store the matched value as a route Param
		rctx.SetParam(ctx, part.param, token)
	}
	return true
}

func (part *regexPart) Eq(other Part) bool {
	if otherRp, ok := other.(*regexPart); ok {
		return otherRp.expr.String() == part.expr.String() && otherRp.param == part.param
	}
	return false
}

func (part *regexPart) ParameterName() string {
	return part.param
}

func (part *regexPart) SetParameterName(s string) {
	part.param = s
}

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
func build_defaultRoute(method, expr string) (*defaultRoute, error) {
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
		part, err := parse(token)
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

// Get the route prefix.
//
// See interface Route.
func (route *defaultRoute) Prefix() string {
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
		if ok := part.Match(req.Context(), token); !ok {
			return nil
		}
		partIdx++
		if next == -1 || partIdx >= route.Length() {
			break
		}
	}
	//return req.WithContext(route.ctx)
	return req
}

func (route *defaultRoute) Attach(mw middleware.Middleware) {
	route.middleware = append(route.middleware, mw)
}

func (route *defaultRoute) Require(v require.Required) {
	route.required = append(route.required, v)
}

func (route *defaultRoute) Middleware() []middleware.Middleware {
	return route.middleware
}

func (route *defaultRoute) Required() []require.Required {
	return route.required
}
