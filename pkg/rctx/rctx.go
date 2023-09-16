// Package rctx defines the context structure for Matcha.
//
// See [https://github.com/decentplatforms/matcha/blob/main/docs/context.md].
package rctx

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

const (
	DefaultMaxParams = int(10)
)

type Context struct {
	parent context.Context
	params *routeParams
	err    error
}

var rctxPool = &sync.Pool{
	New: func() any {
		return &Context{}
	},
}

func new(parent context.Context, maxParams int) *Context {
	rctx := rctxPool.Get().(*Context)
	rctx.parent = parent
	if rctx.params == nil || rctx.params.cap < maxParams {
		rctx.params = newParams(maxParams)
	} else {
		rctx.params.cap = maxParams
		rctx.params.head = 0
	}
	return rctx
}

// PrepareRequestContext prepares the context of a request for matching.
func PrepareRequestContext(req *http.Request, maxParams int) *http.Request {
	rctx := new(req.Context(), maxParams)
	rctx.params.set(rctx, reserved_fullpath, req.URL.Path)
	return req.WithContext(rctx)
}

// ResetRequestContext resets any values in the context that shouldn't be maintained between attempts to match routes.
// This assumes that the request has a rctx.Context, and returns an error if it does not.
func ResetRequestContext(req *http.Request) error {
	ctx := req.Context()
	rctx, correctType := ctx.(*Context)
	if !correctType {
		return errors.New("request must have *rctx.Context when resetting")
	}
	rctx.params.head = 0
	return nil
}

func ReturnRequestContext(req *http.Request) {
	if rctx, ok := req.Context().(*Context); ok {
		rctx.parent = nil
		rctx.err = nil
		for i := range rctx.params.rps {
			rctx.params.rps[i].key = ""
			rctx.params.rps[i].value = ""
		}
		rctx.params.reserved.reset()
		rctxPool.Put(rctx)
	}
}

// PARAMETER IMPLEMENTATION

// GetParam gets a parameter by its key string.
// This automatically converts key to its underlying context key type.
// GetParam is unique in that it's the only native function that doesn't fail on non-*rctx.Context types;
// it's possible that the context type changes between the route match and handling, so for other contexts
// the key is passed to ctx.Value(paramKey(key)).
func GetParam(ctx context.Context, key string) string {
	if rctx, ok := ctx.(*Context); ok {
		v := rctx.params.get(paramKey(key))
		if v != "" {
			return v
		} else if rctx.parent != nil {
			ctx = rctx.parent
		} else {
			return ""
		}
	}
	if val, ok := ctx.Value(paramKey(key)).(string); ok {
		return val
	}
	return ""
}

// SetParam sets a parameter with a key string.
// This automatically converts key to its underlying context key type.
// Context params have a max value determined at creation, and this returns an error if the user attempts to exceed
// the maximum number of params.
func SetParam(ctx context.Context, key, value string) error {
	if rctx, ok := ctx.(*Context); ok {
		return rctx.params.set(rctx, paramKey(key), value)
	}
	return errors.New("cannot SetParam on non-rctx Context")
}

// CONTEXT IMPLEMENTATION

// rctx.Context does not natively support deadlines.
// If the parent context has a deadline, that will be returned.
//
// See interface context.Context.
func (ctx *Context) Deadline() (time.Time, bool) {
	if ctx.parent != nil {
		return ctx.parent.Deadline()
	}
	return time.Time{}, false
}

// rctx.Context does not natively support doneness signals.
// If the parent context has a doneness signal, that will be returned.
//
// See interface context.Context.
func (ctx *Context) Done() <-chan struct{} {
	if ctx.parent != nil {
		return ctx.parent.Done()
	}
	return nil
}

// Return the current error on the context, or the error on the parent if applicable.
//
// See interface context.Context.
func (ctx *Context) Err() error {
	if ctx.err != nil {
		return ctx.err
	} else if ctx.parent != nil {
		return ctx.parent.Err()
	}
	return nil
}

// Get a value.
//
// See interface context.Context.
func (ctx *Context) Value(key any) any {
	if pkey, ok := key.(paramKey); ok {
		v := ctx.params.get(pkey)
		if v != "" {
			return v
		} else if ctx.parent != nil {
			return ctx.parent.Value(key)
		} else {
			return ""
		}
	} else if ctx.parent != nil {
		return ctx.parent.Value(key)
	} else {
		return nil
	}
}
