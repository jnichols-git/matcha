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
