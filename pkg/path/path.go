// package path manages the tokenization of http paths
// TODO: this should belong to the data repo
package path

import (
	"strings"
)

// Return the next token from a path, starting at position last, and the position to use with the next call.
// Returns "" if no token could be found
func Next(path string, last int) (string, int) {
	if last+1 > len(path) {
		return "", -1
	}
	start := last
	for {
		idx := strings.Index(path[start+1:], "/")
		if idx == -1 {
			break
		}
		end := start + idx + 1
		// Return if path token isn't 'empty' (/)
		if end-start > 1 {
			return path[start:end], end
		}
		// check
		start = end
	}
	return path[start:], -1
}
