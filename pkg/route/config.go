package route

import (
	"github.com/cloudretic/router/pkg/cors"
	"github.com/cloudretic/router/pkg/middleware"
)

// RouteConfigFuncs can be applied to a Route at creation.
type ConfigFunc func(Route) error

func CORSHeaders(aco *cors.AccessControlOptions) ConfigFunc {
	return func(r Route) error {
		r.Attach(cors.CORSMiddleware(aco))
		return nil
	}
}

func WithMiddleware(mw middleware.Middleware) ConfigFunc {
	return func(r Route) error {
		r.Attach(mw)
		return nil
	}
}
