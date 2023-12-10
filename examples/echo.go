package main

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2"
)

func echo(w http.ResponseWriter, req *http.Request) {
	name := matcha.RouteParam(req, "name")
	w.Write([]byte("Hello, " + name + "!\n"))
}

func EchoExample() {
	rt := matcha.Router()
	rt.HandleFunc(http.MethodGet, "/hello/:name", echo)
	http.ListenAndServe(":3000", rt.Handler())
}
