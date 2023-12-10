package main

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2"
)

func ValidateName(w http.ResponseWriter, req *http.Request) *http.Request {
	if name := matcha.RouteParam(req, "name"); name[0] != 'A' {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Names must start with 'A'.\n"))
		return nil
	}
	return req
}

func MiddlewareExample() {
	router := matcha.Router()
	nameRoute, _ := matcha.Route(http.MethodGet, "/hello/:name")
	nameRoute.Use(ValidateName)
	router.HandleRouteFunc(nameRoute, echo)
	http.ListenAndServe(":3000", router)
}
