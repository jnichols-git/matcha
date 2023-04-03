package route

import "github.com/cloudretic/router/pkg/middleware"

// RouteConfigFuncs can be applied to a Route at creation.
type ConfigFunc func(Route) error

func WithMiddleware(mw middleware.Middleware) ConfigFunc {
	return func(r Route) error {
		r.Attach(mw)
		return nil
	}
}