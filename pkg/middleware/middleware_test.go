package middleware

import (
	"errors"
	"fmt"
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
		if ExecuteMiddleware([]Middleware{m}, w, r) != r {
			t.Error("ExpectQueryParam did not recognize foo was provided")
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
			t.Error("ExpectQueryParam should not have recognized foo was provided")
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

func TestLogRequests(t *testing.T) {
	t.Run("log request with no origin", func(t *testing.T) {
		var builder strings.Builder
		mw := LogRequests(&builder)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		req = ExecuteMiddleware([]Middleware{mw}, w, req)
		if req == nil {
			t.Fatal("request was nil, should be unchanged")
		}

		log, err := ParseLog(builder.String())
		if err != nil {
			t.Errorf("error parsing log entry: %v", err)
		}
		if log.Origin != "" {
			t.Errorf(
				"log includes origin where there should be none, got %s",
				log.Origin,
			)
		}
		if log.Method != http.MethodGet {
			t.Errorf(
				"log should be for a GET request, got method %s",
				log.Method,
			)
		}
		logURL := log.URL.String()
		reqURL := req.URL.String()
		if logURL != reqURL {
			t.Errorf(
				"log should have same URL as request; expected %s but got %s",
				reqURL,
				logURL,
			)
		}
	})

	t.Run("log request with an origin", func(t *testing.T) {
		var builder strings.Builder
		mw := LogRequests(&builder)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		req.Header.Set("Origin", "origin.com")
		req = ExecuteMiddleware([]Middleware{mw}, w, req)
		if req == nil {
			t.Fatal("request was nil, should be unchanged")
		}

		log, err := ParseLog(builder.String())
		if err != nil {
			t.Errorf("error parsing log entry: %v", err)
		}
		if log.Origin != "origin.com" {
			t.Errorf(
				"log should have an origin of origin.com, got %s",
				log.Origin,
			)
		}
		if log.Method != http.MethodGet {
			t.Errorf(
				"log should be for a GET request, got method %s",
				log.Method,
			)
		}
		logURL := log.URL.String()
		reqURL := req.URL.String()
		if logURL != reqURL {
			t.Errorf(
				"log should have same URL as request; expected %s but got %s",
				reqURL,
				logURL,
			)
		}
	})
}

func TestLogRequestsIf(t *testing.T) {
	t.Run("log request with origin cloudretic.com", func(t *testing.T) {
		var builder strings.Builder
		mw := LogRequestsIf(func(r *http.Request) bool {
			return r.Header.Get("Origin") == "cloudretic.com"
		}, &builder)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		req.Header.Set("Origin", "cloudretic.com")
		req = ExecuteMiddleware([]Middleware{mw}, w, req)
		if req == nil {
			t.Fatal("request was nil, should be unchanged")
		}
		log, err := ParseLog(builder.String())
		if err != nil {
			t.Errorf("error parsing log entry: %v", err)
		}
		if log.Origin != "cloudretic.com" {
			t.Errorf(
				"log should have an origin of cloudretic.com, got %s",
				log.Origin,
			)
		}
		if log.Method != http.MethodGet {
			t.Errorf(
				"log should be for a GET request, got method %s",
				log.Method,
			)
		}
		logURL := log.URL.String()
		reqURL := req.URL.String()
		if logURL != reqURL {
			t.Errorf(
				"log should have same URL as request; expected %s but got %s",
				reqURL,
				logURL,
			)
		}

		builder.Reset()
		w = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		req.Header.Set("Origin", "other.com")
		req = ExecuteMiddleware([]Middleware{mw}, w, req)
		if req == nil {
			t.Fatal("request was nil, should be unchanged")
		}
		if builder.String() != "" {
			t.Errorf("No log should occur on this origin, got %s", builder.String())
		}
	})
}

type badWriter struct{}

func (bw *badWriter) Write([]byte) (int, error) {
	return 0, errors.New("no")
}

func TestFailure(t *testing.T) {
	mw1 := LogRequests(&badWriter{})
	req := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
	req = mw1(nil, req)
	if req == nil {
		t.Errorf("request should not be rejected on failure")
	}
	badlog := "incorrectly formatted log"
	l, err := ParseLog(badlog)
	if l != nil || err == nil {
		t.Errorf("bad log should return nil, err, got %v, %s", l, err)
	}
	badurl := "google.com/" + string([]byte{0x7f})
	badlog = "0 - GET " + badurl
	l, err = ParseLog(badlog)
	if l != nil || err == nil {
		t.Errorf("bad log should return nil, err, got %v, %s", l, err)
	}
}

func TestHandlerAsMiddleware(t *testing.T) {
	expectedResponse := "Hello, world!"
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, expectedResponse)
	}
	mw := Handler(http.HandlerFunc(handler))
	req := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
	rec := httptest.NewRecorder()
	req = mw(rec, req)
	if req == nil {
		t.Errorf("request should be unchanged.")
	}
	if res := rec.Body.String(); res != "Hello, world!" {
		t.Errorf(
			"response should be \"%s\", got \"%s\"",
			expectedResponse,
			res,
		)
	}
}