package cors

import (
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
}
