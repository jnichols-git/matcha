package route

import (
	"net/http"
	"reflect"
	"testing"
)

func TestStringRoute(t *testing.T) {
	rt, err := New(http.MethodGet, "/test")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/test", nil)
	if req = rt.MatchAndUpdateContext(req); req == nil {
		t.Errorf("Expected route to match")
	}
	req, _ = http.NewRequest(http.MethodGet, "http://url.com/test2", nil)
	if req = rt.MatchAndUpdateContext(req); req != nil {
		t.Errorf("Expected route to not match")
	}
}

func TestWildcardRoute(t *testing.T) {
	rt, err := New(http.MethodGet, "/[param1]/[param2]/[param3]")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/test1/test2/test3", nil)
	if req = rt.MatchAndUpdateContext(req); req == nil {
		t.Errorf("Expected route to match")
	}
	if p1 := GetParam(req.Context(), "param1"); p1 == "" || p1 != "test1" {
		t.Errorf("Expected route param param1=test1; got %s", p1)
	}
	if p2 := GetParam(req.Context(), "param2"); p2 == "" || p2 != "test2" {
		t.Errorf("Expected route param param2=test2; got %s", p2)
	}
	if p3 := GetParam(req.Context(), "param3"); p3 == "" || p3 != "test3" {
		t.Errorf("Expected route param param3=test3; got %s", p3)
	}
}

func TestRegexRoute(t *testing.T) {
	rt, err := New(http.MethodGet, "/{[a-zA-Z]{4}}")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/test", nil)
	if req = rt.MatchAndUpdateContext(req); req == nil {
		t.Errorf("Expected route to match")
	}
	req, _ = http.NewRequest(http.MethodGet, "http://url.com/t3st", nil)
	if req = rt.MatchAndUpdateContext(req); req != nil {
		t.Errorf("Expected route to not match")
	}
}

func TestPartialRoute(t *testing.T) {
	t.Run("basics", func(t *testing.T) {
		rt, err := New(http.MethodGet, `/partial/+`)
		if reflect.TypeOf(rt) != reflect.TypeOf(&partialRoute{}) {
			t.Fatalf("/partial/+ should create a partialRoute, got %s", reflect.TypeOf(rt).String())
		}
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/partial/any/path", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("Expected route to match")
		}
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/partial", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("Partial routes should match their roots")
		}
	})
	t.Run("filename", func(t *testing.T) {
		rt, err := New(http.MethodGet, `/file/[filename]{\w+(?:\.\w+)?}+`)
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/file/README.md", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("Expected route to match")
		} else {
			param := GetParam(req.Context(), "filename")
			if param == "" {
				t.Errorf("Expected a filename param")
			} else {
				if param != "/README.md" {
					t.Errorf("Expected filename param %s, got %s", "/README.md", param)
				}
			}
		}
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/file/complex/path/file.txt", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("Expected route to match")
		} else {
			param := GetParam(req.Context(), "filename")
			if param == "" {
				t.Errorf("Expected a filename param")
			} else {
				if param != "/complex/path/file.txt" {
					t.Errorf("Expected filename param %s, got %s", "/complex/path/file.txt", param)
				}
			}
		}
	})
}

func TestMethodRoute(t *testing.T) {
	rt, err := New(http.MethodPost, "/test")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodPost, "http://url.com/test", nil)
	if req = rt.MatchAndUpdateContext(req); req == nil {
		t.Errorf("Expected POST to match")
	}
	req, _ = http.NewRequest(http.MethodGet, "http://url.com/test", nil)
	if req = rt.MatchAndUpdateContext(req); req != nil {
		t.Errorf("Expected GET not to match")
	}
}
