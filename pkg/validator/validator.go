// Package validator defines the ValidatorFunc, which allows switching of routes based on request properties.
// This serves a distinct purpose compared to middleware; it can't modify requests or send responses, and can only
// return control to the router by rejecting requests with a false return value.
package validator

import (
	"net/http"
	"strings"
)

type Validator func(req *http.Request) bool

func Hosts(hn ...string) Validator {
	return func(req *http.Request) bool {
		rh := strings.Split(req.Host, ":")[0]
		if rh == "" {
			return false
		}
		for _, h := range hn {
			if rh == h {
				return true
			}
		}
		return false
	}
}
