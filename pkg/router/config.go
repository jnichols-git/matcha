package router

import (
	"net/http"

	"github.com/cloudretic/router/pkg/route"
)

// ConfigFuncs run on Routers, usually to add a route or attach middleware.
type ConfigFunc func(rt Router) error

// Add a Route for the Router to handle
func WithRoute(r route.Route, h http.Handler) ConfigFunc {
	return func(rt Router) error {
		rt.AddRoute(r, h)
		return nil
	}
}

// Add a handler for requests that are not handled by any other route in the Router
func WithNotFound(h http.Handler) ConfigFunc {
	return func(rt Router) error {
		rt.AddNotFound(h)
		return nil
	}
}
