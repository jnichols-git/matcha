package bench

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/jnichols-git/matcha/v2/cors"
	"github.com/jnichols-git/matcha/v2/internal/rctx"
	"github.com/jnichols-git/matcha/v2/require"
	"github.com/jnichols-git/matcha/v2/teaware"
)

type benchRoute struct {
	method   string
	path     string
	testPath string
	mws      []teaware.Middleware
	rqs      []require.Required
}

var aco = &cors.Options{
	AllowOrigin:  []string{"jnichols.info"},
	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodDelete},
	AllowHeaders: []string{"client_id"},
}

func mwID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("X-Matcha-Request-ID", strconv.FormatInt(rand.Int63(), 10))
		next.ServeHTTP(w, r)
	})
}

func mwIsUserParam(userParam string) teaware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := rctx.GetParam(r.Context(), userParam)
			is := r.Header.Get("X-Platform-User-ID")
			if is != user {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("user " + is + " unauthorized"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

var api_mws = []teaware.Middleware{teaware.Options(aco), mwID}
var api_mws_auth = []teaware.Middleware{mwIsUserParam("user"), teaware.Options(aco), mwID}
var api_rqs = []require.Required{require.Hosts("[.*]")}
