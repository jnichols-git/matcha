package middleware

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"net/url"
)

// Middleware runs on any incoming request. Attachment behavior is defined by the structure it's attached to (route vs. router).
//
// Returns an *http.Request; the middleware can set router params or reject a request by returning nil.
type Middleware func(http.ResponseWriter, *http.Request) *http.Request

// Returns a middleware that checks for the presence of a query parameter.
func ExpectQueryParam(name string) Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		if r.URL.Query().Has(name) {
			return r
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing query param: %s", name)
		return nil
	}
}

// The string used to indicate an absent origin in log entries.
// 
// Origins take one of the following forms, so must not take any of these forms:
//     "null"
//     "<scheme>://<hostname>"
//     "<scheme>://<hostname>:<port>"
const OriginAbsent = "-"

// Returns a middleware that logs the details of an incoming request.
func LogRequests(w io.Writer) Middleware {
	return func(_ http.ResponseWriter, r *http.Request) *http.Request {
		logRequest(w, r)
		return r
	}
}

// Returns a middleware that logs the details of an incoming request only if
// test(request) == true.
func LogRequestsIf(test func(*http.Request) bool, w io.Writer) Middleware {
	return func(_ http.ResponseWriter, r *http.Request) *http.Request {
		if test(r) {
			logRequest(w, r)
		}
		return r
	}
}

func ParseLog(s string) (*LogEntry, error) {
	var log LogEntry
	var rawTimestamp int64
	var rawOrigin string
	var rawURL string
	_, err := fmt.Sscanf(
		s,
		"%d %s %s %s",
		&rawTimestamp,
		&rawOrigin,
		&log.Method,
		&rawURL,
	)
	if err != nil {
		return nil, err
	}

	log.Timestamp = time.Unix(0, rawTimestamp)

	if rawOrigin == OriginAbsent {
		log.Origin = ""
	} else {
		log.Origin = rawOrigin
	}

	log.URL, err = url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return &log, nil
}

type LogEntry struct {
	Timestamp time.Time
	// Will be "" if no origin header was given on the request.
	Origin string
	Method string
	URL *url.URL
}

func logRequest(w io.Writer, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = OriginAbsent
	}
	fmt.Fprintf(
		w, 
		"%d %s %s %s\n",
		time.Now().UnixNano(),
		origin,
		r.Method,
		r.URL,
	)
}