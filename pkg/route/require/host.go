package require

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudretic/matcha/pkg/regex"
)

// getReqHost gets the host, port for an inbound server request.
// If no port is detected, this will make a best-effort guess at port 80 or 443 based on information
// about TLS encryption on the request; if the scheme is https, port 443 will be used, and vice versa.
func getReqHost(req *http.Request) (string, string) {
	var host, port string
	toks := strings.Split(req.Host, ":")
	if len(toks) == 2 {
		host, port = toks[0], toks[1]
	} else if len(toks) == 1 {
		if req.URL.Scheme == "https" {
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

// splitHostPort gets the scheme, host, port for a value in Hosts or HostPorts.
// It provides the tokens that represent each component. See the documentation for HostPorts for the
// most detail.
func splitHostPort(hp string) (scheme, host, port string) {
	toks := strings.Split(hp, "://")
	if len(toks) == 2 {
		scheme = toks[0]
		hp = toks[1]
	}
	toks = strings.Split(hp, ":")
	if len(toks) == 2 {
		host = toks[0]
		port = toks[1]
	} else {
		host = hp
	}
	return
}

// Hosts checks a request against a list of host patterns.
// Hosts should be provided as a string or pattern (see package regex), and information about scheme and ports
// will be ignored; if you want those features, see HostPorts.
func Hosts(hns ...string) Required {
	hmfs := make([]Required, 0, len(hns))
	for _, hn := range hns {
		var hf func(str string) bool
		_, host, _ := splitHostPort(hn)
		hpatt, isPatt, err := regex.CompilePattern(host)
		hf = func(inHost string) bool {
			_, inHost, _ = splitHostPort(inHost)
			if !isPatt || err != nil {
				return inHost == host
			} else {
				return hpatt.Match(inHost)
			}
		}
		hmfs = append(hmfs, func(req *http.Request) bool {
			return hf(req.Host)
		})
	}
	return func(req *http.Request) bool {
		for _, v := range hmfs {
			if v(req) {
				return true
			}
		}
		return false
	}
}

// HostPorts checks a request against a list of host:port patterns.
// Hosts should be provided in the format "scheme://host:port", where:
//   - Scheme is either the string http or https (optional, case-sensitive)
//   - Host is either a string or a Pattern (see package regex)
//   - Port is one or more numbers or ranges of numbers (ex. 1-10), delimited by commas (optional)
//
// When provided with a port, Hosts will always match against those ports, regardless of the scheme.
// Ranges are inclusive on both ends. Otherwise, it will match ports 80 and 443 for HTTP and HTTPS
// respectively, defaulting to HTTP. You must provide a scheme or port number if you want to match HTTPS
// requests to a specific host.
func HostPorts(hns ...string) Required {
	hmfs := make([]Required, 0, len(hns))
	for _, hn := range hns {
		var hf, pf func(str string) bool
		scheme, host, port := splitHostPort(hn)
		hpatt, isPatt, err := regex.CompilePattern(host)
		hf = func(inHost string) bool {
			if !isPatt || err != nil {
				return inHost == host
			} else {
				return hpatt.Match(host)
			}
		}
		if port != "" {
			ps := strings.Split(port, ",")
			pfs := make([]func(inPort int64) bool, 0, len(ps))
			for _, p := range ps {
				if prange := strings.Split(p, "-"); len(prange) == 2 {
					start, err := strconv.ParseInt(prange[0], 10, 64)
					if err != nil {
						continue
					}
					end, err := strconv.ParseInt(prange[1], 10, 64)
					if err != nil {
						continue
					}
					pfs = append(pfs, func(inPort int64) bool {
						return start <= inPort && inPort <= end
					})
				} else {
					target, err := strconv.ParseInt(p, 10, 64)
					if err != nil {
						continue
					}
					pfs = append(pfs, func(inPort int64) bool {
						return inPort == target
					})
				}
			}
			pf = func(str string) bool {
				pn, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return false
				}
				for _, f := range pfs {
					if f(pn) {
						return true
					}
				}
				return false
			}
		} else {
			p := 80
			if scheme == "https" {
				p = 443
			}
			pf = func(str string) bool {
				pn, err := strconv.ParseInt(str, 10, 64)
				if err != nil {
					return false
				}
				return int(pn) == p
			}
		}
		hmfs = append(hmfs, func(req *http.Request) bool {
			host, port := getReqHost(req)
			return hf(host) && pf(port)
		})
	}
	return func(req *http.Request) bool {
		for _, v := range hmfs {
			if v(req) {
				return true
			}
		}
		return false
	}
}
