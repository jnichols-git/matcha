package route

import (
	"github.com/jnichols-git/matcha/v2/internal/route/require"
	"github.com/jnichols-git/matcha/v2/pkg/cors"
	"github.com/jnichols-git/matcha/v2/pkg/middleware"
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
func WithMiddleware(mws ...middleware.Middleware) ConfigFunc {
	return func(r Route) error {
		r.Attach(mws...)
		return nil
	}
}

func Require(rs ...require.Required) ConfigFunc {
	return func(r Route) error {
		r.Require(rs...)
		return nil
	}
}
