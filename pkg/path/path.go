// package path manages the tokenization of http paths
// TODO: this should belong to the data repo
package path

import (
	"strings"

	"github.com/cloudretic/go-collections/pkg/slices"
)

// Tokenize a string over / in-order
// Filters out any blank results
func TokenizeString(s string) []string {
	tokens := strings.Split(s, "/")
	tokens, _ = slices.FilterFunc(tokens, func(token string) bool { return token == "" })
	return tokens
}
