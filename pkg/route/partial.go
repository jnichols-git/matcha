package route

import (
	"net/http"
	"strings"

	"github.com/cloudretic/router/pkg/middleware"
	"github.com/cloudretic/router/pkg/path"
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

// partialEndPart assumes that it's starting at the first partial token.
// For example, in route /file/[filename]{.+}, partialEndPart will start on any token after file
func (part *partialEndPart) Match(rmc *routeMatchContext, token string) bool {
	ok := part.subPart.Match(rmc, token)
	if !ok {
		return false
	}
	if part.param == "" {
		return true
	}
	rmc.params[part.param] = rmc.params[part.param] + token
	// If there's a match, get the current path from params and append the token
	return true
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
	origExpr string
	mws      []middleware.Middleware
	parts    []Part
	rmc      *routeMatchContext
}

// Tokenize and parse a route expression into a partialRoute.
//
// See interface Route.
func build_partialRoute(expr string) (*partialRoute, error) {
	route := &partialRoute{
		origExpr: expr,
		mws:      make([]middleware.Middleware, 0),
		parts:    make([]Part, 0),
		rmc:      newRMC(),
	}

	tokenCt := strings.Count(expr, "/")
	var token string
	var partIdx int
	for next := 0; next < len(expr); {
		token, next = path.Next(expr, next)
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
		if pp, ok := part.(paramPart); ok {
			pn := pp.ParameterName()
			if pn != "" {
				route.rmc.params[pn] = ""
			}
		}
		route.parts = append(route.parts, part)
		partIdx++
		if next == -1 {
			break
		}
	}
	return route, nil
}

// Get a string value unique to the route.
//
// See interface Route.
func (route *partialRoute) Hash() string {
	return route.origExpr
}

// Get the length of the route.
// For partialRoutes, this is the number of *absolute* parts; the adaptive part at the end is excluded.
// This ensures that when matching for longest route, the more specialized route is always picked.
//
// See interface Route.
func (route *partialRoute) Length() int {
	return len(route.parts) - 1
}

// Attach middleware to the route. Middleware is handled in attachment order.
//
// See interface Route.
func (route *partialRoute) Attach(mw middleware.Middleware) {
	route.mws = append(route.mws, mw)
}

// Match a request and update its context.
// If the request path is longer than the route, partialRoute will do two things:
//   - Check each token beyond the last against the last Part
//   - If the last part is a Wildcard, stores the leftover route as the parameter
//
// See interface Route.
func (route *partialRoute) MatchAndUpdateContext(req *http.Request) *http.Request {
	//req = req.Clone(req.Context())
	route.rmc.reset()
	expr := req.URL.Path
	// check length; tokens should be > parts
	//tokens := path.TokenizeString(req.URL.Path)
	if strings.Count(expr, "/") < len(route.parts)-1 {
		return nil
	}
	// Run any attached middleware
	for _, mw := range route.mws {
		if req = mw(req); req == nil {
			return nil
		}
	}
	var token string
	var partIdx int
	for next := 0; next < len(expr); {
		part := route.parts[partIdx]
		token, next = path.Next(expr, next)
		if ok := part.Match(route.rmc, token); !ok {
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
	return route.rmc.apply(req)
}
