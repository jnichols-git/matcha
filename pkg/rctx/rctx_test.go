package rctx

import (
	"fmt"
	"net/http"
	"testing"
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

func TestParams(t *testing.T) {
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
