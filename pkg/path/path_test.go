// Copyright 2023 Decent Platforms
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

package path

import "testing"

func TestNext(t *testing.T) {
	path := ""
	tk, next := Next(path, 0)
	if tk != "/" || next != -1 {
		t.Errorf("Empty path should return '/', -1, got '%s', %d", tk, next)
	}
	path = "/"
	tk, next = Next(path, 0)
	if tk != "/" || next != -1 {
		t.Errorf("Root path should return '/', -1, got '%s', %d", tk, next)
	}
	tk, next = Next(path, 10)
	if tk != "" || next != -1 {
		t.Errorf("Root path should return '', -1, got '%s', %d", tk, next)
	}
	path = "/path/to/file.txt"
	expected := []string{"/path", "/to", "/file.txt"}
	i := 0
	for next = 0; next != -1; {
		tk, next = Next(path, next)
		if tk != expected[i] {
			t.Errorf("Expected '%s' at %d, got '%s'", expected[i], i, tk)
		}
		i++
	}
	path = "/consec///slash"
	expected = []string{"/consec", "/slash"}
	i = 0
	for next = 0; next != -1; {
		tk, next = Next(path, next)
		if tk != expected[i] {
			t.Errorf("Expected '%s' at %d, got '%s'", expected[i], i, tk)
		}
		i++
	}
}

func BenchmarkNext(b *testing.B) {
	path := "/path/to/file.txt"
	next := 0
	for i := 0; i < b.N; i++ {
		_, next = Next(path, next)
		if next == -1 {
			next = 0
		}
	}
}

func TestMakePartial(t *testing.T) {
	if px := MakePartial("/hello", ""); px != "/hello/+" {
		t.Error("/hello/+", px)
	}
	if px := MakePartial("/hello/", ""); px != "/hello/+" {
		t.Error("/hello/+", px)
	}
	if px := MakePartial("/hello/+", ""); px != "/hello/+" {
		t.Error("/hello/+", px)
	}
	if px := MakePartial("/hello", "next"); px != "/hello/[next]+" {
		t.Error("/hello/[next]+", px)
	}
}
