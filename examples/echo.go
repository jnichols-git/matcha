package examples

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2"
	"github.com/jnichols-git/matcha/v2/internal/rctx"
)

func echoAdmin(w http.ResponseWriter, req *http.Request) {
	name := rctx.GetParam(req.Context(), "name")
	w.Write([]byte("Hello, admin " + name + "!"))
}

func echo(w http.ResponseWriter, req *http.Request) {
	name := rctx.GetParam(req.Context(), "name")
	w.Write([]byte("Hello, " + name + "!"))
}

func EchoExample() {
	rt := matcha.Router()
	rt.HandleFunc(http.MethodGet, "/hello/:name{admin:.+}", echoAdmin)
	rt.HandleFunc(http.MethodGet, "/hello/:name", echo)
	http.ListenAndServe(":3000", rt)
}
