package matcha

import (
	"github.com/jnichols-git/matcha/v2/internal/route"
	"github.com/jnichols-git/matcha/v2/internal/router"
)

func Route(method, expr string) (r *route.Route, err error) {
	return route.New(method, expr)
}

func Router() (r *router.Router) {
	return router.Default()
}
