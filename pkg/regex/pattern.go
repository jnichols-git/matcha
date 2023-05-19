package regex

import (
	"errors"
	"regexp"
	"strings"
)

// findAllRegexGroups finds all regex expressions contained in {}s.
// These groups are returned in [][]int, an array of int paipatt representing the start
// and end of the regex groups. Returns an error if unbalanced brackets are found.
func findAllRegexGroups(expr string) ([][]int, error) {
	out := make([][]int, 0)
	start, end := -1, -1
	depth := 0
	for i, r := range expr {
		if r == '{' {
			if depth == 0 {
				start = i
			}
			depth++
		}
		if r == '}' {
			depth--
			if depth == 0 {
				end = i
			}
		}
		if start != -1 && end != -1 {
			out = append(out, []int{start, end + 1})
			start, end = -1, -1
		}
	}
	if depth != 0 {
		return nil, errors.New("unbalanced brackets in pattern " + expr)
	}
	return out, nil
}

// pmf is a pattern match function.
// pmfs take in a string and start position and output the end position of the token
// matched, or -1 if no match is found. this allows for easy token length tracking,
// since regex matches will frequently be of variable length.
type pmf func(string, int) int

// Patterns are a combination of static string components and regex validation.
// Any piece of the string contained in brackets {} will be matched as regex, while any outside of brackets
// will be matched as itself; for example, {.*}.cloudretic.{.*} matches any subdomain and top-level domain for
// cloudretic. Patterns *must* contain some regex. It performs generally equivalent to the exact same regex in
// brackets (see regex_bench_test.go) and better if static string parts are swapped out of the expression entirely.
type Pattern struct {
	fs      []pmf
	statics []string
}

// matchf generates a pmf for a given string token.
// If given a token in {}s, it will generate a pmf that matches against the contained regex; otherwise, it
// will match against the provided raw string.
//
// Returns pmf, isStatic, err, where isStatic represents the choice to use a static string or not; Patterns
// need this during compilation to store static parts, used to delimit inputs during matching.
func matchf(tk string) (pmf, bool, error) {
	if tk[0] == '{' && tk[len(tk)-1] == '}' {
		expr, err := regexp.Compile(tk[1 : len(tk)-1])
		if err != nil {
			return nil, false, err
		}
		return func(s string, i int) int {
			match := expr.FindString(s[i:])
			if match == "" {
				return -1
			}
			return i + strings.Index(s, match) + len(match)
		}, false, nil
	} else {
		return func(s string, i int) int {
			if end := i + len(tk); len(s) < end || tk != s[i:end] {
				return -1
			} else {
				return end
			}
		}, true, nil
	}
}

// resolve adds the pmf for the pattern within the set to the Pattern.
// This returns an error if the pmf fails to compile.
func resolve(patt *Pattern, expr string, set []int) error {
	f, static, err := matchf(expr[set[0]:set[1]])
	if err != nil {
		return err
	} else if static {
		patt.statics = append(patt.statics, expr[set[0]:set[1]])
	}
	patt.fs = append(patt.fs, f)
	return nil
}

// CompilePattern compiles a string pattern expression to a *regex.Pattern.
//
// Returns patt, isPatt, err, where err != nil if the pattern is invalid and !isPatt
// if the provided expression is valid, but not a pattern (a static string).
func CompilePattern(expr string) (*Pattern, bool, error) {
	// Static strings begone
	if !strings.ContainsAny(expr, "{}") {
		return nil, false, nil
	}
	patt := &Pattern{
		fs:      make([]pmf, 0),
		statics: make([]string, 0),
	}
	regexIndices, err := findAllRegexGroups(expr)
	if err != nil || len(regexIndices) == 0 {
		return nil, false, err
	}
	var set, next []int
	// Start with the first set
	set = regexIndices[0]
	for setidx := 0; setidx < len(regexIndices)-1; {
		// For each set, resolve the token in the set
		err = resolve(patt, expr, set)
		if err != nil {
			return nil, false, err
		}
		// Move to the next set OR whatever's in between this set and the next.
		next = regexIndices[setidx+1]
		if set[1] == next[0] {
			setidx++
		} else {
			next = []int{set[1], next[0]}
		}
		set = next
	}
	// Resolve the last set.
	err = resolve(patt, expr, set)
	if err != nil {
		return nil, false, err
	}
	return patt, true, nil
}

// findStatic finds the start index of the idx'th static element of the pattern.
func (patt *Pattern) findStatic(in string, idx int) int {
	if idx >= len(patt.statics) {
		return -1
	}
	return strings.Index(in, patt.statics[idx])
}

// Match matches a pattern to a string.
func (patt *Pattern) Match(str string) bool {
	static := 0
	i := 0
	for _, f := range patt.fs {
		tk := str
		// Limit tokens to static parts, if they're available. Avoids overconsumption from regex.
		if max := patt.findStatic(str, static); max != -1 {
			tk = tk[:max]
			static++
		}
		i = f(tk, i)
		if i == -1 {
			return false
		}
	}
	return true
}
