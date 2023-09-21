// Copyright 2023 Decent Platforms
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

package tree

import (
	"math/rand"
	"net/http"
	"testing"

	"github.com/decentplatforms/matcha/pkg/route"
)

var routes []string = []string{
	"/1/1/1", "/1/1/2", "/1/1/3",
	"/1/2/1", "/1/2/2", "/1/2/3",
	"/1/3/1", "/1/3/2", "/1/3/3",
	"/2/1/1", "/2/1/2", "/2/1/3",
	"/2/2/1", "/2/2/2", "/2/2/3",
	"/2/3/1", "/2/3/2", "/2/3/3",
	"/3/1/1", "/3/1/2", "/3/1/3",
	"/3/2/1", "/3/2/2", "/3/2/3",
	"/3/3/1", "/3/3/2", "/3/3/3",
}

func testTree() *RouteTree {
	rtree := New()
	for _, rs := range routes {
		r := route.Declare(http.MethodGet, rs)
		rtree.Add(r)
	}
	return rtree
}

var tt *RouteTree = testTree()

func BenchmarkTreeMatchSingle(b *testing.B) {
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	for i := 0; i < b.N; i++ {
		ri := rand.Int() % len(routes)
		req.URL.Path = routes[ri]
		tt.Match(req)
	}
}
func BenchmarkTreeMatchAll(b *testing.B) {
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	for i := 0; i < b.N; i++ {
		for _, r := range routes {
			req.URL.Path = r
			tt.Match(req)
		}
	}
}
