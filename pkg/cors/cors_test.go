package cors

import (
	"net/http"
	"testing"
)

var aco1 = &Options{
	AllowOrigin:      []string{"*"},
	AllowMethods:     []string{"*"},
	AllowHeaders:     []string{"*"},
	ExposeHeaders:    []string{"*"},
	MaxAge:           1000,
	AllowCredentials: false,
}
var aco2 = &Options{
	AllowOrigin:      []string{"origin.com"},
	AllowMethods:     []string{http.MethodGet, http.MethodPost},
	AllowHeaders:     []string{"x-Header-1", "X-Header-2", "x-Header-3"},
	ExposeHeaders:    []string{"x-Header-Out"},
	MaxAge:           1000,
	AllowCredentials: false,
}

var simple_request = &http.Request{
	Method: http.MethodGet,
	Header: http.Header{
		"Origin": {"origin.com"},
	},
}

var preflight_request = &http.Request{
	Method: http.MethodOptions,
	Header: http.Header{
		"Origin":                         {"origin.com"},
		"Access-Control-Request-Method":  {http.MethodPost},
		"Access-Control-Request-Headers": {"X-Header-1", "x-Header-2"},
	},
}

func TestGetRequest(t *testing.T) {
	crh := GetRequest(simple_request)
	if crh.Origin != "origin.com" {
		t.Errorf("expected origin 'origin.com', got '%s'", crh.Origin)
	}
	crh = GetRequest(preflight_request)
	if crh.Origin != "origin.com" {
		t.Errorf("expected origin 'origin.com', got '%s'", crh.Origin)
	}
	if crh.RequestMethod != http.MethodPost {
		t.Errorf("expected method '%s', got '%s'", http.MethodPost, crh.RequestMethod)
	}
	if len(crh.RequestHeaders) != 2 || crh.RequestHeaders[0] != "X-Header-1" || crh.RequestHeaders[1] != "x-Header-2" {
		t.Errorf("expected headers %v, got %v", []string{"X-Header-1", "x-Header-2"}, crh.RequestHeaders)
	}
}

func BenchmarkGetRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetRequest(preflight_request)
	}
}

func TestReflectRequest(t *testing.T) {
	crh := GetRequest(simple_request)
	resp_headers := http.Header{}
	ReflectRequest(aco1, crh, resp_headers)
	if len(resp_headers.Values(AllowOrigin)) != 1 || resp_headers.Values(AllowOrigin)[0] != "origin.com" {
		t.Errorf("expected allow-origin 'origin.com', got %v", resp_headers.Values(AllowOrigin))
	}
	if len(resp_headers.Values(AllowMethods)) != 1 || resp_headers.Values(AllowMethods)[0] != http.MethodGet {
		t.Errorf("expected allow-methods 'GET', got %v", resp_headers.Values(AllowMethods))
	}
	if len(resp_headers.Values(AllowHeaders)) != 0 {
		t.Errorf("expected no allowed additional headers, got %v", resp_headers.Values(AllowHeaders))
	}
	if resp_headers.Get(AllowCredentials) != "false" {
		t.Errorf("expected no allowed credentials")
	}

	crh = GetRequest(preflight_request)
	resp_headers = http.Header{}
	ReflectRequest(aco2, crh, resp_headers)
	if len(resp_headers.Values(AllowOrigin)) != 1 || resp_headers.Values(AllowOrigin)[0] != "origin.com" {
		t.Errorf("expected allow-origin 'origin.com', got %v", resp_headers.Values(AllowOrigin))
	}
	if len(resp_headers.Values(AllowMethods)) != 1 || resp_headers.Values(AllowMethods)[0] != http.MethodPost {
		t.Errorf("expected allow-methods 'POST', got %v", resp_headers.Values(AllowMethods))
	}
	if len(resp_headers.Values(AllowHeaders)) != 2 || resp_headers.Values(AllowHeaders)[0] != "X-Header-1" || resp_headers.Values(AllowHeaders)[1] != "x-Header-2" {
		t.Errorf("expected allow-headers X-Header-1 and X-Header-2, got %v", resp_headers.Values(AllowHeaders))
	}
	if len(resp_headers.Values(ExposeHeaders)) != 1 || resp_headers.Values(ExposeHeaders)[0] != "x-Header-Out" {
		t.Errorf("expected expose-headers x-Header-Out, got %v", resp_headers.Values(ExposeHeaders))
	}
	if resp_headers.Get(AllowCredentials) != "false" {
		t.Errorf("expected no allowed credentials")
	}
}

func BenchmarkReflectRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := http.Header{}
		crh := GetRequest(simple_request)
		ReflectRequest(aco1, crh, h)
	}
}
