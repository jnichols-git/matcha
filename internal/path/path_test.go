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
	path = "/trailing/slash/"
	expected = []string{"/trailing", "/slash", "/"}
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
	if px := MakePartial("/hello", "next"); px != "/hello/{next}+" {
		t.Error("/hello/{next}+", px)
	}
}
