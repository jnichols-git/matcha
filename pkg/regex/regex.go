// Package regex wraps the Golang regexp package to provide more convenient access to some features.
package regex

import (
	"regexp"

	"slices"
)

// Convenience function for getting groups from a regex match
// Returns nil if no match is made, otherwise returns the groups that were matched
// via regexp.FindAllStringsSubmatch
func Groups(r *regexp.Regexp, token string) []string {
	matches := r.FindAllStringSubmatch(token, -1)
	if matches == nil {
		return nil
	} else {
		return slices.DeleteFunc(matches[0][1:], func(s string) bool {
			return s == ""
		})
	}
}
