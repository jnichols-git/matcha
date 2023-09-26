// Copyright 2023 Matcha Authors
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

// GetRequestIDHeader gets the request ID from a request if it was set in the header.
// Headers don't require any special behavior to access; this function mostly exists for consistency with
// GetRequestIDContext.
func GetRequestIDHeader(req *http.Request, key string) string {
	return req.Header.Get(key)
}

// GetRequestIDContext gets the request ID from a request if it was set in context.
// Since the key is typed, you *must* use this function to access the request ID. It's recommended that you
// wrap this with whatever key you use for your project as a convenience.
func GetRequestIDContext(req *http.Request, key string) string {
	if rid, ok := req.Context().Value(requestIdKey(key)).(string); ok {
		return rid
	}
	return ""
}
