package regex

import "testing"

func TestPattern(t *testing.T) {
	rs, err := CompilePattern("[(api|www)].jnichols.[.*]")
	if err != nil {
		t.Errorf("expected expression to compile, got %s", err)
	}
	if ok := rs.Match("api.jnichols.info"); !ok {
		t.Error("expected match")
	}
	if ok := rs.Match("blog.jnichols.info"); ok {
		t.Error("expected no match")
	}
	if ok := rs.Match("api.google.com"); ok {
		t.Error("expected no match")
	}
	if ok := rs.Match("jnichols.info"); ok {
		t.Error("expected no match")
	}
	rs, err = CompilePattern("[.{4}.+]")
	if err != nil {
		t.Error(err)
	}
	if ok := rs.Match("abcde"); !ok {
		t.Error("expected match")
	}
	if ok := rs.Match("abcd"); ok {
		t.Error("expected no match")
	}
	rs, _ = CompilePattern("[.+].jnichols.info")
	if ok := rs.Match("jnichols.info"); ok {
		t.Error("expected no match")
	}
	if ok := rs.Match("www.jnichols.info:80"); ok {
		t.Error("expected no match")
	}

	rs, err = CompilePattern("api.jnichols.info")
	if err != nil {
		t.Errorf("expected no error")
	}
	_, err = CompilePattern("[.+].jnichols.[.+")
	if err == nil {
		t.Errorf("should fail with unbalanced braces")
	}
	_, err = CompilePattern("[[}{.*]")
	if err == nil {
		t.Errorf("should fail with invalid regex")
	}
	_, err = CompilePattern("[[]")
	if err == nil {
		t.Errorf("should fail with invalid regex")
	}
}
