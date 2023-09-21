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
	"context"
	"net/http"
	"sync"
	"sync/atomic"
)

type idt bool

var ids atomic.Int64
var idp = &sync.Pool{
	New: func() any {
		id := ids.Add(1)
		return id
	},
}
var handling = make(map[int64]*http.Request)
var hlock = &sync.RWMutex{}

// Handler allows integration of traditional handler-chain-ware.
func Handler(create func(next http.Handler) http.Handler) Middleware {
	next := func(w http.ResponseWriter, req *http.Request) {
		id := req.Context().Value(idt(true)).(int64)
		hlock.Lock()
		handling[id] = req
		hlock.Unlock()
	}
	h := create(http.HandlerFunc(next))
	return func(w http.ResponseWriter, req *http.Request) *http.Request {
		id := idp.Get().(int64)
		ctx := context.WithValue(req.Context(), idt(true), id)
		req = req.WithContext(ctx)
		h.ServeHTTP(w, req)
		hlock.RLock()
		out := handling[id]
		hlock.RUnlock()
		idp.Put(id)
		return out
	}
}
