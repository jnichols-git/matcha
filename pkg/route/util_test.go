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

package route

import (
	"net/http"
	"testing"
)

func TestNumParams(t *testing.T) {
	r1 := Declare(http.MethodGet, "/static/route")
	if np := NumParams(r1); np != 0 {
		t.Errorf("expected 0 params, got %d", np)
	}
	r2 := Declare(http.MethodGet, "/[wc1]{.+}/[wc2]")
	if np := NumParams(r2); np != 2 {
		t.Errorf("expected 2 params, got %d", np)
	}
	r3 := Declare(http.MethodGet, "/[wc1]{.+}/[wc2]/[wc3]+")
	if np := NumParams(r3); np != 3 {
		t.Errorf("expected 3 params, got %d", np)
	}
	r4 := Declare(http.MethodGet, "/{.+}/[wc2]/+")
	if np := NumParams(r4); np != 1 {
		t.Errorf("expected 1 params, got %d", np)
	}
}
