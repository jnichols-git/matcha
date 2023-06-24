package middleware

import (
	"fmt"
	"net/http"

	"github.com/cloudretic/matcha/pkg/regex"
)

// ExpectQueryParam checks for the presence of a query parameter.
// Requests are rejected if the query parameter `name` is not present, or if the value
// doesn't match the provided patterns. `patts` can be left empty to permit any or no
// value assigned to `name`. Invalid patterns are silently discarded.
//
// See package Pattern for more details on pattern construction.
func ExpectQueryParam(name string, patts ...string) Middleware {
	var fs []func(v string) bool
	if len(patts) == 0 {
		fs = append(fs, func(_ string) bool {
			return true
		})
	}
	for _, patt := range patts {
		value, isPatt, err := regex.CompilePattern(patt)
		if err != nil {
			continue
		}
		var f func(v string) bool
		if isPatt {
			f = value.Match
		} else {
			f = func(v string) bool {
				return patt == v
			}
		}
		fs = append(fs, f)
	}
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		q := r.URL.Query()
		if !q.Has(name) {
			fmt.Fprintf(w, "invalid value for query param %s", name)
			return nil
		}
		v := r.URL.Query().Get(name)
		r.URL.Query()
		for _, f := range fs {
			if f(v) {
				return r
			}
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid value for query param %s", name)
		return nil
	}
}

// ExpectHeader checks for the presence of a header.
// Requests are rejected if the header `name` is not present, or if the value
// doesn't match the provided patterns, if any. `patts` can be left empty to permit
// any value assigned to `name`, but headers must have a value to be permitted.
// Invalid patterns are silently discarded.
//
// See package Pattern for more details on pattern construction.
func ExpectHeader(name string, patts ...string) Middleware {
	var fs []func(v string) bool
	if len(patts) == 0 {
		fs = append(fs, func(_ string) bool {
			return true
		})
	}
	for _, patt := range patts {
		value, isPatt, err := regex.CompilePattern(patt)
		if err != nil {
			continue
		}
		var f func(v string) bool
		if isPatt {
			f = value.Match
		} else {
			f = func(v string) bool {
				return patt == v
			}
		}
		fs = append(fs, f)
	}
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		h := r.Header
		v := h.Get(name)
		if v == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid value for header " + name))
			return nil
		}
		for _, f := range fs {
			if f(v) {
				return r
			}
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid value for header " + name))
		return nil
	}
}
