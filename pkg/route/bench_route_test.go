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
// 435 ns/op, 848 B/op, 9 allocs/op
func BenchmarkStringRoute(b *testing.B) {
	rt := Declare("/a/b/c/d/e/f/g/h")
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/a/b/c/d/e/f/g/h", nil)
	var out *http.Request
	for i := 0; i < b.N; i++ {
		out = rt.MatchAndUpdateContext(req)
	}
	use(out)
}

// Wildcard route
//
// 1259 ns/op, 3536 B/op, 41 allocs/op
func BenchmarkWildcardRoute(b *testing.B) {
	rt := Declare("/[a]/[b]/[c]/[d]/[e]/[f]/[g]/[h]")
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/a/b/c/d/e/f/g/h", nil)
	var out *http.Request
	for i := 0; i < b.N; i++ {
		out = rt.MatchAndUpdateContext(req)
	}
	use(out)
}

// Regex route
//
// 787.5 ns/op, 853 B/op, 9 allocs/op
func BenchmarkRegexRoute(b *testing.B) {
	rt := Declare("/{.+}/{.+}/{.+}/{.+}/{.+}/{.+}/{.+}/{.+}")
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/a/b/c/d/e/f/g/h", nil)
	var out *http.Request
	for i := 0; i < b.N; i++ {
		out = rt.MatchAndUpdateContext(req)
	}
	use(out)
}
