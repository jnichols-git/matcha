package middleware

import (
	"fmt"
	"net/http"
)

// Middleware runs on any incoming request. Attachment behavior is defined by the structure it's attached to (route vs. router).
//
// Returns an *http.Request; the middleware can set router params or reject a request by returning nil.
type Middleware func(http.ResponseWriter, *http.Request) *http.Request

// ExpectQueryParam returns a middleware that checks for the presence of a query parameter.
func ExpectQueryParam(name string) Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		if r.URL.Query().Has(name) {
			return r
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing query param: %s", name)
		return nil
	}
}