package matcha

import "github.com/jnichols-git/matcha/v2/internal/route"

func Route(method, expr string) (r *route.Route, err error) {
	return route.New(method, expr)
}
