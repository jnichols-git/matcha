// package regex wraps the Golang regexp package
package regex

import (
	"regexp"
)

// Convenience function for getting groups from a regex match
// Returns nil if no match is made, otherwise returns the groups that were matched
// via regexp.FindAllStringsSubmatch
func Groups(r *regexp.Regexp, token string) []string {
	matches := r.FindAllStringSubmatch(token, -1)
	if matches == nil {
		return nil
	} else {
		return matches[0][1:]
	}
}
