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
