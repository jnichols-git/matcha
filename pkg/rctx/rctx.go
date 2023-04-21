package rctx

import (
	"context"
	"errors"
	"net/http"
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

// PrepareRequestContext prepares the context of a request for matching.
func PrepareRequestContext(req *http.Request, maxParams int) *http.Request {
	rctx := &Context{
		parent: req.Context(),
		params: newParams(maxParams),
		err:    nil,
	}
	return req.WithContext(rctx)
}

func Head(req *http.Request) (int, error) {
	ctx := req.Context()
	rctx, correctType := ctx.(*Context)
	if !correctType {
		return -1, errors.New("request must have *rctx.Context")
	}
	return rctx.params.head, nil
}

// ResetRequestContext resets the param head to the provided value.
// This effectively deletes any params assigned after that head, allowing for partial resets.
func ResetRequestContextHead(req *http.Request, to int) error {
	ctx := req.Context()
	rctx, correctType := ctx.(*Context)
	if !correctType {
		return errors.New("request must have *rctx.Context when resetting")
	}
	if rctx.params.head < to {
		return errors.New("cannot move the context head forward")
	}
	rctx.params.head = to
	return nil
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

// PARAMETER IMPLEMENTATION

// GetParam gets a parameter by its key string.
// This automatically converts key to its underlying context key type.
// GetParam is unique in that it's the only native function that doesn't fail on non-*rctx.Context types;
// it's possible that the context type changes between the route match and handling, so for other contexts
// the key is passed to ctx.Value(paramKey(key)).
func GetParam(ctx context.Context, key string) string {
	if rctx, ok := ctx.(*Context); ok {
		return rctx.params.get(paramKey(key))
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
		return rctx.params.set(paramKey(key), value)
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
		return ctx.params.get(pkey)
	} else if ctx.parent != nil {
		return ctx.parent.Value(key)
	} else {
		return nil
	}
}
