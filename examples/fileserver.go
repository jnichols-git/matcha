package main

import (
	"net/http"
	"os"

	"github.com/jnichols-git/matcha/v2"
)

type fileServer struct {
	root string
}

func (fs *fileServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := matcha.RouteParam(req, "filepath")
	dat, err := os.ReadFile(fs.root + path)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("File " + path + " does not exist.\n"))
		return
	}
	w.Write(dat)
}

func FileServerExample() {
	rt := matcha.Router()
	rt.Handle(http.MethodGet, "/files/:filepath+", &fileServer{"./examples/"})
	http.ListenAndServe(":3000", rt.Handler())
}
