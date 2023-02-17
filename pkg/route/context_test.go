package route

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestResetOnto(t *testing.T) {
	base := context.WithValue(context.Background(), string("test"), string("OK"))
	ctx := newRMC()
	ctx.ResetOnto(base)
	if ctx.parent == nil {
		t.Fatal("ResetOnto didn't set parent")
	}
	if testVal := ctx.parent.Value(string("test")); testVal != string("OK") {
		t.Fatalf("expected parent value for 'test' to be 'OK'; got %s", testVal)
	}
}

func TestAllocate(t *testing.T) {
	ctx := newRMC()
	ctx.Allocate("test")
	if _, ok := ctx.vals[param("test")]; !ok {
		t.Fatal("allocate failed to set key")
	}
}

func TestSetParamIfAllocated(t *testing.T) {
	ctx := newRMC()
	err := ctx.SetParamIfAllocated(param("test"), "OK")
	if err == nil {
		t.Error("SetParamIfAllocated should fail for a non-allocated param")
	}
	ctx.Allocate("test")
	err = ctx.SetParamIfAllocated(param("test"), "OK")
	if err != nil {
		t.Errorf("SetParamIfAllocated should pass for an allocated param, got %s", err)
	}
	if val, ok := ctx.vals[param("test")]; !ok || val != "OK" {
		t.Errorf("Expected ctx.vals[param('test')] = 'OK', true, got %s, %t", val, ok)
	}
}

func TestDeadline(t *testing.T) {
	ctx := newRMC()
	if ctxdl, ok := ctx.Deadline(); ok {
		t.Errorf("expected ctx.Deadline to fail with no base; got %v", ctxdl)
	}
	dl := time.Now().Add(time.Second * 10)
	base, cancel := context.WithDeadline(context.Background(), dl)
	defer cancel()
	ctx.ResetOnto(base)
	if ctxdl, ok := ctx.Deadline(); !ok {
		t.Errorf("expected ctx.Deadline to be OK if base has deadline")
	} else if !ctxdl.Equal(dl) {
		t.Errorf("expected ctx.Deadline to be equal to base deadline %v; got %v", dl, ctxdl)
	}
}

func TestDone(t *testing.T) {
	ctx := newRMC()
	if ctxd := ctx.Done(); ctxd != nil {
		t.Errorf("expected ctx.Done to fail with no base")
	}
	base, cancel := context.WithCancel(context.Background())
	ctx.ResetOnto(base)
	ctxd := ctx.Done()
	if ctxd == nil {
		t.Fatalf("expected ctx.Done to have a value with cancellable but uncancelled base")
	}
	cancel()
	select {
	case <-ctxd:
	default:
		t.Errorf("Expected ctxd to send after cancellation")
	}
}

func TestErr(t *testing.T) {
	ctx := newRMC()
	ctx.err = errors.New("test err")
	if err := ctx.Err(); err == nil || !errors.Is(err, ctx.err) {
		t.Errorf("ctx.Err not properly returning own error")
	}
	ctx.err = nil
	if err := ctx.Err(); err != nil {
		t.Errorf("ctx.Err should return nil with nil error, got %s", err)
	}
	base := context.Background()
	ctx.ResetOnto(base)
	if err := ctx.Err(); err != nil {
		t.Errorf("ctx.Err should return nil with nil and base with nil, got %s", err)
	}
}

func TestValue(t *testing.T) {
	ctx := newRMC()
	ctx.Allocate("test")
	base := context.WithValue(context.Background(), string("pass"), string("OK"))
	ctx.ResetOnto(base)
	err := ctx.SetParamIfAllocated(param("test"), "OK")
	if err != nil {
		t.Fatal(err)
	}
	if v := ctx.Value(param("test")); v == nil {
		t.Errorf("ctx should have value for param('test')")
	} else if v.(string) != "OK" {
		t.Errorf("ctx.Value(param('test')) should be 'OK', got '%v'", v)
	}
	if v := ctx.Value(string("pass")); v == nil {
		t.Errorf("ctx should have value for string ('pass') by pass-through to base")
	} else if v.(string) != "OK" {
		t.Errorf("ctx.Value(string('pass')) should be 'OK', got '%v'", v)
	}
}

func TestGetSetParam(t *testing.T) {
	ctx := newRMC()
	ctx.Allocate("test")
	if param := GetParam(ctx, "test"); param != "" {
		t.Errorf("expected empty param 'test', got '%s'", param)
	}
	setParam(ctx, "test", "OK")
	if param := GetParam(ctx, "test"); param == "" || param != "OK" {
		t.Errorf("expected param 'test' == 'OK', got '%s'", param)
	}
	base := context.WithValue(ctx, "ignore", "ignore")
	if param := GetParam(base, "test"); param != "OK" {
		t.Errorf("expected plain context to pass through to rmc for 'test' == 'ok', got '%s'", param)
	}
	ctx.ResetOnto(nil)
	if param := GetParam(ctx, "test"); param != "" {
		t.Errorf("expected empty param 'test', got '%s'", param)
	}
	base = context.WithValue(ctx, "ignore", "ignore")
	if param := GetParam(base, "test"); param != "" {
		t.Errorf("expected plain context to pass through to rmc for 'test' == '', got '%s'", param)
	}
	base = context.Background()
	if param := GetParam(base, "test"); param != "" {
		t.Errorf("expected plain context with no pass through to give 'test' == '', got '%s'", param)
	}
}
