package route

import (
	"net/http"
)

// routeRequestHandler contains information for a service and the handler to use.
type RequestHandler struct {
	w          http.ResponseWriter
	req        *http.Request
	useHandler http.Handler
}

func (rrh *RequestHandler) Serve() {
	rrh.useHandler.ServeHTTP(rrh.w, rrh.req)
}

// Use a route sequentially.
// Returns the routeRequestHandler for this route, if the request matches.
func UseRoute(r Route, w http.ResponseWriter, req *http.Request) *RequestHandler {
	matched := r.MatchAndUpdateContext(req)
	if matched == nil {
		return nil
	}
	return &RequestHandler{
		w:          w,
		req:        matched,
		useHandler: r.Handler(),
	}
}
