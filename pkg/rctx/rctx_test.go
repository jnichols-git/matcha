package rctx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
		req = PrepareRequestContext(req, DefaultMaxParams)
		ctx := req.Context()
		for i := 0; i < DefaultMaxParams; i++ {
			p := fmt.Sprintf("p%d", i)
			v := fmt.Sprintf("v%d", i)
			err := SetParam(ctx, p, v)
			if err != nil {
				t.Error(err)
			}
		}
		err = SetParam(ctx, "errp", "errv")
		if err == nil {
			t.Error("should fail to set more than DefaultMaxParams")
		}
		for i := 0; i < DefaultMaxParams; i++ {
			p := fmt.Sprintf("p%d", i)
			v := fmt.Sprintf("v%d", i)
			got := GetParam(ctx, p)
			if got != v {
				t.Errorf("expected %s, got %s", v, got)
			}
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
	ctx, _ = context.WithDeadline(ctx, dl)
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
	// Test nil parent
	req = &http.Request{}
	req = PrepareRequestContext(req, DefaultMaxParams)
	if ctxdl, _ := req.Context().Deadline(); !ctxdl.IsZero() {
		t.Errorf("expected zero deadline, got %v", ctxdl)
	}
	if err := req.Context().Err(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if v := req.Context().Value("testKey"); v != nil {
		t.Errorf("expected nil, got %s", v)
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
