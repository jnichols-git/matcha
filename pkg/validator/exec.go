package validator

import "net/http"

func ExecuteValidators(req *http.Request, vs []Validator) bool {
	for _, v := range vs {
		if !v(req) {
			return false
		}
	}
	return true
}
