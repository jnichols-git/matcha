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
