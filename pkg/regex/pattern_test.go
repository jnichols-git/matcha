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
