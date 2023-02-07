package router

import (
	"net/http"

	"github.com/cloudretic/router/pkg/route"
)

// ConfigFuncs run on Routers, usually to add a route or attach middleware.
type ConfigFunc func(rt Router)

func WithRoute(r route.Route, h http.Handler) ConfigFunc {
	return func(rt Router) {
		rt.AddRoute(r, h)
	}
}

func WithNotFound(h http.Handler) ConfigFunc {
	return func(rt Router) {
		rt.AddNotFound(h)
	}
}
