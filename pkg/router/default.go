package router

import (
	"net/http"

	"github.com/CloudRETIC/router/pkg/route"
)

type defaultRouter struct {
	routes route.RouteSet
}

func Default() *defaultRouter {
	return &defaultRouter{}
}

func (rt *defaultRouter) AddRoute(r route.Route) {
	rt.routes = append(rt.routes, r)
}

func (handler *defaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, r := range handler.routes {
		handler := route.UseRoute(r, w, req)
		if handler != nil {
			handler.Serve()
			return
		}
	}
}
