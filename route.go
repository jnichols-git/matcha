package matcha

import "github.com/jnichols-git/matcha/v2/internal/route"

type Route struct {
	method  string
	expr    string
	parts   []route.Part
	partial bool
}
