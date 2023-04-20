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
	rps  []*routeParam
	cap  int
	head int
}

// PARAMETERS

func newParams(size int) *routeParams {
	rps := &routeParams{
		rps:  make([]*routeParam, size),
		cap:  size,
		head: 0,
	}
	for i := 0; i < size; i++ {
		rps.rps[i] = &routeParam{}
	}
	return rps
}

func (rps *routeParams) get(key paramKey) string {
	for i := rps.head; i > 0; i-- {
		kv := rps.rps[i-1]
		if kv.key == key {
			return kv.value
		}
	}
	return ""
}

func (rps *routeParams) set(key paramKey, value string) error {
	if rps.head >= rps.cap {
		return errors.New("placeholder error; over capacity")
	}
	rps.rps[rps.head].key = key
	rps.rps[rps.head].value = value
	rps.head++
	return nil
}
