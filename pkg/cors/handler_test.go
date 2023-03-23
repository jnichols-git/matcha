package cors

import (
	"net/http"
	"net/url"
	"testing"
)

func TestHandler(t *testing.T) {
	test_expr := "/invalid/rou!e"
	_, _, err := PreflightHandler(test_expr, aco2)
	if err == nil {
		t.Fatal("expected failure to build invalid route")
	}
	test_expr = "/static/route"
	r, h, err := PreflightHandler(test_expr, aco2)
	if err != nil {
		t.Fatal(err)
	}
	adp := &testAdapter{}
	w, req, res, _ := adp.Adapt(preflight_request)
	req.URL = &url.URL{Path: "/static/route"}
	if r.MatchAndUpdateContext(req) == nil {
		t.Error("path should match")
	}
	h.ServeHTTP(w, req)
	if res.StatusCode != 204 {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}
	ao := res.Headers[AllowOrigin]
	if len(ao) != 1 || ao[0] != req.Header.Get("Origin") {
		t.Errorf("expected allow-origin '%s', got %v", req.Header.Get("Origin"), ao)
	}
	am := res.Headers[AllowMethods]
	if len(am) != 1 || am[0] != http.MethodPost {
		t.Errorf("expected allow-method '%s', got %v", req.Method, am)
	}
	ah := res.Headers[AllowHeaders]
	if len(ah) != 2 || ah[0] != "X-Header-1" || ah[1] != "X-Header-2" {
		t.Errorf("expected allowed headers 'X-Header-1,X-Header-2', got %v", ah)
	}
	eh := res.Headers[ExposeHeaders]
	if len(eh) != 1 || eh[0] != "X-Header-Out" {
		t.Errorf("expected exposed header X-Header-Out, got %v", eh)
	}
}
