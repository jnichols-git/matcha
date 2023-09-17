package rctx

import (
	"context"
)

const FULLPATH = string("matcha_fullpath")
const key_reserved_fullpath = paramKey(FULLPATH)
const MOUNTPROXYTO = string("matcha_mountProxyTo")
const key_reserved_mountProxyTo = paramKey(MOUNTPROXYTO)

type reservedParams struct {
	fullpath     string
	mountProxyTo string
}

// get gets a reserved rctx param.
// Returns the value, and true if key is reserved.
func (rps *reservedParams) get(key paramKey) (string, bool) {
	switch key {
	case key_reserved_fullpath:
		return rps.fullpath, true
	case key_reserved_mountProxyTo:
		return rps.mountProxyTo, true
	default:
		return "", false
	}
}

// set sets a reserved rctx param.
// Returns true if the key is reserved.
func (rps *reservedParams) set(parent context.Context, key paramKey, value string) (bool, error) {
	switch key {
	case key_reserved_fullpath:
		if rps.fullpath == "" {
			if parent == nil {
				rps.fullpath = value
			} else if orig := parent.Value(key_reserved_fullpath); orig != nil && orig.(string) != "" {
				rps.fullpath = orig.(string)
			} else {
				rps.fullpath = value
			}
		}
		return true, nil
	case key_reserved_mountProxyTo:
		rps.mountProxyTo = value
		return true, nil
	default:
		return false, nil
	}
}

// reset resets the reserved params.
// Some reserved params have special behavior when set multiple times; this sets back to
// default values so that behavior can be replicated on pooled context.
func (rps *reservedParams) reset() {
	rps.fullpath = ""
	rps.mountProxyTo = ""
}
