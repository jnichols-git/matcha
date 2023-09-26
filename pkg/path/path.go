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

// Package path manages the tokenization of http paths.
package path

import (
	"strings"
)

// Return the next token from a path, starting at position last, and the position to use with the next call.
// Next considers multiple consecutive slashes to act as a single slash.
func Next(path string, last int) (string, int) {
	if path == "" && last == 0 {
		path = "/"
	}
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

// MakePartial gives the partial equivalent of a route.
// This effectively appends /+ to the path.
func MakePartial(path string, param string) string {
	if param != "" {
		param = "[" + param + "]"
	}
	i := len(path) - 1
	if path[i-1:] == "/+" {
		path = path[:i-1]
	} else if path[i] == '/' {
		path = path[:i]
	}
	return path + "/" + param + "+"
}
