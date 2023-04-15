package middleware

import (
	"fmt"
	"io"
	"net/http"
)

// Middleware runs on any incoming request. Attachment behavior is defined by the structure it's attached to (route vs. router).
//
// Returns an *http.Request; the middleware can set router params or reject a request by returning nil.
type Middleware func(http.ResponseWriter, *http.Request) *http.Request

// Returns a middleware that checks for the presence of a query parameter.
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

// Returns a middleware that logs the details of an incoming request.
func LogRequests(w io.Writer) Middleware {
	return func(_ http.ResponseWriter, r *http.Request) *http.Request {
		logRequest(w, r)
		return r
	}
}

// Returns a middleware that logs the details of an incoming request only if
// test(request) == true.
func LogRequestsIf(test func(*http.Request) bool, w io.Writer) Middleware {
	return func(_ http.ResponseWriter, r *http.Request) *http.Request {
		if test(r) {
			logRequest(w, r)
		}
		return r
	}
}

func logRequest(w io.Writer, r *http.Request) {
	fmt.Fprintf(w, "%s %v\n", r.Method, r.URL)
}