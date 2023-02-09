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

func newRMC() *routeMatchContext {
	return &routeMatchContext{
		vals: map[param]string{},
	}
}

func (rmc *routeMatchContext) Clone() *routeMatchContext {
	return &routeMatchContext{}
}

func (rmc *routeMatchContext) ResetOnto(ctx context.Context) {
	rmc.parent = ctx
	for k := range rmc.vals {
		rmc.vals[k] = ""
	}
}

func (rmc *routeMatchContext) Allocate(key string) {
	rmc.vals[param(key)] = ""
}

func (rmc *routeMatchContext) SetParamIfAllocated(key param, value string) error {
	if _, ok := rmc.vals[key]; ok {
		rmc.vals[key] = value
		return nil
	}
	return errors.New(fmt.Sprintf("key %s not pre-allocated", key))
}

// =====CONTEXT IMPLEMENTATION=====

// routeMatchContext does not currently support deadlines.
func (rmc *routeMatchContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

// routeMatchContext does not currently support doneness signals.
func (rmc *routeMatchContext) Done() <-chan struct{} {
	return nil
}

func (rmc *routeMatchContext) Err() error {
	return rmc.err
}
func (rmc *routeMatchContext) Value(key any) any {
	if pkey, ok := key.(param); ok {
		return rmc.vals[pkey]
	} else {
		return rmc.parent.Value(key)
	}
}

func setParam(rmc *routeMatchContext, key, val string) {
	rmc.SetParamIfAllocated(param(key), val)
}

func GetParam(c context.Context, key string) string {
	if rmc, ok := c.(*routeMatchContext); ok {
		val, _ := rmc.vals[param(key)]
		return val
	}
	return ""
}
