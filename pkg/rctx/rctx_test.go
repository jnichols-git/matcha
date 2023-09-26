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
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestPrepareRequestContext(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://test.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = PrepareRequestContext(req, DefaultMaxParams)
	ctx := req.Context()
	if _, ok := ctx.(*Context); !ok {
		t.Error("should be of type *rctx.Context")
	}
	if req.Method != http.MethodPost {
		t.Error("method shouldn't change")
	}
}

func TestNative(t *testing.T) {
	t.Run("with rctx", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "http://test.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = PrepareRequestContext(req, 4)
		ctx := req.Context()
		for i := 0; i < 4; i++ {
			p := fmt.Sprintf("p%d", i)
			v := fmt.Sprintf("v%d", i)
			err := SetParam(ctx, p, v)
			if err != nil {
				t.Error(err)
			}
		}
		err = SetParam(ctx, "errp", "errv")
		if err == nil {
			t.Error("should fail to set more than 4 params")
		}
		for i := 0; i < 4; i++ {
			p := fmt.Sprintf("p%d", i)
			v := fmt.Sprintf("v%d", i)
			got := GetParam(ctx, p)
			if got != v {
				t.Errorf("expected %s, got %s", v, got)
			}
		}
		SetParam(ctx, "p0", "diffValue")
		if v := GetParam(ctx, "p0"); v != "diffValue" {
			t.Errorf("expected diffValue, got %s", v)
		}
		ResetRequestContext(req)
		got := GetParam(ctx, "p0")
		if got != "" {
			t.Error("should be empty")
		}
		err = SetParam(ctx, "new-param", "new-value")
		if err != nil {
			t.Error(err)
		} else if val := GetParam(ctx, "new-param"); val != "new-value" {
			t.Errorf("expected new-value, got %s", val)
		}
		ReturnRequestContext(req)
		req, err = http.NewRequest(http.MethodPost, "http://test.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = PrepareRequestContext(req, 3)
		ctx = req.Context()
		if l := ctx.(*Context).params.cap; l != 3 {
			t.Errorf("should have space for 3 params exactly")
		}
		for _, kv := range ctx.(*Context).params.rps {
			if kv.key != "" || kv.value != "" {
				t.Errorf("values should be cleared from pool contexts")
			}
		}
	})
	t.Run("with ctx", func(t *testing.T) {})
}

// Test native functions on a regular context.Context.
// This should generally fail; use of context.Context is wildly inefficient and shouldn't be supported by most native functions.
// See below for the context.Context implementations, which are supported.
func TestNativeWithHttp(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://test.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = SetParam(req.Context(), "testParam", "testValue")
	if err == nil {
		t.Error("expected set to fail on default request context")
	}
	val := GetParam(req.Context(), "testParam")
	if val != "" {
		t.Error("expected get to be empty")
	}
	// GetParam should work on non-*rctx.Context contexts to avoid issues with changing context in handlers.
	ctx := context.WithValue(req.Context(), paramKey("testParam"), "testValue")
	val = GetParam(ctx, "testParam")
	if val == "" {
		t.Error("get should work on a regular context.Context")
	}
	// Reset without preparing context; should fail, as default context doesn't work here
	err = ResetRequestContext(req)
	if err == nil {
		t.Error("expected reset to fail on default request context")
	}
}

func TestImplementContext(t *testing.T) {
	// Test passthrough
	req, err := http.NewRequest(http.MethodPost, "http://test.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	dl := time.Now().Add(time.Second * 1)
	ctx := req.Context()
	ctx, cancel := context.WithDeadline(ctx, dl)
	ctx = context.WithValue(ctx, "testKey", "testValue")
	req = req.WithContext(ctx)
	req = PrepareRequestContext(req, DefaultMaxParams)
	if ctxdl, _ := req.Context().Deadline(); !ctxdl.Equal(dl) {
		t.Errorf("expected equal deadlines, got %v and %v", ctxdl, dl)
	}
	if err := req.Context().Err(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if v := req.Context().Value("testKey"); v != "testValue" {
		t.Errorf("expected testValue, got %s", v)
	}
	time.Sleep(time.Second)
	if done := req.Context().Done(); done == nil {
		t.Error("done should not be nil")
	} else {
		stop := make(chan struct{})
		mtx := &sync.Mutex{}
		go func() {
			time.Sleep(time.Second)
			stop <- struct{}{}
		}()
		mtx.Lock()
		go func() {
		wait:
			for {
				select {
				case <-stop:
					t.Error("done should trigger before the stop signal here")
					break wait
				case <-done:
					break wait
				}
			}
			mtx.Unlock()
		}()
		mtx.Lock()
		mtx.Unlock()
	}
	cancel()
	// Test nil parent
	req = &http.Request{}
	req = PrepareRequestContext(req, DefaultMaxParams)
	if ctx, ok := req.Context().(*Context); !ok {
		t.Fatal("need *rctx.Context")
	} else {
		ctx.parent = nil
	}
	if ctxdl, _ := req.Context().Deadline(); !ctxdl.IsZero() {
		t.Errorf("expected zero deadline, got %v", ctxdl)
	}
	if err := req.Context().Err(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if v := req.Context().Value("testKey"); v != nil {
		t.Errorf("expected nil, got %s", v)
	}
	if done := req.Context().Done(); done != nil {
		t.Errorf("expected nil")
	}
}

// Test to make sure values persist when context is overridden
func TestOverrideContext(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://test.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = PrepareRequestContext(req, DefaultMaxParams)
	if ctx, ok := req.Context().(*Context); !ok {
		t.Fatal("need *rctx.Context")
	} else {
		ctx.err = errors.New("test error")
		SetParam(ctx, "testParam", "testCorrectValue")
		req = req.WithContext(context.WithValue(ctx, "testParam", "testIncorrectValue"))
	}
	ctx := req.Context()
	if err := ctx.Err(); err == nil {
		t.Error("expected error to pass through")
	}
	if v := ctx.Value(paramKey("testParam")); v != "testCorrectValue" {
		t.Errorf("got %s", v)
	}
}

func TestNested(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/path", nil)
	ctx := context.WithValue(req.Context(), paramKey("p3"), "v3")
	req = req.WithContext(ctx)

	req = PrepareRequestContext(req, 1)
	err := SetParam(req.Context(), "p1", "v1")
	if err != nil {
		t.Fatal(err)
	}
	req = PrepareRequestContext(req, 1)
	err = SetParam(req.Context(), "p2", "v2")
	if err != nil {
		t.Fatal(err)
	}

	if p1 := GetParam(req.Context(), "p1"); p1 != "v1" {
		t.Error("v1", p1)
	}
	if p2 := GetParam(req.Context(), "p2"); p2 != "v2" {
		t.Error("v2", p2)
	}
	if p3 := GetParam(req.Context(), "p3"); p3 != "v3" {
		t.Error("v3", p3)
	}
	if p4 := GetParam(req.Context(), "p4"); p4 != "" {
		t.Error("", p4)
	}

	// Test nil parents
	rctx := req.Context().(*Context)
	rctx.parent.(*Context).parent = nil
	if p3 := GetParam(req.Context(), "p3"); p3 != "" {
		t.Error("", p3)
	}

	rctx.parent = nil
	if p1 := GetParam(req.Context(), "p1"); p1 != "" {
		t.Error("", p1)
	}
}

func TestFullPath(t *testing.T) {
	// Regular case
	req := httptest.NewRequest(http.MethodGet, "http://test.com/test/path", nil)
	req = PrepareRequestContext(req, 0)
	if fp := GetParam(req.Context(), PARAM_FULLPATH); fp != "/test/path" {
		t.Error(fp)
	}
	// Empty
	req = httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	req = PrepareRequestContext(req, 0)
	if fp := GetParam(req.Context(), PARAM_FULLPATH); fp != "" {
		t.Error(fp)
	}
	// Nil; should pass and fullpath should be empty
	req = &http.Request{}
	req = PrepareRequestContext(req, 0)
	if fp := GetParam(req.Context(), PARAM_FULLPATH); fp != "" {
		t.Error(fp)
	}
	// Nested; should be equal to fullpath of first context
	req = httptest.NewRequest(http.MethodGet, "http://test.com/test/path", nil)
	req = PrepareRequestContext(req, 0)
	req.URL, _ = url.Parse("/path")
	req = PrepareRequestContext(req, 0)
	if fp := GetParam(req.Context(), PARAM_FULLPATH); fp != "/test/path" {
		t.Error(fp)
	}
	// Nested with non-rctx context in between
	req = httptest.NewRequest(http.MethodGet, "http://test.com/test/path", nil)
	req = PrepareRequestContext(req, 0)
	req = req.WithContext(context.WithValue(req.Context(), "someKey", "someValue"))
	if fp := GetParam(req.Context(), PARAM_FULLPATH); fp != "/test/path" {
		t.Error(fp)
	}
	req.URL, _ = url.Parse("/path")
	req = PrepareRequestContext(req, 0)
	if fp := GetParam(req.Context(), PARAM_FULLPATH); fp != "/test/path" {
		t.Error(fp)
	}
}
