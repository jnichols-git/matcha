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
	"testing"
)

func use(any) {}

func BenchmarkHostPattern(b *testing.B) {
	b.Run("regexp.Regexp", func(b *testing.B) {
		expr := regexp.MustCompile(`.+\.decentplatforms\..+`)
		for i := 0; i < b.N; i++ {
			use(expr.MatchString("www.decentplatforms.com"))
		}
	})
	b.Run("regex.Pattern-full", func(b *testing.B) {
		rs, isrs, err := CompilePattern(`{.+\.decentplatforms\..+}`)
		if err != nil || !isrs {
			b.Fatal(err)
		}
		for i := 0; i < b.N; i++ {
			use(rs.Match("www.decentplatforms.com"))
		}
	})
	b.Run("regex.Pattern-partial", func(b *testing.B) {
		rs, isrs, err := CompilePattern(`{.+}.decentplatforms.{.+}`)
		if err != nil || !isrs {
			b.Fatal(err)
		}
		for i := 0; i < b.N; i++ {
			use(rs.Match("www.decentplatforms.com"))
		}
	})
}
