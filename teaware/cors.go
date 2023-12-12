package teaware

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2/cors"
)

func Options(aco *cors.Options) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cors.SetCORSResponseHeaders(w, r, aco)
			next.ServeHTTP(w, r)
		})
	}
}
