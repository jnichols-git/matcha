package route

import (
	"errors"
	"regexp"

	"github.com/jnichols-git/matcha/v2/regex"
)

const (
	regexp_part = string(`^(?:(\/[\w\.\~]*)|\/(\:\w+)?(\[.+\])?(\+)?)$`)
)

var regexp_part_compiled *regexp.Regexp = regexp.MustCompile(regexp_part)

type Part struct {
	expr    string
	param   string
	pattern *regexp.Regexp
	multi   bool
}

// Match checks if a Part matches an input token.
// Tokens should not contain a leading or trailing forward-slash.
func (p Part) Match(token string) (matched bool) {
	clean := token[1:]
	matched = (p.expr != "" && p.expr == token) ||
		(p.pattern != nil && p.pattern.FindString(clean) == clean) ||
		(p.expr == "" && p.pattern == nil)
	return
}

// Eq checks if a Part is equal to another Part.
// This is used in tree-building for routers.
func (p Part) Eq(other Part) (eq bool) {
	eq = p.expr == other.expr &&
		p.param == other.param &&
		p.pattern != nil && other.pattern != nil && p.pattern.String() == other.pattern.String()
	return
}

func (p Part) Nil() (isNil bool) {
	isNil = p.expr == "" && p.param == "" && p.pattern == nil
	return
}

// Parameter returns the name of a URI parameter attached to this Part, if there is one.
func (p Part) Parameter() string {
	return p.param
}

// Multi returns true if the part should be permitted to match multiple route tokens.
func (p Part) Multi() bool {
	return p.multi
}

func ParsePart(token string) (p Part, err error) {
	groups := regex.Groups(regexp_part_compiled, token)
	if groups == nil {
		err = errors.New("provided Part expression " + token + " does not match regex " + regexp_part)
		return
	}
	for _, group := range groups {
		switch group[0] {
		case ':':
			p.param = group[1:]
		case '[':
			p.pattern, err = regexp.Compile(group[1 : len(group)-1])
		case '+':
			p.multi = true
		default:
			p.expr = group
		}
	}
	if p.expr != "" && (p.param != "" || p.pattern != nil || p.multi) {
		err = errors.Join(err, errors.New("Parts cannot have an expression and a param/pattern/multi modifier"))
	}
	return
}
