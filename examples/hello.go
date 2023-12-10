package main

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2"
)

func sayHello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello, World!\n"))
}

func HelloExample() {
	rt := matcha.Router()
	rt.HandleFunc(http.MethodGet, "/hello", sayHello)
	http.ListenAndServe(":3000", rt.Handler())
}
