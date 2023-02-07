package route

import (
	"net/http"

	"github.com/cloudretic/go-collections/pkg/slices"
)

// RouteConfigFuncs can be applied to a Route at creation.
type RouteConfigFunc func(Route)

func WithMethods(methods ...string) func(Route) {
	return func(r Route) {
		r.Attach(func(r *http.Request) *http.Request {
			method := r.Method
			if !slices.Contains(methods, method) {
				return nil
			}
			return r
		})
	}
}
