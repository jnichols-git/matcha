package middleware

import (
	"net/http/httptest"
	"testing"
)

func TestExpectQueryParam(t *testing.T) {
	t.Run("foo=bar", func(t *testing.T) {
		m := ExpectQueryParam("foo", "bar")
		r := httptest.NewRequest("GET", "http://example.com?foo=bar", nil)
		w := httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != r {
			t.Error("ExpectQueryParam did not recognize foo was provided")
		}
		r = httptest.NewRequest("GET", "http://example.com?foo=baz", nil)
		w = httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != nil {
			t.Error("ExpectQueryParam should not have recognized foo was provided")
		}
	})

	t.Run("foo=bar,baz", func(t *testing.T) {
		m := ExpectQueryParam("foo", "{bar|baz}", "{[}")
		r := httptest.NewRequest("GET", "http://example.com?foo=bar", nil)
		w := httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != r {
			t.Error("ExpectQueryParam did not recognize foo was provided")
		}
		r = httptest.NewRequest("GET", "http://example.com?foo=baz", nil)
		w = httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != r {
			t.Error("ExpectQueryParam did not recognize foo was provided")
		}
		r = httptest.NewRequest("GET", "http://example.com?foo=bop", nil)
		w = httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != nil {
			t.Error("ExpectQueryParam should not have recognized foo was provided")
		}
	})

	t.Run("foo=", func(t *testing.T) {
		m := ExpectQueryParam("foo")
		r := httptest.NewRequest("GET", "http://example.com?foo=", nil)
		w := httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != r {
			t.Error("ExpectQueryParam did not recognize foo was provided")
		}
	})

	t.Run("foo without equals sign", func(t *testing.T) {
		m := ExpectQueryParam("foo")
		r := httptest.NewRequest("GET", "http://example.com?foo", nil)
		w := httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != r {
			t.Error("ExpectQueryParam did not recognize foo was provided")
		}
	})

	t.Run("foo absent", func(t *testing.T) {
		m := ExpectQueryParam("foo")
		r := httptest.NewRequest("GET", "http://example.com?bar=foo", nil)
		w := httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != nil {
			t.Error("ExpectQueryParam should not have recognized foo was provided")
		}
	})
}

func TestExpectHeader(t *testing.T) {
	t.Run("foo: bar", func(t *testing.T) {
		m := ExpectHeader("foo", "bar")
		r := httptest.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("foo", "bar")
		w := httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != r {
			t.Error("ExpectHeader did not recognize foo was provided")
		}
		r = httptest.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("foo", "baz")
		w = httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != nil {
			t.Error("ExpectHeader should not have recognized foo was provided")
		}
	})

	t.Run("foo: bar,baz", func(t *testing.T) {
		m := ExpectHeader("foo", "{bar|baz}", "{[}")
		r := httptest.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("foo", "bar")
		w := httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != r {
			t.Error("ExpectHeader did not recognize foo was provided")
		}
		r = httptest.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("foo", "baz")
		w = httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != r {
			t.Error("ExpectHeader did not recognize foo was provided")
		}
		r = httptest.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("foo", "bop")
		w = httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != nil {
			t.Error("ExpectHeader should not have recognized foo was provided")
		}
	})

	t.Run("foo", func(t *testing.T) {
		m := ExpectHeader("foo")
		r := httptest.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("foo", "bar")
		w := httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != r {
			t.Error("ExpectHeader should have recognized foo was provided")
		}
		r = httptest.NewRequest("GET", "http://example.com", nil)
		r.Header.Set("foo", "")
		w = httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != nil {
			t.Error("ExpectHeader should not have recognized foo was provided")
		}
	})

	t.Run("foo absent", func(t *testing.T) {
		m := ExpectHeader("foo")
		r := httptest.NewRequest("GET", "http://example.com", nil)
		w := httptest.NewRecorder()
		if ExecuteMiddleware([]Middleware{m}, w, r) != nil {
			t.Error("ExpectHeader should not have recognized foo was provided")
		}
	})
}
