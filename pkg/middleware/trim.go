package middleware

import (
	"net/http"
	"strings"
)

// TrimPrefix trims a static prefix from the path of an inbound request.
// If the prefix doesn't exist, the request is unmodified. If you want to reject requests
// without the prefix, use TrimPrefixStrict.
func TrimPrefix(prefix string) Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		path := r.URL.Path
		if strings.Index(path, prefix) == 0 {
			r.URL.Path = path[len(prefix):]
		}
		return r
	}
}

// TrimPrefixStrict trims a static prefix from the path of an inbound request.
// If the prefix doesn't exist, the request is rejected and the errMsg is sent as a response.
// An empty errMsg will generate an error message "expected path prefix [prefix]".
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
