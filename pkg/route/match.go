package route

import (
	"context"
	"net/http"
)

//type routeMatchParam [2]string

type routeMatchContext struct {
	params map[string]string
}

func newRMC() *routeMatchContext {
	return &routeMatchContext{
		params: make(map[string]string),
	}
}

func (rmc *routeMatchContext) reset() {
	for k := range rmc.params {
		rmc.params[k] = ""
	}
}

func (rmc *routeMatchContext) apply(req *http.Request) *http.Request {
	ctx := req.Context()
	for k, v := range rmc.params {
		ctx = context.WithValue(ctx, k, v)
	}
	return req.WithContext(ctx)
}
