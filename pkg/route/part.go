package route

import (
	"context"
	"fmt"
	"regexp"

	"github.com/cloudretic/matcha/pkg/regex"
)

const (
	// part matching
	regexp_wildcard = string(`/\[(.*?)\](.*)`)
	regexp_regex    = string(`[/\]]{(.*)}`)
	// handy constants to have around
	regexp_anyWord = string(`\w+`)
)

// Regex used for parsing tokens
var regexp_wildcard_compiled *regexp.Regexp = regexp.MustCompile(regexp_wildcard)
var regexp_regex_compiled *regexp.Regexp = regexp.MustCompile(regexp_regex)

// Compiled regex for any word
var regexp_anyWord_compiled *regexp.Regexp = regexp.MustCompile(regexp_anyWord)

// Parts are the main body of a Route, and are an interface defining
// a Match function against tokens in a request URL.
type Part interface {
	// Match should return nil if the Part doesn't match the token.
	// If it does, it should return the request, with any modifications done on
	// behalf of the Part (usually wildcard tokens)
	Match(ctx context.Context, token string) bool
	// Compare to another part.
	// Should return equal iff the result of Match would be the exact same, given the same context and token.
	Eq(other Part) bool
}

// paramParts may or may not store some parameter.
// This is for internal use in package route only, so that extensions of Part/Route can specialize behavior
// for Parts that do or don't have parameters.
type paramPart interface {
	ParameterName() string
	SetParameterName(string)
}

// Parse a token into a route Part.
func parse(token string) (Part, error) {
	// wildcard check
	if groups := regex.Groups(regexp_wildcard_compiled, token); groups != nil {
		// There must be at least one group here.
		wildcardExpr := groups[0]
		// If there's another group, we need to specialize further.
		// Otherwise, it's a regular wildcardPart.
		if len(groups) > 1 {
			// regex check
			if groups := regex.Groups(regexp_regex_compiled, token); groups != nil {
				regexExpr := groups[0]
				return build_regexPart(wildcardExpr, regexExpr)
			}
		}
		if len(wildcardExpr)+3 != len(token) {
			return nil, fmt.Errorf("error parsing expression %s: got a wildcard part with a non-regex addition, which is invalid", token)
		}
		return build_wildcardPart(wildcardExpr)
	}

	// If we get here, it's not a wildcard.

	// regex check
	if groups := regex.Groups(regexp_regex_compiled, token); groups != nil {
		regexExpr := groups[0]
		return build_regexPart("", regexExpr)
	}

	// Not wildcard or regex; just return as stringPart
	return build_stringPart(token)
}
