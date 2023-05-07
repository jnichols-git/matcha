package regex

import "testing"

func TestPattern(t *testing.T) {
	rs, isrs, err := CompilePattern("{.{3}}.cloudretic.{.*}")
	if err != nil {
		t.Errorf("expected pattern to compile, got %s", err)
	} else if !isrs {
		t.Errorf("expected pattern to compile to rich string")
	}
	if ok := rs.Match("api.cloudretic.com"); !ok {
		t.Error("expected match")
	}
	if ok := rs.Match("blog.cloudretic.com"); ok {
		t.Error("expected no match")
	}
	rs, isrs, err = CompilePattern("api.cloudretic.com")
	if err != nil {
		t.Errorf("expected no error")
	}
	if isrs {
		t.Errorf("static string is not a Pattern")
	}
}
