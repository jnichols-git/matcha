package cors

import (
	"net/http"
	"testing"
)

func TestCORSMiddleware(t *testing.T) {
	adp := &testAdapter{}
	w, req, res, _ := adp.Adapt(simple_request)
	mid := CORSMiddleware(aco1)
	mid(w, req)
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
	mid = CORSMiddleware(aco2)
	mid(w, req)
	ao = res.Headers[AllowOrigin]
	if len(ao) != 1 || ao[0] != req.Header.Get("Origin") {
		t.Errorf("expected allow-origin '%s', got %v", req.Header.Get("Origin"), ao)
	}
	am = res.Headers[AllowMethods]
	if len(am) != 1 || am[0] != http.MethodPost {
		t.Errorf("expected allow-method '%s', got %v", req.Method, am)
	}
	ah = res.Headers[AllowHeaders]
	if len(ah) != 2 || ah[0] != "X-Header-1" || ah[1] != "X-Header-2" {
		t.Errorf("expected allowed headers 'X-Header-1,X-Header-2', got %v", ah)
	}
	eh = res.Headers[ExposeHeaders]
	if len(eh) != 1 || eh[0] != "X-Header-Out" {
		t.Errorf("expected exposed header X-Header-Out, got %v", eh)
	}
}

func BenchmarkCORSMiddleware(b *testing.B) {
	adp := &testAdapter{}
	w, req, _, _ := adp.Adapt(preflight_request)
	mid := CORSMiddleware(aco2)
	for i := 0; i < b.N; i++ {
		_ = mid(w, req)
	}
}
