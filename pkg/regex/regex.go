// Copyright 2023 Matcha Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package regex wraps the Golang regexp package to provide more convenient access to some features.
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
