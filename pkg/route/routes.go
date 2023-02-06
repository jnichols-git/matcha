package route

import (
	"net/http"

	"github.com/CloudRETIC/router/pkg/path"
)

type defaultRoute struct {
	parts   []Part
	handler http.Handler
}

func build_defaultRoute(expr string, handler http.Handler) (*defaultRoute, error) {
	tokens := path.TokenizeString(expr)
	route := &defaultRoute{
		parts:   make([]Part, 0),
		handler: handler,
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

func (route *defaultRoute) MatchAndUpdateContext(req *http.Request) *http.Request {
	req = req.Clone(req.Context())
	tokens := path.TokenizeString(req.URL.Path)
	if len(tokens) != len(route.parts) {
		return nil
	}
	for i, part := range route.parts {
		if req = part.Match(req, tokens[i]); req == nil {
			return nil
		}
	}
	return req
}

func (route *defaultRoute) Handler() http.Handler {
	return route.handler
}
