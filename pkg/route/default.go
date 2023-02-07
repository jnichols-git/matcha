package route

import (
	"net/http"

	"github.com/cloudretic/router/pkg/middleware"
	"github.com/cloudretic/router/pkg/path"
)

type defaultRoute struct {
	origExpr string
	mws      []middleware.Middleware
	parts    []Part
}

func build_defaultRoute(expr string) (*defaultRoute, error) {
	tokens := path.TokenizeString(expr)
	route := &defaultRoute{
		origExpr: expr,
		mws:      make([]middleware.Middleware, 0),
		parts:    make([]Part, 0),
	}
	for _, token := range tokens {
		if part, err := parse(token); err != nil {
			return nil, err
		} else {
			route.parts = append(route.parts, part)
		}
	}
	return route, nil
}

func (route *defaultRoute) Hash() string {
	return route.origExpr
}

func (route *defaultRoute) Attach(middleware.Middleware) {

}

func (route *defaultRoute) MatchAndUpdateContext(req *http.Request) *http.Request {
	req = req.Clone(req.Context())
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
		if req = part.Match(req, tokens[i]); req == nil {
			return nil
		}
	}
	return req
}
