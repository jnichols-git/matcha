package regex

import "testing"

func TestPattern(t *testing.T) {
	rs, isrs, err := CompilePattern("{(api|www)}.cloudretic.{.*}")
	if err != nil {
		t.Errorf("expected expression to compile, got %s", err)
	} else if !isrs {
		t.Errorf("expected expression to compile to pattern")
	}
	if ok := rs.Match("api.cloudretic.com"); !ok {
		t.Error("expected match")
	}
	if ok := rs.Match("blog.cloudretic.com"); ok {
		t.Error("expected no match")
	}
	if ok := rs.Match("api.google.com"); ok {
		t.Error("expected no match")
	}
	if ok := rs.Match("cloudretic.com"); ok {
		t.Error("expected no match")
	}
	rs, isrs, err = CompilePattern("{.{4}}{.+}")
	if err != nil {
		t.Error(err)
	} else if !isrs {
		t.Errorf("expected expression to compile to pattern")
	}
	if ok := rs.Match("abcde"); !ok {
		t.Error("expected match")
	}
	if ok := rs.Match("abcd"); ok {
		t.Error("expected no match")
	}
	rs, _, _ = CompilePattern("{.+}.cloudretic.com")
	if ok := rs.Match("cloudretic.com"); ok {
		t.Error("expected no match")
	}
	if ok := rs.Match("www.cloudretic.com:80"); ok {
		t.Error("expected no match")
	}

	rs, isrs, err = CompilePattern("api.cloudretic.com")
	if err != nil {
		t.Errorf("expected no error")
	}
	if isrs {
		t.Errorf("static string is not a Pattern")
	}
	_, _, err = CompilePattern("{.+}.cloudretic.{.+")
	if err == nil {
		t.Errorf("should fail with unbalanced braces")
	}
	_, _, err = CompilePattern("{[}{.*}")
	if err == nil {
		t.Errorf("should fail with invalid regex")
	}
	_, _, err = CompilePattern("{[}")
	if err == nil {
		t.Errorf("should fail with invalid regex")
	}
}
