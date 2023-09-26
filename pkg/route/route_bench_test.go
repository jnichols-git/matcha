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

package route

import (
	"net/http"
	"testing"

	"github.com/decentplatforms/matcha/pkg/rctx"
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
	rt := Declare(http.MethodGet, "/[a]/[b]/[c]/[d]/[e]/[f]/[g]/[h]")
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
