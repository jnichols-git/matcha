package router

import (
	"net/http"

	"github.com/CloudRETIC/router/pkg/route"
)

type Router interface {
	AddRoute(r route.Route)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

func Handle(rt Router, expr string, handler http.Handler) error {
	r, err := route.New(expr, handler)
	if err != nil {
		return err
	}
	rt.AddRoute(r)
	return nil
}
