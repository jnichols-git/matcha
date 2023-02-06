package route

import (
	"net/http"
	"regexp"

	"github.com/CloudRETIC/router/pkg/router/params"
)

// LITERALS

// Literal route Parts match and pass without additionally transforming
// the request.

// stringPart; literal string part
type stringPart struct {
	val string
}

func build_stringPart(val string) (*stringPart, error) {
	return &stringPart{val}, nil
}

func (part *stringPart) Match(req *http.Request, token string) *http.Request {
	if part.val == token {
		return req
	} else {
		return nil
	}
}

// WILDCARDS

// Wildcard route Parts store parameters for use by the router in handlers.
// They use a syntax of [wildcard] to denote their name, and can additionally
// be qualified by some conditions by splitting with the : character.

// wildcardParts always match, and add the token as a request param.
type wildcardPart struct {
	param string
}

func build_wildcardPart(param string) (*wildcardPart, error) {
	return &wildcardPart{param}, nil
}

func (part *wildcardPart) Match(req *http.Request, token string) *http.Request {
	if part.param != "" {
		req = params.Set(req, part.param, token)
	}
	return req
}

// regexParts match against regular expressions.
// They're created using the syntax [wildcard]:{regex}
type regexPart struct {
	param string
	expr  *regexp.Regexp
}

func build_regexPart(param, expr string) (*regexPart, error) {
	expr_compiled, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	} else {
		return &regexPart{param, expr_compiled}, nil
	}
}

func (part *regexPart) Match(req *http.Request, token string) *http.Request {
	// Match against regex
	matched := part.expr.FindString(token)
	if matched == "" {
		return nil
	}
	// If a parameter is set, act as a wildcard param.
	if part.param != "" {
		// If a token matched, store the matched value as a route Param
		req = params.Set(req, part.param, matched)
	}
	return req
}
