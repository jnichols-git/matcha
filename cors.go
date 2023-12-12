package matcha

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2/cors"
)

// SetCORSHeaders sets CORS headers on the response according to your
// cors.Options.
// This function is for managing CORS on "simple" requests that don't use
// preflight. You may also want to try teaware.Options. If you want to handle
// preflight/OPTIONS requests, use Options to create a handler for it.
func SetCORSHeaders(w http.ResponseWriter, req *http.Request, aco *cors.Options) {
	cors.SetCORSResponseHeaders(w, req, aco)
}

// Options returns an http.Handler that sets CORS headers according to your
// cors.Options.
func Options(aco *cors.Options) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cors.SetCORSResponseHeaders(w, r, aco)
	})
}
