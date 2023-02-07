package router

import (
	"net/http"

	"github.com/cloudretic/router/pkg/middleware"
	"github.com/cloudretic/router/pkg/route"
)

type Router interface {
	Attach(mw middleware.Middleware)
	AddRoute(r route.Route, h http.Handler)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

func Handle(rt Router, r route.Route, handler http.Handler) error {
	rt.AddRoute(r, handler)
	return nil
}
