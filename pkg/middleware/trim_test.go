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
