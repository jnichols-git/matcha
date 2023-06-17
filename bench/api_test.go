package bench

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/cloudretic/matcha/pkg/cors"
	"github.com/cloudretic/matcha/pkg/middleware"
	"github.com/cloudretic/matcha/pkg/rctx"
	"github.com/cloudretic/matcha/pkg/route"
	"github.com/cloudretic/matcha/pkg/route/require"
	"github.com/cloudretic/matcha/pkg/router"
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
		AllowOrigin:  []string{"cloudretic.com"},
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
var api_rqs = []require.Required{require.Hosts("{.*}")}

var apiRoutes = []benchRoute{
	// Get/create posts
	{method: http.MethodGet, path: "/[board]/posts", testPath: "/cloudretic/posts", mws: api_mws, rqs: api_rqs},
	{method: http.MethodPut, path: "/[board]/posts", testPath: "/cloudretic/posts", mws: api_mws, rqs: api_rqs},
	{method: http.MethodGet, path: "/[board]/posts/[order]{new|top}", testPath: "/cloudretic/posts/new", mws: api_mws, rqs: api_rqs},
	// Get/update/delete post
	{method: http.MethodPatch, path: "/[board]/posts/[id]", testPath: "/cloudretic/posts/2719", mws: api_mws, rqs: api_rqs},
	{method: http.MethodDelete, path: "/[board]/posts/[id]", testPath: "/cloudretic/posts/2719", mws: api_mws, rqs: api_rqs},
	// Get/create comments
	{method: http.MethodGet, path: "/[board]/posts/[id]/comments", testPath: "/cloudretic/posts/2719/comments", mws: api_mws, rqs: api_rqs},
	{method: http.MethodPut, path: "/[board]/posts/[id]/comments", testPath: "/cloudretic/posts/2719/comments", mws: api_mws, rqs: api_rqs},
	{method: http.MethodGet, path: "/[board]/posts/[id]/comments/[order]{new|top}", testPath: "/cloudretic/posts/2719/comments/new", mws: api_mws, rqs: api_rqs},
	// Get/update/delete comment
	{method: http.MethodGet, path: "/[board]/posts/[id]/comments/[id]", testPath: "/cloudretic/posts/2719/comments/2719", mws: api_mws, rqs: api_rqs},
	{method: http.MethodPatch, path: "/[board]/posts/[id]/comments/[id]", testPath: "/cloudretic/posts/2719/comments/2719", mws: api_mws, rqs: api_rqs},
	{method: http.MethodDelete, path: "/[board]/posts/[id]/comments/[id]", testPath: "/cloudretic/posts/2719/comments/2719", mws: api_mws, rqs: api_rqs},
	// Get user info
	// Posts/comments
	{method: http.MethodGet, path: "/[user]/posts", testPath: "/jnichols/posts", mws: api_mws, rqs: api_rqs},
	{method: http.MethodGet, path: "/[user]/posts/[order]{new|top}", testPath: "/jnichols/posts/new", mws: api_mws, rqs: api_rqs},
	{method: http.MethodGet, path: "/[user]/comments", testPath: "/jnichols/comments", mws: api_mws, rqs: api_rqs},
	{method: http.MethodGet, path: "/[user]/comments/[order]{new|top}", testPath: "/jnichols/comments/new", mws: api_mws, rqs: api_rqs},
	// Liked/saved (requires client_id header)
	{method: http.MethodGet, path: "/[user]/liked", testPath: "/jnichols/liked", mws: api_mws_auth, rqs: api_rqs},
	{method: http.MethodGet, path: "/[user]/liked/[order]{new|top}", testPath: "/jnichols/liked/new", mws: api_mws_auth, rqs: api_rqs},
	{method: http.MethodGet, path: "/[user]/saved", testPath: "/jnichols/saved", mws: api_mws_auth, rqs: api_rqs},
	{method: http.MethodGet, path: "/[user]/saved/[order]{new|top}", testPath: "/jnichols/saved/new", mws: api_mws_auth, rqs: api_rqs},
}

func choose() *benchRoute {
	idx := rand.Int() % len(apiRoutes)
	return &apiRoutes[idx]
}

func handleOK(w http.ResponseWriter, req *http.Request) {
	id := req.Header.Get("X-Matcha-Request-ID")
	w.Write([]byte(id + " OK"))
}

func TestAPI(t *testing.T) {
	rt := router.Default()
	for _, tr := range apiRoutes {
		r := route.Declare(tr.method, tr.path)
		for _, mw := range tr.mws {
			r.Attach(mw)
		}
		for _, rq := range tr.rqs {
			r.Require(rq)
		}
		rt.HandleRouteFunc(r, handleOK)
	}
	w := httptest.NewRecorder()
	for i := 0; i < len(apiRoutes); i++ {
		br := apiRoutes[i]
		req := httptest.NewRequest(br.method, br.testPath, nil)
		req.Header.Set("X-Platform-User-ID", "jnichols")
		rt.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatal(br.method, br.path, br.testPath, w.Code)
		}
	}
}

func BenchmarkAPI(b *testing.B) {
	rt := router.Default()
	for _, tr := range apiRoutes {
		r, err := route.New(tr.method, tr.path)
		if err != nil {
			b.Fatal(err)
		}
		for _, mw := range tr.mws {
			r.Attach(mw)
		}
		for _, rq := range tr.rqs {
			r.Require(rq)
		}
		rt.HandleRouteFunc(r, handleOK)
	}
	w := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		br := choose()
		req := httptest.NewRequest(br.method, br.testPath, nil)
		req.Header.Set("X-Platform-User-ID", "jnichols")
		rt.ServeHTTP(w, req)
	}
}
