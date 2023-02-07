package route

import (
	"net/http"

	"github.com/cloudretic/go-collections/pkg/slices"
)

// RouteConfigFuncs can be applied to a Route at creation.
type ConfigFunc func(Route) error

// Filter out requests that don't have a method included in methods
func WithMethods(methods ...string) ConfigFunc {
	return func(r Route) error {
		r.Attach(func(r *http.Request) *http.Request {
			if !slices.Contains(methods, r.Method) {
				return nil
			}
			return r
		})
		return nil
	}
}
