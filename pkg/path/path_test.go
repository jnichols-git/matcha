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
