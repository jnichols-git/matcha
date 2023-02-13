package route

import (
	"net/http"
	"testing"
)

func use(any) {}

// Benchmarking
// Benchmarks are done on 8-length routes, where each part contains the structure being tested.

// Static route
//
// 185 ns/op, 256 B/op, 1 allocs/op
func BenchmarkStringRoute(b *testing.B) {
	rt := Declare(http.MethodGet, "/a/b/c/d/e/f/g/h")
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/a/b/c/d/e/f/g/h", nil)
	var out *http.Request
	for i := 0; i < b.N; i++ {
		out = rt.MatchAndUpdateContext(req)
	}
	use(out)
}

// Wildcard route
//
// 396 ns/op, 256 B/op, 1 allocs/op
func BenchmarkWildcardRoute(b *testing.B) {
	rt := Declare(http.MethodGet, "/[a]/[b]/[c]/[d]/[e]/[f]/[g]/[h]")
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/a/b/c/d/e/f/g/h", nil)
	var out *http.Request
	for i := 0; i < b.N; i++ {
		out = rt.MatchAndUpdateContext(req)
	}
	use(out)
}

// Regex route
//
// 469 ns/op, 257 B/op, 1 allocs/op
func BenchmarkRegexRoute(b *testing.B) {
	rt := Declare(http.MethodGet, "/{.+}/{.+}/{.+}/{.+}/{.+}/{.+}/{.+}/{.+}")
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/a/b/c/d/e/f/g/h", nil)
	var out *http.Request
	for i := 0; i < b.N; i++ {
		out = rt.MatchAndUpdateContext(req)
	}
	use(out)
}

// Partial route
//
// 520 ns/op, 257 B/op, 1 allocs/op
func BenchmarkPartialRoute(b *testing.B) {
	rt := Declare(http.MethodGet, "/+")
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/a/b/c/d/e/f/g/h", nil)
	var out *http.Request
	for i := 0; i < b.N; i++ {
		out = rt.MatchAndUpdateContext(req)
	}
	use(out)
}
