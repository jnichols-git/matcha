package route

import (
	"net/http"
	"strings"
)

type Validator func(req *http.Request) bool

func ExecuteValidators(req *http.Request, vs []Validator) bool {
	for _, v := range vs {
		if !v(req) {
			return false
		}
	}
	return true
}

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
