package matcha

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2/internal/rctx"
	"github.com/jnichols-git/matcha/v2/route"
	"github.com/jnichols-git/matcha/v2/router"
)

func Route(method, expr string) (r *route.Route, err error) {
	return route.New(method, expr)
}

func Router() (r *router.Router) {
	return router.Default()
}

func RouteParam(req *http.Request, param string) string {
	return rctx.GetParam(req.Context(), param)
}
