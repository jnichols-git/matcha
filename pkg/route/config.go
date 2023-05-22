package route

import (
	"github.com/cloudretic/matcha/pkg/cors"
	"github.com/cloudretic/matcha/pkg/middleware"
	"github.com/cloudretic/matcha/pkg/route/require"
)

// RouteConfigFuncs can be applied to a Route at creation.
type ConfigFunc func(Route) error

// Attaches middleware to the route that sets CORS headers on matched requests only.
func CORSHeaders(aco *cors.AccessControlOptions) ConfigFunc {
	return func(r Route) error {
		r.Attach(cors.CORSMiddleware(aco))
		return nil
	}
}

// Attaches middleware to the route.
func WithMiddleware(mw middleware.Middleware) ConfigFunc {
	return func(r Route) error {
		r.Attach(mw)
		return nil
	}
}

func Validators(vs ...require.Required) ConfigFunc {
	return func(r Route) error {
		for _, v := range vs {
			r.Validate(v)
		}
		return nil
	}
}
