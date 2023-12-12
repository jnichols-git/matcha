package teaware

import (
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

// Compile generates an http.Handler from a chain of middleware, ending
// in the provided "last" handler.
func Handler(last http.Handler, mws ...Middleware) http.Handler {
	next := last
	for i := len(mws) - 1; i >= 0; i-- {
		next = mws[i](next)
	}
	return next
}
