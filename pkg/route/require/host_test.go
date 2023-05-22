package require

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireHosts(t *testing.T) {
	rq := Hosts("localhost", "{.+}.cloudretic.com")
	// Positive cases
	req := httptest.NewRequest(http.MethodGet, "http://localhost:3000", nil)
	if !rq(req) {
		t.Error("expected match")
	}
	for i := 3001; i <= 4000; i++ {
		url := fmt.Sprintf("http://localhost:%d", i)
		req = httptest.NewRequest(http.MethodGet, url, nil)
		if !rq(req) {
			t.Error("expected match", url)
		}
	}
	req = httptest.NewRequest(http.MethodGet, "http://localhost:4500", nil)
	if !rq(req) {
		t.Error("expected match")
	}
	req = httptest.NewRequest(http.MethodGet, "https://www.cloudretic.com:443", nil)
	if !rq(req) {
		t.Error("expected match")
	}
	req = httptest.NewRequest(http.MethodGet, "https://api.cloudretic.com:443", nil)
	if !rq(req) {
		t.Error("expected match")
	}
	req = httptest.NewRequest(http.MethodGet, "https://api.cloudretic.com", nil)
	if !rq(req) {
		t.Error("expected match")
	}
	req = httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	if !rq(req) {
		t.Error("expected match")
	}
	req = httptest.NewRequest(http.MethodGet, "http://api.cloudretic.com", nil)
	if !rq(req) {
		t.Error("expected match")
	}
}

func TestRequireHostPorts(t *testing.T) {
	rq := HostPorts("localhost:3000", "localhost:3001-4000,4500", "https://{.+}.cloudretic.com")
	// Positive cases
	req := httptest.NewRequest(http.MethodGet, "http://localhost:3000", nil)
	if !rq(req) {
		t.Error("expected match")
	}
	for i := 3001; i <= 4000; i++ {
		url := fmt.Sprintf("http://localhost:%d", i)
		req = httptest.NewRequest(http.MethodGet, url, nil)
		if !rq(req) {
			t.Error("expected match", url)
		}
	}
	req = httptest.NewRequest(http.MethodGet, "http://localhost:4500", nil)
	if !rq(req) {
		t.Error("expected match")
	}
	req = httptest.NewRequest(http.MethodGet, "https://www.cloudretic.com:443", nil)
	if !rq(req) {
		t.Error("expected match")
	}
	req = httptest.NewRequest(http.MethodGet, "https://api.cloudretic.com:443", nil)
	if !rq(req) {
		t.Error("expected match")
	}
	req = httptest.NewRequest(http.MethodGet, "https://api.cloudretic.com", nil)
	if !rq(req) {
		t.Error("expected match")
	}
	// Negative cases
	req = httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	if rq(req) {
		t.Error("expected no match")
	}
	req = httptest.NewRequest(http.MethodGet, "http://api.cloudretic.com", nil)
	if rq(req) {
		t.Error("expected no match")
	}
}
