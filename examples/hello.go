package examples

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2/internal/router"
)

func sayHello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func HelloExample() {
	rt := router.Default()
	rt.HandleFunc(http.MethodGet, "/hello", sayHello)
	// or:
	// rt.Handle(http.MethodGet, "/hello", http.HandlerFunc(sayHello))
	http.ListenAndServe(":3000", rt)
}
