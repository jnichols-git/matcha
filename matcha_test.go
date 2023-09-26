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

package matcha

/*
 * The central package exists to alias common commands and reduce imports; all this behavior is tested elsewhere.
 * The tests here are mostly a formality to maintain coverage.
 */

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/decentplatforms/matcha/pkg/rctx"
)

func TestRouter(t *testing.T) {
	rt := Router()
	if rt == nil {
		t.Fatal("nil router")
	}
}

func TestRoute(t *testing.T) {
	r := Route(http.MethodGet, "/test/route")
	if r == nil {
		t.Fatal("nil route")
	}
}

func TestGetParam(t *testing.T) {
	r := Route(http.MethodGet, "/test/[param]")
	req := httptest.NewRequest(http.MethodGet, "/test/value", nil)
	req = rctx.PrepareRequestContext(req, 1)
	r.MatchAndUpdateContext(req)
	if value := GetParam(req.Context(), "param"); value != "value" {
		t.Error(value)
	}
}
