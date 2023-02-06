// package regex wraps the Golang regexp package
package regex

import (
	"regexp"
)

func Groups(r *regexp.Regexp, token string) []string {
	matches := r.FindAllStringSubmatch(token, -1)
	if matches == nil {
		return nil
	} else {
		return matches[0][1:]
	}
}
