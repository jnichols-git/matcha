package cors

import (
	"fmt"
	"net/http"

	"github.com/cloudretic/router/pkg/route"
)

// Generate a preflight request handler for a route expression.
func PreflightHandler(expr string, aco *AccessControlOptions) (route.Route, http.HandlerFunc, error) {
	r, err := route.New(http.MethodOptions, expr)
	if err != nil {
		return nil, nil, err
	}
	f := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hello")
		SetCORSResponseHeaders(w, r, aco)
		w.WriteHeader(http.StatusNoContent)
	}
	return r, f, nil
}
