package middleware

import (
	"net/http"
	"strings"
)

func TrimPrefix(prefix string) Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		path := r.URL.Path
		if strings.Index(path, prefix) == 0 {
			r.URL.Path = path[len(prefix):]
		}
		return r
	}
}

func TrimPrefixStrict(prefix string, errMsg string) Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		path := r.URL.Path
		if strings.Index(path, prefix) == 0 {
			r.URL.Path = path[len(prefix):]
		} else {
			w.WriteHeader(http.StatusBadRequest)
			if errMsg == "" {
				errMsg = "expected path prefix " + prefix
			}
			w.Write([]byte(errMsg))
			return nil
		}
		return r
	}
}
