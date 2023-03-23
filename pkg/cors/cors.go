package cors

import (
	"net/http"
	"strconv"
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

// Access-Control request headers
// Can be pulled from requests to check against access-control
type CORSRequestHeaders struct {
	Origin         string
	RequestMethod  string
	RequestHeaders []string
	Cookie         string
}

// Get CORS request headers from an HTTP request.
func GetCORSRequestHeaders(req *http.Request) (crh *CORSRequestHeaders) {
	crh = &CORSRequestHeaders{}
	crh.Origin = req.Header.Get(Origin)
	crh.RequestMethod = req.Header.Get(RequestMethod)
	if len(crh.RequestMethod) == 0 {
		crh.RequestMethod = req.Method
	}
	crh.RequestHeaders = req.Header.Values(RequestHeaders)
	return
}

// Access-Control options
// Used to control cross-origin resource sharing on routes
type AccessControlOptions struct {
	AllowOrigin      []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	MaxAge           float64
	AllowCredentials bool
}

// Create a deep copy of aco that reflects the headers in crh.
// If there is no option in aco that can reflect crh, the output will have an empty field; this is intended behavior
// and indicates to the HTTP client that resource sharing should not be allowed for this request
func ReflectCORSRequestHeaders(aco *AccessControlOptions, crh *CORSRequestHeaders) *AccessControlOptions {
	out := &AccessControlOptions{
		AllowOrigin:      make([]string, 1),
		AllowMethods:     make([]string, 1),
		AllowHeaders:     make([]string, len(crh.RequestHeaders)),
		ExposeHeaders:    make([]string, len(aco.ExposeHeaders)),
		MaxAge:           0,
		AllowCredentials: aco.AllowCredentials,
	}
	if len(aco.AllowOrigin) == 1 && aco.AllowOrigin[0] == "*" {
		out.AllowOrigin = []string{crh.Origin}
	} else {
		for _, allowedOrigin := range aco.AllowOrigin {
			if crh.Origin == allowedOrigin {
				out.AllowOrigin = []string{crh.Origin}
				break
			}
		}
	}
	if len(aco.AllowMethods) == 1 && aco.AllowMethods[0] == "*" {
		out.AllowMethods = []string{crh.RequestMethod}
	} else {
		for _, allowedMethod := range aco.AllowMethods {
			if crh.RequestMethod == allowedMethod {
				out.AllowMethods = []string{crh.RequestMethod}
				break
			}
		}
	}
	if len(aco.AllowHeaders) == 1 && aco.AllowHeaders[0] == "*" {
		for i, requestedHeader := range crh.RequestHeaders {
			out.AllowHeaders[i] = requestedHeader
		}
	} else {
		hct := 0
	allowed:
		for _, allowedHeader := range aco.AllowHeaders {
			for _, requestedHeader := range crh.RequestHeaders {
				if allowedHeader == requestedHeader {
					out.AllowHeaders[hct] = requestedHeader
					hct++
					continue allowed
				}
			}
		}
		out.AllowHeaders = out.AllowHeaders[:hct]
	}
	// There's not a great way to check which headers need to be exposed, so this is returned as * if that's provided.
	for i, exposedHeader := range aco.ExposeHeaders {
		out.ExposeHeaders[i] = exposedHeader
	}
	return out
}

func SetCORSResponseHeaders(w http.ResponseWriter, req *http.Request, aco *AccessControlOptions) {
	crh := GetCORSRequestHeaders(req)
	res := ReflectCORSRequestHeaders(aco, crh)
	h := w.Header()
	h.Del(AllowOrigin)
	h.Del(AllowMethods)
	h.Del(AllowHeaders)
	h.Del(ExposeHeaders)
	for _, origin := range res.AllowOrigin {
		h.Add(AllowOrigin, origin)
	}
	for _, method := range res.AllowMethods {
		h.Add(AllowMethods, method)
	}
	for _, header := range res.AllowHeaders {
		h.Add(AllowHeaders, header)
	}
	for _, header := range res.ExposeHeaders {
		h.Add(ExposeHeaders, header)
	}
	h.Set(MaxAge, strconv.FormatFloat(aco.MaxAge, 'f', 0, 64))
	h.Set(AllowCredentials, strconv.FormatBool(aco.AllowCredentials))
}
