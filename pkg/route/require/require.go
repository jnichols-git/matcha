package require

import (
	"net/http"
)

type Required func(req *http.Request) bool

// ExecuteRequireds executes a list of route validators on a request.
// It only returns true if every validator provided returns true.
func Execute(req *http.Request, vs []Required) bool {
	for _, v := range vs {
		if !v(req) {
			return false
		}
	}
	return true
}
