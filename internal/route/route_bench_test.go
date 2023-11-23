package route

import (
	"net/http"
	"testing"

	"github.com/jnichols-git/matcha/v2/pkg/rctx"
)

func use(any) {}

// Benchmarking
// Benchmarks are done on 8-length routes, where each part contains the structure being tested.

// Static route
//
// 189 ns/op, 256 B/op, 1 allocs/op
func BenchmarkStringRoute(b *testing.B) {
	rt := Declare(http.MethodGet, "/a/b/c/d/e/f/g/h")
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/a/b/c/d/e/f/g/h", nil)
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	var out *http.Request
	for i := 0; i < b.N; i++ {
		out = rt.MatchAndUpdateContext(req)
	}
	use(out)
}

// Wildcard route
//
// 407 ns/op, 256 B/op, 1 allocs/op
func BenchmarkWildcardRoute(b *testing.B) {
	rt := Declare(http.MethodGet, "/{a}/{b}/{c}/{d}/{e}/{f}/{g}/{h}")
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/a/b/c/d/e/f/g/h", nil)
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	var out *http.Request
	for i := 0; i < b.N; i++ {
		out = rt.MatchAndUpdateContext(req)
	}
	use(out)
}

// Regex route
//
// 476 ns/op, 257 B/op, 1 allocs/op
func BenchmarkRegexRoute(b *testing.B) {
	rt := Declare(http.MethodGet, "/{.+}/{.+}/{.+}/{.+}/{.+}/{.+}/{.+}/{.+}")
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/a/b/c/d/e/f/g/h", nil)
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	var out *http.Request
	for i := 0; i < b.N; i++ {
		out = rt.MatchAndUpdateContext(req)
	}
	use(out)
}

// Partial route
//
// 533 ns/op, 257 B/op, 1 allocs/op
func BenchmarkPartialRoute(b *testing.B) {
	rt := Declare(http.MethodGet, "/+")
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/a/b/c/d/e/f/g/h", nil)
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	var out *http.Request
	for i := 0; i < b.N; i++ {
		out = rt.MatchAndUpdateContext(req)
	}
	use(out)
}
