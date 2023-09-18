package rctx

import (
	"context"
	"testing"
)

func TestReserved(t *testing.T) {
	// Full path should set exactly once each reset
	rps := reservedParams{}
	ctx := context.Background()
	rps.set(ctx, key_reserved_fullpath, "/test/path")
	if fp, reserved := rps.get(key_reserved_fullpath); !reserved || fp != "/test/path" {
		t.Error(fp, reserved)
	}
	rps.set(ctx, key_reserved_fullpath, "/other/test/path")
	if fp, reserved := rps.get(key_reserved_fullpath); !reserved || fp != "/test/path" {
		t.Error(fp, reserved)
	}
	rps.reset()
	rps.set(nil, key_reserved_fullpath, "/other/test/path")
	if fp, reserved := rps.get(key_reserved_fullpath); !reserved || fp != "/other/test/path" {
		t.Error(fp, reserved)
	}
	// A parent context with this value should be accessed instead of the child context.
	ctx = context.WithValue(ctx, key_reserved_fullpath, "/parent/path")
	rps.reset()
	rps.set(ctx, key_reserved_fullpath, "/test/path")
	if fp, reserved := rps.get(key_reserved_fullpath); !reserved || fp != "/parent/path" {
		t.Error(fp, reserved)
	}
	// Mount proxy-to path should be freely settable
	rps.set(ctx, key_reserved_mountProxyTo, "/test/path")
	if fp, reserved := rps.get(key_reserved_mountProxyTo); !reserved || fp != "/test/path" {
		t.Error(fp, reserved)
	}
	rps.set(ctx, key_reserved_mountProxyTo, "/other/test/path")
	if fp, reserved := rps.get(key_reserved_mountProxyTo); !reserved || fp != "/other/test/path" {
		t.Error(fp, reserved)
	}
	// Other keys should be ignored
	if reserved, err := rps.set(ctx, paramKey("other-key"), "anyvalue"); reserved || err != nil {
		t.Error(reserved, err)
	}
}
