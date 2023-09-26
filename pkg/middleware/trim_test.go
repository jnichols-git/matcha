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

package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrimPrefix(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/path/to/resource", nil)

	reg := TrimPrefix("/path/to")
	r1 := reg(w, req.Clone(context.Background()))
	if r1.URL.Path != "/resource" {
		t.Error("/resource", r1.URL.Path)
	}
	reg_wrong := TrimPrefix("/other/path")
	r2 := reg_wrong(w, req.Clone(context.Background()))
	if r2.URL.Path != "/path/to/resource" {
		t.Error("/path/to/resource", r1.URL.Path)
	}

	strict := TrimPrefixStrict("/path/to", "")
	r3 := strict(w, req.Clone(context.Background()))
	if r3.URL.Path != "/resource" {
		t.Error("/resource", r1.URL.Path)
	}
	strict_wrong := TrimPrefixStrict("/other/path", "")
	r4 := strict_wrong(w, req.Clone(context.Background()))
	if r4 != nil {
		t.Error("expected nil request")
	}
	if w.Code != http.StatusBadRequest {
		t.Error(http.StatusBadRequest, w.Code)
	}
}
