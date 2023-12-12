package main

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2"
)

func ValidateName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if name := matcha.RouteParam(r, "name"); name[0] != 'A' {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Names must start with 'A'.\n"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func MiddlewareExample() {
	router := matcha.Router()
	nameRoute, _ := matcha.Route(http.MethodGet, "/hello/:name")
	nameRoute.Use(ValidateName)
	router.HandleRouteFunc(nameRoute, echo)
	http.ListenAndServe(":3000", router)
}
