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

package rctx

import (
	"net/http"
	"net/url"
	"testing"
)

func use(any) {}

func BenchmarkNewParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		params := newParams(DefaultMaxParams)
		use(params)
	}
}

func BenchmarkPrepare(b *testing.B) {
	url, _ := url.Parse("/")
	req := &http.Request{
		URL: url,
	}
	for i := 0; i < b.N; i++ {
		req = PrepareRequestContext(req, DefaultMaxParams)
		//req = req.WithContext(context.Background())
		// req = req.WithContext(context.Background())
		req = &http.Request{
			URL: url,
		}
	}
}

func BenchmarkSetGetSingleParam(b *testing.B) {
	req := &http.Request{}
	req = PrepareRequestContext(req, DefaultMaxParams)
	for i := 0; i < b.N; i++ {
		SetParam(req.Context(), "paramKey", "paramVal")
		GetParam(req.Context(), "paramKey")
		ResetRequestContext(req)
	}
}
