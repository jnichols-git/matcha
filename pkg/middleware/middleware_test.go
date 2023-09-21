// Copyright 2023 Decent Platforms
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
	t.Run("log request with origin decentplatforms.com", func(t *testing.T) {
		var builder strings.Builder
		mw := LogRequestsIf(func(r *http.Request) bool {
			return r.Header.Get("Origin") == "decentplatforms.com"
		}, &builder)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "https://example.com/", nil)
		req.Header.Set("Origin", "decentplatforms.com")
		req = ExecuteMiddleware([]Middleware{mw}, w, req)
		if req == nil {
			t.Fatal("request was nil, should be unchanged")
		}
		log, err := ParseLog(builder.String())
		if err != nil {
			t.Errorf("error parsing log entry: %v", err)
		}
		if log.Origin != "decentplatforms.com" {
			t.Errorf(
				"log should have an origin of decentplatforms.com, got %s",
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
