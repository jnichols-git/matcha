package bench

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/jnichols-git/matcha/v2/internal/rctx"
	"github.com/jnichols-git/matcha/v2/pkg/cors"
	"github.com/jnichols-git/matcha/v2/pkg/middleware"
	"github.com/jnichols-git/matcha/v2/pkg/require"
)

type benchRoute struct {
	method   string
	path     string
	testPath string
	mws      []middleware.Middleware
	rqs      []require.Required
}

func mwCORS() middleware.Middleware {
	return cors.CORSMiddleware(&cors.AccessControlOptions{
		AllowOrigin:  []string{"jnichols.info"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{"client_id"},
	})
}

func mwID(w http.ResponseWriter, req *http.Request) *http.Request {
	req.Header.Add("X-Matcha-Request-ID", strconv.FormatInt(rand.Int63(), 10))
	return req
}

func mwIsUserParam(userParam string) middleware.Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		user := rctx.GetParam(r.Context(), userParam)
		is := r.Header.Get("X-Platform-User-ID")
		if is != user {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("user " + is + " unauthorized"))
			return nil
		}
		return r
	}
}

var api_mws = []middleware.Middleware{mwCORS(), mwID}
var api_mws_auth = []middleware.Middleware{mwIsUserParam("user"), mwCORS(), mwID}
var api_rqs = []require.Required{require.Hosts("[.*]")}
