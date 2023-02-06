// package path manages the tokenization of http paths
// TODO: this should belong to the data repo
package path

import (
	"strings"
)

// Tokenize a string over / in-order
// Filters out any blank results
func TokenizeString(s string) []string {
	tokens := strings.Split(s, "/")
	out := make([]string, 0)
	for _, token := range tokens {
		if token != "" {
			out = append(out, token)
		}
	}
	return out
}
