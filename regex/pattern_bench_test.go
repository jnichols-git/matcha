package regex

import (
	"regexp"
	"testing"
)

func use(any) {}

func BenchmarkHostPattern(b *testing.B) {
	b.Run("regexp.Regexp", func(b *testing.B) {
		expr := regexp.MustCompile(`.+\.jnichols\..+`)
		for i := 0; i < b.N; i++ {
			use(expr.MatchString("www.jnichols.info"))
		}
	})
	b.Run("regex.Pattern-full", func(b *testing.B) {
		rs, err := CompilePattern(`[.+\.jnichols\..+]`)
		if err != nil {
			b.Fatal(err)
		}
		for i := 0; i < b.N; i++ {
			use(rs.Match("www.jnichols.info"))
		}
	})
	b.Run("regex.Pattern-partial", func(b *testing.B) {
		rs, err := CompilePattern(`[.+].jnichols.[.+]`)
		if err != nil {
			b.Fatal(err)
		}
		for i := 0; i < b.N; i++ {
			use(rs.Match("www.jnichols.info"))
		}
	})
	b.Run("regex.Pattern-string", func(b *testing.B) {
		rs, err := CompilePattern(`www.jnichols.info`)
		if err != nil {
			b.Fatal(err)
		}
		for i := 0; i < b.N; i++ {
			use(rs.Match("www.jnichols.info"))
		}
	})
}