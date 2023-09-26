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

package regex

import (
	"regexp"
	"strings"
	"testing"
)

func TestGroups(t *testing.T) {
	r := regexp.MustCompile(`^([a-zA-Z]\w+)\.([a-zA-Z]+)`)
	fns := []string{
		"file.txt",
		"file.go",
		"file2.txt",
	}
	for _, fn := range fns {
		expected := strings.Split(fn, ".")
		groups := Groups(r, fn)
		for i := range expected {
			if expected[i] != groups[i] {
				t.Errorf("Expected '%s' at %d, got '%s'", expected[i], i, groups[i])
			}
		}
	}
	nomatch := []string{
		"2file.txt",
	}
	for _, fn := range nomatch {
		groups := Groups(r, fn)
		if groups != nil {
			t.Errorf("Should return nil for non-matching values, got %v", groups)
		}
	}
}
