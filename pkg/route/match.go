package route

import (
	"net/http"

	"github.com/cloudretic/router/pkg/router/params"
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
	return params.Set(req, rmc.params)
}
