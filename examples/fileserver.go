package examples

import (
	"net/http"
	"os"

	"github.com/jnichols-git/matcha/v2/internal/rctx"
	"github.com/jnichols-git/matcha/v2/internal/router"
)

type fileServer struct {
	root string
}

func (fs *fileServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := rctx.GetParam(req.Context(), "filepath")
	dat, err := os.ReadFile(fs.root + path)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("File " + path + " does not exist."))
		return
	}
	w.Write(dat)
}

func FileServer(dir string) {
	rt := router.Default()
	rt.Handle(http.MethodGet, "/files/{filepath}+", &fileServer{dir})
	http.ListenAndServe(":3000", rt)
}
