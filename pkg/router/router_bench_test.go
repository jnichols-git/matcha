// Copyright 2023 Matcha Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/decentplatforms/matcha/pkg/route"
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
	rt := Declare(Default(),
		HandleRoute(route.Declare(http.MethodGet, "/"), okHandler("root")),
		HandleRoute(route.Declare(http.MethodGet, "/[wildcard]"), rpHandler("wildcard")),
		HandleRoute(route.Declare(http.MethodGet, `/route/{[a-zA-Z]+}`), okHandler("letters")),
		HandleRoute(route.Declare(http.MethodGet, `/route/[id]{[\w]{4}}`), rpHandler("id")),
		HandleRoute(route.Declare(http.MethodGet, `/static/file/[filename]{\w+(?:\.\w+)?}+`), rpHandler("filename")),
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
		ri := rand.Int() % len(benchReqs)
		r := benchReqs[ri]
		rt.ServeHTTP(mockWriter, r)
	}
}

// Basic router benchmark.
// For more involved benchmarks, see /bench. This serves as a baseline value, not a robust example under load.
func BenchmarkBasicRouter(b *testing.B) {
	rt := Declare(Default(),
		HandleRoute(route.Declare(http.MethodGet, "/"), okHandler("root")),
		HandleRoute(route.Declare(http.MethodGet, "/[wildcard]"), rpHandler("wildcard")),
		HandleRoute(route.Declare(http.MethodGet, `/route/{[a-zA-Z]+}`), okHandler("letters")),
		HandleRoute(route.Declare(http.MethodGet, `/route/[id]{[\w]{4}}`), rpHandler("id")),
		HandleRoute(route.Declare(http.MethodGet, `/static/file/[filename]{\w+(?:\.\w+)?}+`), rpHandler("filename")),
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
