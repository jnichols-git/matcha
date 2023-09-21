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

package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

func noop(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func getID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if id, err := strconv.Atoi(r.Header.Get("X-Request-ID")); err == nil {
			r = r.WithContext(context.WithValue(r.Context(), "Request-ID", id))
		}
		next.ServeHTTP(w, r)
	})
}

func TestHandlerConcurrent(t *testing.T) {
	mw1 := Handler(noop)
	mw2 := Handler(getID)
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		a := i
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("X-Request-ID", strconv.FormatInt(int64(a), 10))
			req = mw1(nil, req)
			req = mw2(nil, req)
			if req == nil {
				t.Errorf("nil request")
			} else if id := req.Context().Value("Request-ID"); id != a {
				t.Error("expected", a, "got", id)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkHandlerNoop(b *testing.B) {
	mw := Handler(getID)
	wg := &sync.WaitGroup{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1; j++ {
			wg.Add(1)
			go func() {
				req = mw(nil, req)
				wg.Done()
			}()
			wg.Wait()
		}
	}
}

func BenchmarkHandlerID(b *testing.B) {
	mw := Handler(getID)
	wg := &sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1; j++ {
			wg.Add(1)
			a := j
			go func() {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("X-Request-ID", strconv.FormatInt(int64(a), 10))
				req = mw(nil, req)
				if req == nil {
					b.Errorf("nil request")
				} else if id := req.Context().Value("Request-ID"); id != a {
					b.Error("expected", a, "got", id)
				}
				wg.Done()
			}()
			wg.Wait()
		}
	}
}

func TestHandlerAsMiddleware(t *testing.T) {
	expectedResponse := "Hello, world!"
	create := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, expectedResponse)
			next.ServeHTTP(w, r)
		})
	}
	mw := Handler(create)
	req := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
	rec := httptest.NewRecorder()
	req = mw(rec, req)
	if req == nil {
		t.Errorf("request should be unchanged.")
	}
	if res := rec.Body.String(); res != "Hello, world!" {
		t.Errorf(
			"response should be \"%s\", got \"%s\"",
			expectedResponse,
			res,
		)
	}
}
