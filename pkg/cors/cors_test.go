package cors

import (
	"net/http"
	"testing"
)

var aco1 = &AccessControlOptions{
	AllowOrigin:      []string{"*"},
	AllowMethods:     []string{"*"},
	AllowHeaders:     []string{"*"},
	ExposeHeaders:    []string{"*"},
	MaxAge:           1000,
	AllowCredentials: false,
}
var aco2 = &AccessControlOptions{
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

func TestGetCORSRequestHeaders(t *testing.T) {
	crh := GetCORSRequestHeaders(simple_request)
	if crh.Origin != "origin.com" {
		t.Errorf("expected origin 'origin.com', got '%s'", crh.Origin)
	}
	crh = GetCORSRequestHeaders(preflight_request)
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

func BenchmarkGetCORSRequestHeaders(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetCORSRequestHeaders(preflight_request)
	}
}

func TestReflectCORSRequestHeaders(t *testing.T) {
	crh := GetCORSRequestHeaders(simple_request)
	resp_headers := ReflectCORSRequestHeaders(aco1, crh)
	if len(resp_headers.AllowOrigin) != 1 || resp_headers.AllowOrigin[0] != "origin.com" {
		t.Errorf("expected allow-origin 'origin.com', got %v", resp_headers.AllowOrigin)
	}
	if len(resp_headers.AllowMethods) != 1 || resp_headers.AllowMethods[0] != http.MethodGet {
		t.Errorf("expected allow-methods 'GET', got %v", resp_headers.AllowMethods)
	}
	if len(resp_headers.AllowHeaders) != 0 {
		t.Errorf("expected no allowed additional headers, got %v", resp_headers.AllowHeaders)
	}
	if resp_headers.AllowCredentials {
		t.Errorf("expected no allowed credentials")
	}

	crh = GetCORSRequestHeaders(preflight_request)
	resp_headers = ReflectCORSRequestHeaders(aco2, crh)
	if len(resp_headers.AllowOrigin) != 1 || resp_headers.AllowOrigin[0] != "origin.com" {
		t.Errorf("expected allow-origin 'origin.com', got %v", resp_headers.AllowOrigin)
	}
	if len(resp_headers.AllowMethods) != 1 || resp_headers.AllowMethods[0] != http.MethodPost {
		t.Errorf("expected allow-methods 'POST', got %v", resp_headers.AllowMethods)
	}
	if len(resp_headers.AllowHeaders) != 2 || resp_headers.AllowHeaders[0] != "X-Header-1" || resp_headers.AllowHeaders[1] != "x-Header-2" {
		t.Errorf("expected allow-headers X-Header-1 and X-Header-2, got %v", resp_headers.AllowHeaders)
	}
	if len(resp_headers.ExposeHeaders) != 1 || resp_headers.ExposeHeaders[0] != "x-Header-Out" {
		t.Errorf("expected expose-headers x-Header-Out, got %v", resp_headers.ExposeHeaders)
	}
	if resp_headers.AllowCredentials {
		t.Errorf("expected no allowed credentials")
	}
}

func BenchmarkReflectCORSRequestHeaders(b *testing.B) {
	for i := 0; i < b.N; i++ {
		crh := GetCORSRequestHeaders(simple_request)
		_ = ReflectCORSRequestHeaders(aco1, crh)
	}
}
