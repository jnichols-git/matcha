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

package cors

import (
	"net/http"
	"testing"
)

type basicResp struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

type basicWriter struct {
	Resp *basicResp
}

func (tw *basicWriter) Header() http.Header {
	return tw.Resp.Headers
}

func (tw *basicWriter) Write(b []byte) (int, error) {
	if tw.Resp.StatusCode == 0 {
		tw.WriteHeader(200)
	}
	tw.Resp.Body = append(tw.Resp.Body, b...)
	return len(b), nil
}

func (tw *basicWriter) WriteHeader(statusCode int) {
	if tw.Resp.StatusCode < 200 {
		tw.Resp.StatusCode = statusCode
	}
}

type testAdapter struct {
}

func (adapter *testAdapter) Adapt(in *http.Request) (http.ResponseWriter, *http.Request, *basicResp, error) {
	target := &basicResp{
		Headers: make(map[string][]string),
	}
	w := &basicWriter{
		Resp: target,
	}
	return w, in, target, nil
}

// using aco1, aco2, simple_request, preflight_request from cors_test.go

func TestSetCORSResponseHeaders(t *testing.T) {
	t.Run("aco1", func(t *testing.T) {
		adp := &testAdapter{}
		w, req, res, _ := adp.Adapt(simple_request)
		SetCORSResponseHeaders(w, req, aco1)
		ao := res.Headers[AllowOrigin]
		if len(ao) != 1 || ao[0] != req.Header.Get("Origin") {
			t.Errorf("expected allow-origin '%s', got %v", req.Header.Get("Origin"), ao)
		}
		am := res.Headers[AllowMethods]
		if len(am) != 1 || am[0] != req.Method {
			t.Errorf("expected allow-method '%s', got %v", req.Method, am)
		}
		ah := res.Headers[AllowHeaders]
		if len(ah) != 0 {
			t.Errorf("expected no allowed headers, got %v", ah)
		}
		eh := res.Headers[ExposeHeaders]
		if len(eh) != 1 || eh[0] != "*" {
			t.Errorf("expected all headers exposed, got %v", eh)
		}

		w, req, res, _ = adp.Adapt(preflight_request)
		SetCORSResponseHeaders(w, req, aco1)
		ao = res.Headers[AllowOrigin]
		if len(ao) != 1 || ao[0] != req.Header.Get("Origin") {
			t.Errorf("expected allow-origin '%s', got %v", req.Header.Get("Origin"), ao)
		}
		am = res.Headers[AllowMethods]
		if len(am) != 1 || am[0] != http.MethodPost {
			t.Errorf("expected allow-method '%s', got %v", http.MethodPost, am)
		}
		ah = res.Headers[AllowHeaders]
		if len(ah) != 2 || ah[0] != "X-Header-1" || ah[1] != "x-Header-2" {
			t.Errorf("expected allowed headers 'X-Header-1,x-Header-2', got %v", ah)
		}
		eh = res.Headers[ExposeHeaders]
		if len(eh) != 1 || eh[0] != "*" {
			t.Errorf("expected exposed headers '*', got %v", eh)
		}
	})
	t.Run("aco2", func(t *testing.T) {
		adp := &testAdapter{}
		w, req, res, _ := adp.Adapt(simple_request)
		SetCORSResponseHeaders(w, req, aco2)
		ao := res.Headers[AllowOrigin]
		if len(ao) != 1 || ao[0] != req.Header.Get("Origin") {
			t.Errorf("expected allow-origin '%s', got %v", req.Header.Get("Origin"), ao)
		}
		am := res.Headers[AllowMethods]
		if len(am) != 1 || am[0] != req.Method {
			t.Errorf("expected allow-method '%s', got %v", req.Method, am)
		}
		ah := res.Headers[AllowHeaders]
		if len(ah) != 0 {
			t.Errorf("expected no allowed headers, got %v", ah)
		}
		eh := res.Headers[ExposeHeaders]
		if len(eh) != 1 || eh[0] != "x-Header-Out" {
			t.Errorf("expected exposed header x-Header-Out, got %v", eh)
		}

		w, req, res, _ = adp.Adapt(preflight_request)
		SetCORSResponseHeaders(w, req, aco2)
		ao = res.Headers[AllowOrigin]
		if len(ao) != 1 || ao[0] != req.Header.Get("Origin") {
			t.Errorf("expected allow-origin '%s', got %v", req.Header.Get("Origin"), ao)
		}
		am = res.Headers[AllowMethods]
		if len(am) != 1 || am[0] != http.MethodPost {
			t.Errorf("expected allow-method '%s', got %v", http.MethodPost, am)
		}
		ah = res.Headers[AllowHeaders]
		if len(ah) != 2 || ah[0] != "X-Header-1" || ah[1] != "x-Header-2" {
			t.Errorf("expected allowed headers X-Header-1,x-Header-2, got %v", ah)
		}
		eh = res.Headers[ExposeHeaders]
		if len(eh) != 1 || eh[0] != "x-Header-Out" {
			t.Errorf("expected exposed headers x-Header-Out, got %v", eh)
		}
	})
}

func BenchmarkSetCORSResponseHeaders(b *testing.B) {
	adp := &testAdapter{}
	w, req, _, _ := adp.Adapt(simple_request)
	for i := 0; i < b.N; i++ {
		SetCORSResponseHeaders(w, req, aco2)
	}
}
