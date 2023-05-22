package route

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudretic/matcha/pkg/regex"
)

type Validator func(req *http.Request) bool

// ExecuteValidators executes a list of route validators on a request.
// It only returns true if every validator provided returns true.
func ExecuteValidators(req *http.Request, vs []Validator) bool {
	for _, v := range vs {
		if !v(req) {
			return false
		}
	}
	return true
}

// getReqHost gets the host, port for an inbound server request.
// If no port is detected, this will make a best-effort guess at port 80 or 443 based on information
// about TLS encryption on the request; if a TLS handshake was completed, the port will be 443, and vice versa.
// This shouldn't be used for parsing anything but inbound requests.
func getReqHost(req *http.Request) (string, string) {
	var host, port string
	toks := strings.Split(req.Host, ":")
	if len(toks) == 2 {
		host, port = toks[0], toks[1]
	} else if len(toks) == 1 {
		if req.TLS.HandshakeComplete {
			port = "443"
		} else {
			port = "80"
		}
		host = toks[0]
	} else {
		return "", ""
	}
	return host, port
}

// Hosts checks a request against a list of host patterns.
func Hosts(hns ...string) Validator {
	hmfs := make([]Validator, 0, len(hns))
	for _, hn := range hns {
		toks := strings.Split(hn, ":")
		var hf, pf func(str string) bool
		if len(toks) == 0 {
			continue
		}
		if len(toks) >= 1 {
			h := toks[0]
			hpatt, isPatt, err := regex.CompilePattern(h)
			hf = func(host string) bool {
				if !isPatt || err != nil {
					return host == h
				} else {
					return hpatt.Match(host)
				}
			}
		}
		if len(toks) >= 2 {
			p := toks[1]
			ppatt, isPatt, err := regex.CompilePattern(p)
			pf = func(port string) bool {
				if !isPatt || err != nil {
					return port == p
				} else {
					return ppatt.Match(port)
				}
			}
		} else {
			pf = func(str string) bool {
				p, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return false
				}
				return 1 <= p && p <= 65535
			}
		}
		hmfs = append(hmfs, func(req *http.Request) bool {
			host, port := getReqHost(req)
			return hf(host) && pf(port)
		})
	}
	return func(req *http.Request) bool {
		for _, v := range hmfs {
			if !v(req) {
				return false
			}
		}
		return true
	}
}
