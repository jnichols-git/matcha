package regex

import "testing"

func TestPattern(t *testing.T) {
	rs, isrs, err := CompilePattern("{(api|www)}.decentplatforms.{.*}")
	if err != nil {
		t.Errorf("expected expression to compile, got %s", err)
	} else if !isrs {
		t.Errorf("expected expression to compile to pattern")
	}
	if ok := rs.Match("api.decentplatforms.com"); !ok {
		t.Error("expected match")
	}
	if ok := rs.Match("blog.decentplatforms.com"); ok {
		t.Error("expected no match")
	}
	if ok := rs.Match("api.google.com"); ok {
		t.Error("expected no match")
	}
	if ok := rs.Match("decentplatforms.com"); ok {
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
	rs, _, _ = CompilePattern("{.+}.decentplatforms.com")
	if ok := rs.Match("decentplatforms.com"); ok {
		t.Error("expected no match")
	}
	if ok := rs.Match("www.decentplatforms.com:80"); ok {
		t.Error("expected no match")
	}

	rs, isrs, err = CompilePattern("api.decentplatforms.com")
	if err != nil {
		t.Errorf("expected no error")
	}
	if isrs {
		t.Errorf("static string is not a Pattern")
	}
	_, _, err = CompilePattern("{.+}.decentplatforms.{.+")
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
