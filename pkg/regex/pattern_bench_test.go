package regex

import (
	"regexp"
	"testing"
)

func use(any) {}

func BenchmarkHostPattern(b *testing.B) {
	b.Run("regexp.Regexp", func(b *testing.B) {
		expr := regexp.MustCompile(`.+\.cloudretic\..+`)
		for i := 0; i < b.N; i++ {
			use(expr.MatchString("www.cloudretic.com"))
		}
	})
	b.Run("regex.Pattern-full", func(b *testing.B) {
		rs, isrs, err := CompilePattern(`{.+\.cloudretic\..+}`)
		if err != nil || !isrs {
			b.Fatal(err)
		}
		for i := 0; i < b.N; i++ {
			use(rs.Match("www.cloudretic.com"))
		}
	})
	b.Run("regex.Pattern-partial", func(b *testing.B) {
		rs, isrs, err := CompilePattern(`{.+}.cloudretic.{.+}`)
		if err != nil || !isrs {
			b.Fatal(err)
		}
		for i := 0; i < b.N; i++ {
			use(rs.Match("www.cloudretic.com"))
		}
	})
}
