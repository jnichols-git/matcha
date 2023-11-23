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
