package route

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type param string

// routeMatchContext is a specialization of context for route parameters.
type routeMatchContext struct {
	parent context.Context
	// track context values
	vals map[param]string
	// current context error
	err error
}

// =====RMC IMPLEMENTATION=====

// Create a new routeMatchContext.
func newRMC() *routeMatchContext {
	return &routeMatchContext{
		vals: map[param]string{},
	}
}

// Reset a routeMatchContext onto a regular context.Context.
// This does two things:
//   - Set the parent context of rmc to ctx, allowing middleware modifications to ctx to continue to apply
//   - Reset all preallocated parameter values
func (rmc *routeMatchContext) ResetOnto(ctx context.Context) {
	rmc.parent = ctx
	for k := range rmc.vals {
		rmc.vals[k] = ""
	}
}

// Allocate space for a route parameter.
// Allocation is required to set route parameters later.
func (rmc *routeMatchContext) Allocate(key string) {
	rmc.vals[param(key)] = ""
}

// Set a route parameter if it was preallocated in the context.
func (rmc *routeMatchContext) SetParamIfAllocated(key param, value string) error {
	if _, ok := rmc.vals[key]; ok {
		rmc.vals[key] = value
		return nil
	}
	return errors.New(fmt.Sprintf("key %s not pre-allocated", key))
}

// =====CONTEXT IMPLEMENTATION=====

// routeMatchContext does not natively support deadlines.
// If the parent context has a deadline, that will be returned.
//
// See interface context.Context.
func (rmc *routeMatchContext) Deadline() (time.Time, bool) {
	if rmc.parent != nil {
		return rmc.parent.Deadline()
	}
	return time.Time{}, false
}

// routeMatchContext does not natively support doneness signals.
// If the parent context has a doneness signal, that will be returned.
//
// See interface context.Context.
func (rmc *routeMatchContext) Done() <-chan struct{} {
	if rmc.parent != nil {
		return rmc.parent.Done()
	}
	return nil
}

// Return the current error on the context, or the error on the parent if applicable.
//
// See interface context.Context.
func (rmc *routeMatchContext) Err() error {
	if rmc.err != nil {
		return rmc.err
	} else if rmc.parent != nil {
		return rmc.parent.Err()
	}
	return nil
}

// Get a value.
//
// See interface context.Context.
func (rmc *routeMatchContext) Value(key any) any {
	if pkey, ok := key.(param); ok {
		return rmc.vals[pkey]
	} else {
		return rmc.parent.Value(key)
	}
}

// =====ROUTE PARAMS=====

// Set a route parameter.
// This is only permitted in package route; if you're defining values elsewhere,
// use context.WithValue.
func setParam(rmc *routeMatchContext, key, val string) {
	rmc.SetParamIfAllocated(param(key), val)
}

// Get a route parameter.
// Returns "" if the parameter does not exist.
func GetParam(c context.Context, key string) string {
	if rmc, ok := c.(*routeMatchContext); ok {
		val, _ := rmc.vals[param(key)]
		return val
	}
	return ""
}
