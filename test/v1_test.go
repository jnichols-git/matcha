package bench

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/jnichols-git/matcha/v2"
	"github.com/jnichols-git/matcha/v2/route"
)

/*
 * ===== MockBoards V1 API =====
 *
 * Use Case: MockBoards is a (fake) forum website setting up their API using Matcha. The V1 API has:
 *
 * - 18 endpoints
 * - 4 endpoints requiring client_id header authorization
 * - 4 endpoints ending in an API parameter with an enumeration (new/top)
 * - Middleware assigning request IDs and CORS headers
 * - Requirements for host matching on all routes
 *
 * This benchmark approximates the performance of Matcha while routing their API.
 *
 * ===== Using This Benchmark =====
 *
 * Run all 3 included benchmarks (API, stripped, offset)
 * For sequential results, subtract the offset from the X/op values.
 * For concurrent results, divide the X/op values by 10, then subtract the offset.
 */

// MockBoards API
var apiRoutes = []benchRoute{
	// Get/create posts
	{method: http.MethodGet, path: "/:board/posts", testPath: "/jnichols/posts", mws: api_mws, rqs: api_rqs},
	{method: http.MethodPut, path: "/:board/posts", testPath: "/jnichols/posts", mws: api_mws, rqs: api_rqs},
	{method: http.MethodGet, path: "/:board/posts/:order[new|top]", testPath: "/jnichols/posts/new", mws: api_mws, rqs: api_rqs},
	// Get/update/delete post
	{method: http.MethodPatch, path: "/:board/posts/:id", testPath: "/jnichols/posts/2719", mws: api_mws, rqs: api_rqs},
	{method: http.MethodDelete, path: "/:board/posts/:id", testPath: "/jnichols/posts/2719", mws: api_mws, rqs: api_rqs},
	// Get/create comments
	{method: http.MethodGet, path: "/:board/posts/:id/comments", testPath: "/jnichols/posts/2719/comments", mws: api_mws, rqs: api_rqs},
	{method: http.MethodPut, path: "/:board/posts/:id/comments", testPath: "/jnichols/posts/2719/comments", mws: api_mws, rqs: api_rqs},
	{method: http.MethodGet, path: "/:board/posts/:id/comments/:order[new|top]", testPath: "/jnichols/posts/2719/comments/new", mws: api_mws, rqs: api_rqs},
	// Get/update/delete comment
	{method: http.MethodGet, path: "/:board/posts/:id/comments/:id", testPath: "/jnichols/posts/2719/comments/2719", mws: api_mws, rqs: api_rqs},
	{method: http.MethodPatch, path: "/:board/posts/:id/comments/:id", testPath: "/jnichols/posts/2719/comments/2719", mws: api_mws, rqs: api_rqs},
	{method: http.MethodDelete, path: "/:board/posts/:id/comments/:id", testPath: "/jnichols/posts/2719/comments/2719", mws: api_mws, rqs: api_rqs},
	// Get user info
	// Posts/comments
	{method: http.MethodGet, path: "/:user/posts", testPath: "/jnichols/posts", mws: api_mws, rqs: api_rqs},
	{method: http.MethodGet, path: "/:user/posts/:order[new|top]", testPath: "/jnichols/posts/new", mws: api_mws, rqs: api_rqs},
	{method: http.MethodGet, path: "/:user/comments", testPath: "/jnichols/comments", mws: api_mws, rqs: api_rqs},
	{method: http.MethodGet, path: "/:user/comments/:order[new|top]", testPath: "/jnichols/comments/new", mws: api_mws, rqs: api_rqs},
	// Liked/saved (requires client_id header)
	{method: http.MethodGet, path: "/:user/liked", testPath: "/jnichols/liked", mws: api_mws_auth, rqs: api_rqs},
	{method: http.MethodGet, path: "/:user/liked/:order[new|top]", testPath: "/jnichols/liked/new", mws: api_mws_auth, rqs: api_rqs},
	{method: http.MethodGet, path: "/:user/saved", testPath: "/jnichols/saved", mws: api_mws_auth, rqs: api_rqs},
	{method: http.MethodGet, path: "/:user/saved/:order[new|top]", testPath: "/jnichols/saved/new", mws: api_mws_auth, rqs: api_rqs},
}

func choose() *benchRoute {
	idx := rand.Int() % len(apiRoutes)
	return &apiRoutes[idx]
}

func handleOK(w http.ResponseWriter, req *http.Request) {
	id := req.Header.Get("X-Matcha-Request-ID")
	w.Write([]byte(id + " OK"))
}

// Just to check!
func TestAPIv1(t *testing.T) {
	rt := matcha.Router()
	for _, tr := range apiRoutes {
		r, _ := matcha.Route(tr.method, tr.path)
		r.Use(tr.mws...)
		r.Require(tr.rqs...)
		rt.HandleRouteFunc(r, handleOK)
	}
	h := rt
	w := httptest.NewRecorder()
	for i := 0; i < len(apiRoutes); i++ {
		br := apiRoutes[i]
		req := httptest.NewRequest(br.method, br.testPath, nil)
		req.Header.Set("X-Platform-User-ID", "jnichols")
		h.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Fatal(br.method, br.path, br.testPath, w.Code)
		}
	}
}

// BENCHMARKS: MockBoards API

func BenchmarkAPIv1(b *testing.B) {
	rt := matcha.Router()
	for _, tr := range apiRoutes {
		r, err := route.New(tr.method, tr.path)
		if err != nil {
			b.Fatal(err)
		}
		r.Use(tr.mws...)
		r.Require(tr.rqs...)
		rt.HandleRouteFunc(r, handleOK)
	}
	h := rt
	b.Run(b.Name()+"-sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			w := httptest.NewRecorder()
			br := choose()
			req := httptest.NewRequest(br.method, br.testPath, nil)
			req.Header.Set("X-Platform-User-ID", "jnichols")
			h.ServeHTTP(w, req)
		}
	})
	b.Run(b.Name()+"-concurrent-10", func(b *testing.B) {
		wg := &sync.WaitGroup{}
		for i := 0; i < b.N; i++ {
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					w := httptest.NewRecorder()
					br := choose()
					req := httptest.NewRequest(br.method, br.testPath, nil)
					req.Header.Set("X-Platform-User-ID", "jnichols")
					h.ServeHTTP(w, req)
					wg.Done()
				}()
			}
			wg.Wait()
		}
	})
}

func BenchmarkStrippedAPI(b *testing.B) {
	rt := matcha.Router()
	for _, tr := range apiRoutes {
		r, err := route.New(tr.method, tr.path)
		if err != nil {
			b.Fatal(err)
		}
		rt.HandleRouteFunc(r, handleOK)
	}
	b.Run(b.Name()+"-sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			w := httptest.NewRecorder()
			br := choose()
			req := httptest.NewRequest(br.method, br.testPath, nil)
			req.Header.Set("X-Platform-User-ID", "jnichols")
			rt.ServeHTTP(w, req)
			if w.Code != 200 {
				b.Fatal(w.Code)
			}
		}
	})
	b.Run(b.Name()+"-concurrent-10", func(b *testing.B) {
		wg := &sync.WaitGroup{}
		for i := 0; i < b.N; i++ {
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					w := httptest.NewRecorder()
					br := choose()
					req := httptest.NewRequest(br.method, br.testPath, nil)
					req.Header.Set("X-Platform-User-ID", "jnichols")
					rt.ServeHTTP(w, req)
					wg.Done()
				}()
			}
			wg.Wait()
		}
	})
}

// BenchmarkOffset is used to calculate the performance offset of generating test values
// for benchmarks.
// This is important; creating new requests is expensive
func BenchmarkOffset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = httptest.NewRecorder()
		br := choose()
		req := httptest.NewRequest(br.method, br.testPath, nil)
		req.Header.Set("X-Platform-User-ID", "jnichols")
	}
}
