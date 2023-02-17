package route

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func invalidConfigFunc(r Route) error {
	return errors.New("invalid config")
}

func TestStringRouteNew(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		rt, err := New(http.MethodGet, "/test")
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/test", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		}
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/test2", nil)
		if req = rt.MatchAndUpdateContext(req); req != nil {
			t.Errorf("expected route to not match")
		}
	})
}
func TestStringRouteDeclare(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		rt := Declare(http.MethodGet, "/test")
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/test", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		}
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/test2", nil)
		if req = rt.MatchAndUpdateContext(req); req != nil {
			t.Errorf("expected route to not match")
		}
	})
}

func TestWildcardRouteNew(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		rt, err := New(http.MethodGet, "/[param1]/[param2]/[param3]")
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/test1/test2/test3", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		}
		if p1 := GetParam(req.Context(), "param1"); p1 == "" || p1 != "test1" {
			t.Errorf("expected route param param1=test1; got %s", p1)
		}
		if p2 := GetParam(req.Context(), "param2"); p2 == "" || p2 != "test2" {
			t.Errorf("expected route param param2=test2; got %s", p2)
		}
		if p3 := GetParam(req.Context(), "param3"); p3 == "" || p3 != "test3" {
			t.Errorf("expected route param param3=test3; got %s", p3)
		}
	})
}
func TestWildcardRouteDeclare(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		rt := Declare(http.MethodGet, "/[param1]/[param2]/[param3]")
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/test1/test2/test3", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		}
		if p1 := GetParam(req.Context(), "param1"); p1 == "" || p1 != "test1" {
			t.Errorf("expected route param param1=test1; got %s", p1)
		}
		if p2 := GetParam(req.Context(), "param2"); p2 == "" || p2 != "test2" {
			t.Errorf("expected route param param2=test2; got %s", p2)
		}
		if p3 := GetParam(req.Context(), "param3"); p3 == "" || p3 != "test3" {
			t.Errorf("expected route param param3=test3; got %s", p3)
		}
	})
}

func TestRegexRouteNew(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		rt, err := New(http.MethodGet, "/{[a-zA-Z]{4}}")
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/test", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		}
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/t3st", nil)
		if req = rt.MatchAndUpdateContext(req); req != nil {
			t.Errorf("expected route to not match")
		}
	})
	t.Run("invalid-regex", func(t *testing.T) {
		rt, err := New(http.MethodGet, "/{[a-zA-Z{4}}")
		if err == nil || rt != nil {
			t.Errorf("expected route to fail with invalid regex")
		}
	})
}
func TestRegexRouteDeclare(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		rt := Declare(http.MethodGet, "/{[a-zA-Z]{4}}")
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/test", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		}
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/t3st", nil)
		if req = rt.MatchAndUpdateContext(req); req != nil {
			t.Errorf("expected route to not match")
		}
	})
	t.Run("invalid-regex", func(t *testing.T) {
		var err error
		defer func() {
			err = recover().(error)
		}()
		rt := Declare(http.MethodGet, "/{[a-zA-Z{4}}")
		if err == nil || rt != nil {
			t.Errorf("expected route to panic with invalid regex")
		}
	})
}

func TestPartialRouteNew(t *testing.T) {
	t.Run("valid-basic", func(t *testing.T) {
		rt, err := New(http.MethodGet, `/partial/+`)
		if reflect.TypeOf(rt) != reflect.TypeOf(&partialRoute{}) {
			t.Fatalf("/partial/+ should create a partialRoute, got %s", reflect.TypeOf(rt).String())
		}
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/partial/any/path", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		}
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/partial", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("partial routes should match their roots")
		}
	})
	t.Run("valid-filename", func(t *testing.T) {
		rt, err := New(http.MethodGet, `/file/[filename]{\w+(?:\.\w+)?}+`)
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/file/README.md", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		} else {
			param := GetParam(req.Context(), "filename")
			if param == "" {
				t.Errorf("expected a filename param")
			} else {
				if param != "/README.md" {
					t.Errorf("expected filename param %s, got %s", "/README.md", param)
				}
			}
		}
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/file/complex/path/file.txt", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		} else {
			param := GetParam(req.Context(), "filename")
			if param == "" {
				t.Errorf("expected a filename param")
			} else {
				if param != "/complex/path/file.txt" {
					t.Errorf("expected filename param %s, got %s", "/complex/path/file.txt", param)
				}
			}
		}
	})
	t.Run("invalid-regex", func(t *testing.T) {
		rt, err := New(http.MethodGet, "/complex/{[regex}+")
		if err == nil || rt != nil {
			t.Errorf("expected route.New to fail for partial with invalid regex")
		}
	})
}

func TestPartialRouteDeclare(t *testing.T) {
	t.Run("valid-basic", func(t *testing.T) {
		rt := Declare(http.MethodGet, `/partial/+`)
		if reflect.TypeOf(rt) != reflect.TypeOf(&partialRoute{}) {
			t.Fatalf("/partial/+ should create a partialRoute, got %s", reflect.TypeOf(rt).String())
		}
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/partial/any/path", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		}
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/partial", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("partial routes should match their roots")
		}
	})
	t.Run("valid-filename", func(t *testing.T) {
		rt := Declare(http.MethodGet, `/file/[filename]{\w+(?:\.\w+)?}+`)
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/file/README.md", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		} else {
			param := GetParam(req.Context(), "filename")
			if param == "" {
				t.Errorf("expected a filename param")
			} else {
				if param != "/README.md" {
					t.Errorf("expected filename param %s, got %s", "/README.md", param)
				}
			}
		}
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/file/complex/path/file.txt", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		} else {
			param := GetParam(req.Context(), "filename")
			if param == "" {
				t.Errorf("Expected a filename param")
			} else {
				if param != "/complex/path/file.txt" {
					t.Errorf("expected filename param %s, got %s", "/complex/path/file.txt", param)
				}
			}
		}
	})
	t.Run("invalid-regex", func(t *testing.T) {
		var err error
		defer func() {
			err = recover().(error)
		}()
		rt := Declare(http.MethodGet, "/complex/{[regex}+")
		if err == nil || rt != nil {
			t.Errorf("expected declare to fail with invalid regex")
		}
	})
}

func TestInvalidConfig(t *testing.T) {
	rt, err := New(http.MethodGet, "/static/path", invalidConfigFunc)
	if err == nil || rt != nil {
		t.Errorf("expected New to fail if ConfigFunc returns error")
	}
}
