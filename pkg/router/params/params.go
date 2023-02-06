// package params handles the use of route parameters.
// route.Parts should use params to communicate values such as wildcard or regex parameters.
package params

import (
	"context"
	"net/http"
)

type param string

func Set(req *http.Request, k, v string) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), param(k), v))
}

func Get(req *http.Request, k string) (string, bool) {
	if v := req.Context().Value(param(k)); v == nil {
		return "", false
	} else {
		s, ok := v.(string)
		return s, ok
	}
}
