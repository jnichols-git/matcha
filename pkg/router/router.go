package router

import (
	"net/http"

	"github.com/CloudRETIC/router/pkg/route"
)

type Router interface {
	AddRoute(r route.Route, h http.Handler)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

func Handle(rt Router, r route.Route, handler http.Handler) error {
	rt.AddRoute(r, handler)
	return nil
}
