package router

import (
	"net/http"
	"testing"

	"github.com/cloudretic/router/pkg/route"
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

// Basic router benchmark.
// For more involved benchmarks, see /bench. This serves as a baseline value, not a robust example under load.
func BenchmarkBasicRouter(b *testing.B) {
	rt := Declare(Default(),
		WithRoute(route.Declare(http.MethodGet, "/"), okHandler("root")),
		WithRoute(route.Declare(http.MethodGet, "/[wildcard]"), rpHandler("wildcard")),
		WithRoute(route.Declare(http.MethodGet, `/route/{[a-zA-Z]+}`), okHandler("letters")),
		WithRoute(route.Declare(http.MethodGet, `/route/[id]{[\w]{4}}`), rpHandler("id")),
		WithRoute(route.Declare(http.MethodGet, `/static/file/[filename]{\w+(?:\.\w+)?}+`), rpHandler("filename")),
	)
	benchReqs := []*http.Request{
		declareReq("/"),
		declareReq("/wc"),
		declareReq("/route/aWord"),
		declareReq("/route/anID"),
		declareReq("/static/file/some/file/path.md"),
	}
	mockWriter := &mockResponseWriter{}
	for i := 0; i < b.N; i++ {
		for _, r := range benchReqs {
			rt.ServeHTTP(mockWriter, r)
		}
	}
}
