package router

import (
	"fmt"
	"net/http"

	"github.com/cloudretic/router/pkg/cors"
	"github.com/cloudretic/router/pkg/middleware"
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

// Give a default set of CORS headers.
func DefaultCORS(aco *cors.AccessControlOptions) ConfigFunc {
	return func(rt Router) error {
		rt.Attach(cors.CORSMiddleware(aco))
		return nil
	}
}

// Handle preflight requests on a given route expression.
// This will respond to OPTIONS requests on the route with a 204 (No Content) and a set of headers
// detailing the provided access control options.
// Fails if the provided route expression is invalid (see Route documentation)
func PreflightCORS(expr string, aco *cors.AccessControlOptions) ConfigFunc {
	return func(rt Router) error {
		r, h, err := cors.PreflightHandler(expr, aco)
		fmt.Println(r.Hash())
		if err != nil {
			return err
		}
		rt.AddRoute(r, h)
		return nil
	}
}

// Attach generic middleware to the Router
func WithMiddleware(mw middleware.Middleware) ConfigFunc {
	return func(rt Router) error {
		rt.Attach(mw)
		return nil
	}
}
