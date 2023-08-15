package router

import (
	"net/http"

	"github.com/decentplatforms/matcha/pkg/cors"
	"github.com/decentplatforms/matcha/pkg/middleware"
	"github.com/decentplatforms/matcha/pkg/route"
)

// ConfigFuncs run on Routers, usually to add a route or attach middleware.
type ConfigFunc func(rt Router) error

// Add a Route for the Router to handle.
//
// AddRoute was deprecated in v1.2.0. Use HandleRoute instead.
func WithRoute(r route.Route, h http.Handler) ConfigFunc {
	return func(rt Router) error {
		rt.AddRoute(r, h)
		return nil
	}
}

func Handle(method, path string, h http.Handler) ConfigFunc {
	return func(rt Router) error {
		return rt.Handle(method, path, h)
	}
}

func HandleFunc(method, path string, h http.HandlerFunc) ConfigFunc {
	return func(rt Router) error {
		return rt.HandleFunc(method, path, h)
	}
}

func HandleRoute(r route.Route, h http.Handler) ConfigFunc {
	return func(rt Router) error {
		rt.HandleRoute(r, h)
		return nil
	}
}

func HandleRouteFunc(r route.Route, h http.HandlerFunc) ConfigFunc {
	return func(rt Router) error {
		rt.HandleRouteFunc(r, h)
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
func DefaultCORSHeaders(aco *cors.AccessControlOptions) ConfigFunc {
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
		r, err := route.New(http.MethodOptions, expr)
		if err != nil {
			return err
		}
		f := func(w http.ResponseWriter, r *http.Request) {
			cors.SetCORSResponseHeaders(w, r, aco)
			w.WriteHeader(http.StatusNoContent)
		}
		rt.AddRoute(r, http.HandlerFunc(f))
		return nil
	}
}

// Attach generic middleware to the Router
func WithMiddleware(mws ...middleware.Middleware) ConfigFunc {
	return func(rt Router) error {
		rt.Attach(mws...)
		return nil
	}
}
