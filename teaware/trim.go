package teaware

import (
	"net/http"
	"strings"
)

func TrimPrefix(prefix string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rest, ok := strings.CutPrefix(r.URL.Path, prefix); ok {
				r.URL.Path = rest
			}
			next.ServeHTTP(w, r)
		})
	}
}
