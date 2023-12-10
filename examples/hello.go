package examples

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2"
)

func sayHello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func HelloExample() {
	rt := matcha.Router()
	rt.HandleFunc(http.MethodGet, "/hello", sayHello)
	// or:
	// rt.Handle(http.MethodGet, "/hello", http.HandlerFunc(sayHello))
	http.ListenAndServe(":3000", rt)
}
