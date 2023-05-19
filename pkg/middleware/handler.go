package middleware

import (
	"net/http"
)

// Handler allows integration of traditional handler-chain-ware.
func Handler(create func(next http.Handler) http.Handler) Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		var out *http.Request
		next := func(w http.ResponseWriter, req *http.Request) {
			out = req
		}
		h := create(http.HandlerFunc(next))
		h.ServeHTTP(w, r)
		return out
	}
}
