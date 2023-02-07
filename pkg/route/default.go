package route

import (
	"net/http"
	"regexp"

	"github.com/cloudretic/router/pkg/middleware"
	"github.com/cloudretic/router/pkg/path"
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
	return &stringPart{val}, nil
}

func (part *stringPart) Match(rmc *routeMatchContext, token string) bool {
	if part.val == token {
		return true
	} else {
		return false
	}
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

func (part *wildcardPart) Match(rmc *routeMatchContext, token string) bool {
	if part.param != "" {
		rmc.params[part.param] = token
	}
	return true
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

func (part *regexPart) Match(rmc *routeMatchContext, token string) bool {
	// Match against regex
	matched := part.expr.FindString(token)
	if matched != token {
		return false
	}
	// If a parameter is set, act as a wildcard param.
	if part.param != "" {
		// If a token matched, store the matched value as a route Param
		rmc.params[part.param] = matched
	}
	return true
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
	origExpr string
	mws      []middleware.Middleware
	parts    []Part
	rmc      *routeMatchContext
}

// Tokenize and parse a route expression into a defaultRoute.
//
// See interface Route.
func build_defaultRoute(expr string) (*defaultRoute, error) {
	tokens := path.TokenizeString(expr)
	route := &defaultRoute{
		origExpr: expr,
		mws:      make([]middleware.Middleware, 0),
		parts:    make([]Part, 0),
		rmc:      newRMC(),
	}
	for _, token := range tokens {
		part, err := parse(token)
		if err != nil {
			return nil, err
		}
		if parampart, ok := part.(paramPart); ok {
			pn := parampart.ParameterName()
			if pn != "" {
				route.rmc.params[pn] = ""
			}
		}
	}
	return route, nil
}

// Get a string value unique to the route.
//
// See interface Route.
func (route *defaultRoute) Hash() string {
	return route.origExpr
}

// Get the length of the route.
// For defaultRoutes, this is the total number of Parts it contains.
//
// See interface Route.
func (route *defaultRoute) Length() int {
	return len(route.parts)
}

// Attach middleware to the route. Middleware is handled in attachment order.
//
// See interface Route.
func (route *defaultRoute) Attach(mw middleware.Middleware) {
	route.mws = append(route.mws, mw)
}

// Match a request and update its context.
//
// See interface Route.
func (route *defaultRoute) MatchAndUpdateContext(req *http.Request) *http.Request {
	// req = req.Clone(req.Context())
	route.rmc.reset()
	// Check for path length
	tokens := path.TokenizeString(req.URL.Path)
	if len(tokens) != len(route.parts) {
		return nil
	}
	// Run any attached middleware
	for _, mw := range route.mws {
		if req = mw(req); req == nil {
			return nil
		}
	}
	for i, part := range route.parts {
		if ok := part.Match(route.rmc, tokens[i]); !ok {
			return nil
		}
	}
	return route.rmc.apply(req)
}
