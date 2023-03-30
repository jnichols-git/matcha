package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExpectQueryParam(t *testing.T) {
	t.Run("foo=bar", func(t *testing.T) {
		m := ExpectQueryParam("foo")
		r, _ := http.NewRequest("GET", "http://example.com?foo=bar", nil)
		w := httptest.NewRecorder()
		if m(w, r) != r {
			t.Error("ExpectQueryParam did not recognize foo was provided")
		}
	})

	t.Run("foo=", func(t *testing.T) {
		m := ExpectQueryParam("foo")
		r, _ := http.NewRequest("GET", "http://example.com?foo=", nil)
		w := httptest.NewRecorder()
		if m(w, r) != r {
			t.Error("ExpectQueryParam did not recognize foo was provided")
		}
	})

	t.Run("foo without equals sign", func(t *testing.T) {
		m := ExpectQueryParam("foo")
		r, _ := http.NewRequest("GET", "http://example.com?foo", nil)
		w := httptest.NewRecorder()
		if m(w, r) != r {
			t.Error("ExpectQueryParam should not have recognized foo was provided")
		}
	})

	t.Run("foo absent", func(t *testing.T) {
		m := ExpectQueryParam("foo")
		r, _ := http.NewRequest("GET", "http://example.com?bar=foo", nil)
		w := httptest.NewRecorder()
		if m(w, r) != nil {
			t.Error("ExpectQueryParam should not have recognized foo was provided")
		}
	})
}
