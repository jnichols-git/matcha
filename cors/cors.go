// Package cors defines a set of useful handlers and middleware for setting access control values.
//
// See [https://github.com/jnichols-git/matcha/v2/blob/main/docs/other-features.md#cross-origin-resource-sharing-cors].
package cors

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	Origin           = string("Origin")
	RequestMethod    = string("Access-Control-Request-Method")
	RequestHeaders   = string("Access-Control-Request-Headers")
	AllowOrigin      = string("Access-Control-Allow-Origin")
	AllowMethods     = string("Access-Control-Allow-Methods")
	AllowHeaders     = string("Access-Control-Allow-Headers")
	ExposeHeaders    = string("Access-Control-Expose-Headers")
	MaxAge           = string("Access-Control-Max-Age")
	AllowCredentials = string("Access-Control-Allow-Credentials")
)

// Request defines the CORS-related fields extracted from an *http.Request.
type Request struct {
	Origin         string
	RequestMethod  string
	RequestHeaders []string
}

// Get CORS request headers from an HTTP request.
func GetRequest(req *http.Request) (crh *Request) {
	crh = &Request{}
	crh.Origin = req.Header.Get(Origin)
	crh.RequestMethod = req.Header.Get(RequestMethod)
	if len(crh.RequestMethod) == 0 {
		crh.RequestMethod = req.Method
	}
	crh.RequestHeaders = req.Header.Values(RequestHeaders)
	return
}

// Options define the set of options that a CORS request may return.
type Options struct {
	AllowOrigin      []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	MaxAge           float64
	AllowCredentials bool
}

// ReflectRequest sets headers in out depending on the provided Options and
// Request.
// To limit the cost of OPTIONS requests, ReflectRequest sends the minimum
// possible permissions.
func ReflectRequest(aco *Options, crh *Request, out http.Header) {
	/*
		out := &Options{
			AllowHeaders:     make([]string, len(crh.RequestHeaders)),
			ExposeHeaders:    make([]string, len(aco.ExposeHeaders)),
			MaxAge:           0,
			AllowCredentials: aco.AllowCredentials,
		}
	*/
	if len(aco.AllowOrigin) == 1 && aco.AllowOrigin[0] == "*" {
		out.Set(AllowOrigin, crh.Origin)
	} else {
		for _, allowedOrigin := range aco.AllowOrigin {
			if crh.Origin == allowedOrigin {
				out.Set(AllowOrigin, crh.Origin)
				break
			}
		}
	}
	if len(aco.AllowMethods) == 1 && aco.AllowMethods[0] == "*" {
		out.Set(AllowMethods, crh.RequestMethod)
	} else {
		for _, allowedMethod := range aco.AllowMethods {
			if crh.RequestMethod == allowedMethod {
				out.Set(AllowMethods, crh.RequestMethod)
				break
			}
		}
	}
	if len(aco.AllowHeaders) == 1 && aco.AllowHeaders[0] == "*" {
		for _, requestedHeader := range crh.RequestHeaders {
			out.Add(AllowHeaders, requestedHeader)
		}
	} else {
	allowed:
		for _, allowedHeader := range aco.AllowHeaders {
			for _, requestedHeader := range crh.RequestHeaders {
				if strings.EqualFold(allowedHeader, requestedHeader) {
					out.Add(AllowHeaders, requestedHeader)
					continue allowed
				}
			}
		}
	}
	// There's not a great way to check which headers need to be exposed, so this is returned as * if that's provided.
	for _, exposedHeader := range aco.ExposeHeaders {
		out.Add(ExposeHeaders, exposedHeader)
	}
	out.Set(MaxAge, strconv.FormatFloat(aco.MaxAge, 'f', 0, 64))
	out.Set(AllowCredentials, strconv.FormatBool(aco.AllowCredentials))
}

// Updates the response headers from http.ResponseWriter to mirror a set of access control options.
// Mirroring provides the minimum amount of permissions needed for the inbound request via ReflectCorsRequestHeaders.
func SetCORSResponseHeaders(w http.ResponseWriter, req *http.Request, aco *Options) {
	crh := GetRequest(req)
	h := w.Header()
	h.Del(AllowOrigin)
	h.Del(AllowMethods)
	h.Del(AllowHeaders)
	h.Del(ExposeHeaders)
	ReflectRequest(aco, crh, h)
}
