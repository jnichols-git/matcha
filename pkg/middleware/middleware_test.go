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
	t.Run("log requests", func(t *testing.T) {
		var builder strings.Builder
		mw := LogRequests(&builder)
		r := mw(nil, httptest.NewRequest(http.MethodGet, "https://example.com/", nil))
		if r == nil {
			t.Fatal("request was nil, should be defined")
		}

		result := builder.String()
		if result != "GET https://example.com/\n" {
			t.Fatalf(
				"incorrect log for request %s %s:\n  got %v",
				r.Method,
				r.URL.String(),
				result,
			)
		}
	})

	t.Run("log requests if", func(t *testing.T) {
		var builder strings.Builder
		methodIsGet := func(r *http.Request) bool {
			return r.Method == http.MethodGet
		}
		mw := LogRequestsIf(methodIsGet, &builder)

		r := mw(nil, httptest.NewRequest(http.MethodGet, "https://example.com/", nil))
		if r == nil {
			t.Fatal("request was nil, should be defined")
		}

		result := builder.String()
		if result != "GET https://example.com/\n" {
			t.Fatalf(
				"incorrect log for request %s %s:\n  got %v",
				r.Method,
				r.URL.String(),
				result,
			)
		}

		builder.Reset()
		r = mw(nil, httptest.NewRequest(http.MethodPost, "https://example.com/", nil))
		if r == nil {
			t.Fatal("request was nil, should be defined")
		}

		result = builder.String()
		if result != "" {
			t.Fatalf(
				"incorrect log for request %s %s:\n  got %v",
				r.Method,
				r.URL.String(),
				result,
			)
		}
	})
}