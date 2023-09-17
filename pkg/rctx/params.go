package rctx

import (
	"errors"
)

type paramKey string

type routeParam struct {
	key   paramKey
	value string
}

type routeParams struct {
	rps      []routeParam
	reserved reservedParams
	cap      int
	head     int
}

// PARAMETERS

func newParams(size int) *routeParams {
	rps := &routeParams{
		rps:      make([]routeParam, size),
		reserved: reservedParams{},
		cap:      size,
		head:     0,
	}
	for i := 0; i < size; i++ {
		rps.rps[i] = routeParam{}
	}
	return rps
}

func (rps *routeParams) get(key paramKey) string {
	if value, reserved := rps.reserved.get(key); reserved {
		return value
	}
	for i := 0; i < rps.head; i++ {
		kv := rps.rps[i]
		if kv.key == key {
			return kv.value
		}
	}
	return ""
}

func (rps *routeParams) set(in *Context, key paramKey, value string) error {
	if reserved, err := rps.reserved.set(in.parent, key, value); reserved {
		return err
	}
	idx := rps.head
	inc := true
	for i := 0; i < rps.head; i++ {
		kv := rps.rps[i]
		if kv.key == key {
			idx = i
			inc = false
		}
	}
	if idx >= rps.cap {
		return errors.New("over capacity")
	}
	if inc {
		rps.head++
	}
	rps.rps[idx].key = key
	rps.rps[idx].value = value
	return nil
}
