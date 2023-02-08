// package params handles the use of route parameters.
// route.Parts should use params to communicate values such as wildcard or regex parameters.
package params

import (
	"context"
	"net/http"
)

type param string

func Set(req *http.Request, kvs map[string]string) *http.Request {
	ctx := req.Context()
	for k, v := range kvs {
		ctx = context.WithValue(ctx, param(k), v)
	}
	return req.WithContext(ctx)
}

func Get(req *http.Request, k string) (string, bool) {
	if v := req.Context().Value(param(k)); v == nil {
		return "", false
	} else {
		s, ok := v.(string)
		return s, ok
	}
}
