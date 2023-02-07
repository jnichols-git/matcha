package route

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/cloudretic/router/pkg/router/params"
)

func TestStringRoute(t *testing.T) {
	rt, err := New("/test")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/test", nil)
	if req = rt.MatchAndUpdateContext(req); req == nil {
		t.Errorf("Expected route to match")
	}
}

func TestWildcardRoute(t *testing.T) {
	rt, err := New("/[any]")
	if err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodGet, "http://url.com/test", nil)
	if req = rt.MatchAndUpdateContext(req); req == nil {
		t.Errorf("Expected route to match")
	}
}

func TestRegexRoute(t *testing.T) {
	rt, err := New("/{[a-zA-Z]{4}}")
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
		rt, err := New(`/partial/+`)
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
		rt, err := New(`/file/[filename]{\w+(?:\.\w+)?}+`)
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/file/README.md", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("Expected route to match")
		} else {
			param, ok := params.Get(req, "filename")
			if !ok {
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
			param, ok := params.Get(req, "filename")
			if !ok {
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
	rt, err := New("/test", WithMethods(http.MethodPost))
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
