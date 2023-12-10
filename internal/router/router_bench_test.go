package router

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/jnichols-git/matcha/v2/internal/route"
)

// mock response writer, taken from go-http-routing-benchmark
type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func declareReq(path string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	req.URL.Path = path
	req.RequestURI = path
	return req
}

// Benchmark single requests. Results should generally average out to the cost of a single handle after b.N runs.
func BenchmarkSingleRequests(b *testing.B) {
	rt := Default()
	rt.HandleRoute(route.Declare(http.MethodGet, "/"), okHandler("root"))
	rt.HandleRoute(route.Declare(http.MethodGet, "/:wildcard"), rpHandler("wildcard"))
	rt.HandleRoute(route.Declare(http.MethodGet, `/route/{[a-zA-Z]+}`), okHandler("letters"))
	rt.HandleRoute(route.Declare(http.MethodGet, `/route/:id{[\w]{4}}`), rpHandler("id"))
	rt.HandleRoute(route.Declare(http.MethodGet, `/static/file/:filename{\w+(?:\.\w+)?}+`), rpHandler("filename"))
	benchReqs := []*http.Request{
		declareReq("/"),
		declareReq("/wc"),
		declareReq("/route/aWord"),
		declareReq("/route/anID"),
		declareReq("/static/file/some/file/path.md"),
	}
	mockWriter := &mockResponseWriter{}
	handler := rt.Compile()
	for i := 0; i < b.N; i++ {
		ri := rand.Int() % len(benchReqs)
		r := benchReqs[ri]
		handler.ServeHTTP(mockWriter, r)
	}
}

// Basic router benchmark.
// For more involved benchmarks, see /bench. This serves as a baseline value, not a robust example under load.
func BenchmarkBasicRouter(b *testing.B) {
	rt := Default()
	rt.HandleRoute(route.Declare(http.MethodGet, "/"), okHandler("root"))
	rt.HandleRoute(route.Declare(http.MethodGet, "/:wildcard"), rpHandler("wildcard"))
	rt.HandleRoute(route.Declare(http.MethodGet, `/route/{[a-zA-Z]+}`), okHandler("letters"))
	rt.HandleRoute(route.Declare(http.MethodGet, `/route/:id{[\w]{4}}`), rpHandler("id"))
	rt.HandleRoute(route.Declare(http.MethodGet, `/static/file/:filename{\w+(?:\.\w+)?}+`), rpHandler("filename"))
	benchReqs := []*http.Request{
		declareReq("/"),
		declareReq("/wc"),
		declareReq("/route/aWord"),
		declareReq("/route/anID"),
		declareReq("/static/file/some/file/path.md"),
	}
	mockWriter := &mockResponseWriter{}
	handler := rt.Compile()
	for i := 0; i < b.N; i++ {
		for _, r := range benchReqs {
			handler.ServeHTTP(mockWriter, r)
		}
	}
}
