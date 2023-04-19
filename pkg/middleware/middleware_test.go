package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExpectQueryParam(t *testing.T) {
	t.Run("foo=bar", func(t *testing.T) {
		m := ExpectQueryParam("foo")
		r := httptest.NewRequest("GET", "http://example.com?foo=bar", nil)
		w := httptest.NewRecorder()
		if m(w, r) != r {
			t.Error("ExpectQueryParam did not recognize foo was provided")
		}
	})

	t.Run("foo=", func(t *testing.T) {
		m := ExpectQueryParam("foo")
		r := httptest.NewRequest("GET", "http://example.com?foo=", nil)
		w := httptest.NewRecorder()
		if m(w, r) != r {
			t.Error("ExpectQueryParam did not recognize foo was provided")
		}
	})

	t.Run("foo without equals sign", func(t *testing.T) {
		m := ExpectQueryParam("foo")
		r := httptest.NewRequest("GET", "http://example.com?foo", nil)
		w := httptest.NewRecorder()
		if m(w, r) != r {
			t.Error("ExpectQueryParam should not have recognized foo was provided")
		}
	})

	t.Run("foo absent", func(t *testing.T) {
		m := ExpectQueryParam("foo")
		r := httptest.NewRequest("GET", "http://example.com?bar=foo", nil)
		w := httptest.NewRecorder()
		if m(w, r) != nil {
			t.Error("ExpectQueryParam should not have recognized foo was provided")
		}
	})
}

func TestLogRequests(t *testing.T) {
	t.Run("log request with no origin", func(t *testing.T) {
		var builder strings.Builder
		mw := LogRequests(&builder)
		req := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		req = mw(nil, req)
		if req == nil {
			t.Fatal("request was nil, should be unchanged")
		}

		log, err := ParseLog(builder.String())
		if err != nil {
			t.Fatalf("error parsing log entry: %v", err)
		}
		if log.Origin != "" {
			t.Fatalf(
				"log includes origin where there should be none, got %s",
				log.Origin,
			)
		}
		if log.Method != http.MethodGet {
			t.Fatalf(
				"log should be for a GET request, got method %s",
				log.Method,
			)
		}
		logURL := log.URL.String()
		reqURL := req.URL.String()
		if logURL != reqURL {
			t.Fatalf(
				"log should have same URL as request; expected %s but got %s",
				reqURL,
				logURL,
			)
		}
	})

	t.Run("log request with an origin", func(t *testing.T) {
		var builder strings.Builder
		mw := LogRequests(&builder)
		req := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		req.Header.Set("Origin", "origin.com")
		req = mw(nil, req)
		if req == nil {
			t.Fatal("request was nil, should be unchanged")
		}

		log, err := ParseLog(builder.String())
		if err != nil {
			t.Fatalf("error parsing log entry: %v", err)
		}
		if log.Origin != "origin.com" {
			t.Fatalf(
				"log should have an origin of origin.com, got %s",
				log.Origin,
			)
		}
		if log.Method != http.MethodGet {
			t.Fatalf(
				"log should be for a GET request, got method %s",
				log.Method,
			)
		}
		logURL := log.URL.String()
		reqURL := req.URL.String()
		if logURL != reqURL {
			t.Fatalf(
				"log should have same URL as request; expected %s but got %s",
				reqURL,
				logURL,
			)
		}
	})
}