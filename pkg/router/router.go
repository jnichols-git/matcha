package router

import (
	"net/http"

	"github.com/cloudretic/router/pkg/middleware"
	"github.com/cloudretic/router/pkg/route"
)

type Router interface {
	Attach(mw middleware.Middleware)
	AddRoute(r route.Route, h http.Handler)
	AddNotFound(h http.Handler)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

func New(cfs ...ConfigFunc) (Router, error) {
	rt := Default()
	for _, cf := range cfs {
		err := cf(rt)
		if err != nil {
			return nil, err
		}
	}
	return rt, nil
}

func Declare(cfs ...ConfigFunc) Router {
	rt := Default()
	for _, cf := range cfs {
		err := cf(rt)
		if err != nil {
			panic(err)
		}
	}
	return rt
}
