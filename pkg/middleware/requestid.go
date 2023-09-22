package middleware

import (
	"context"
	"net/http"
)

type RequestIDGenerator func() string

type requestIdKey string

func RequestID(generator RequestIDGenerator, key string, isContext bool) Middleware {
	return func(_ http.ResponseWriter, r *http.Request) *http.Request {
		if isContext {
			ctx := context.WithValue(r.Context(), requestIdKey(key), generator())
			return r.WithContext(ctx)
		}

		r.Header.Add(key, generator())

		return r
	}
}
