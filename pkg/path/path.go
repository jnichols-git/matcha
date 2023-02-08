// package path manages the tokenization of http paths
// TODO: this should belong to the data repo
package path

import (
	"strings"
)

// Return the next token from a path, starting at position last, and the position to use with the next call.
// Returns "" if no token could be found
func Next(path string, last int) (string, int) {
	end := strings.Index(path[last+1:], "/")
	if end == -1 {
		return path[last:], -1
	} else {
		return path[last : last+end+1], last + end + 1
	}
}
