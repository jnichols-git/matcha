package route

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/cloudretic/matcha/pkg/cors"
	"github.com/cloudretic/matcha/pkg/rctx"
	"github.com/cloudretic/matcha/pkg/route/require"
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
		// hash
		if hash := rt.Hash(); hash != "GET /test" {
			t.Errorf("expected hash '/test', got %s", hash)
		}
		// length
		if length := rt.Length(); length != 1 || length != len(rt.Parts()) {
			t.Errorf("expected length 1, got %d", length)
		}
		// prefix
		if prefix := rt.Prefix(); prefix != "/test" {
			t.Errorf("expected prefix '/test', got '%s'", prefix)
		}
		// method
		if method := rt.Method(); method != http.MethodGet {
			t.Errorf("expected method '%s', got '%s'", http.MethodGet, method)
		}
		// valid request
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/test", nil)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		}
		// incorrect path
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/test2", nil)
		if req = rt.MatchAndUpdateContext(req); req != nil {
			t.Errorf("expected route to not match")
		}
		// incorrect path length
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/static/test", nil)
		if req = rt.MatchAndUpdateContext(req); req != nil {
			t.Errorf("expected route to not match")
		}
		// incorrect method
		req, _ = http.NewRequest(http.MethodPost, "http://url.com/test", nil)
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
		// hash
		if hash := rt.Hash(); hash != "GET /[param1]/[param2]/[param3]" {
			t.Errorf("expected hash '/[param1]/[param2]/[param3]', got %s", hash)
		}
		// length
		if length := rt.Length(); length != 3 || length != len(rt.Parts()) {
			t.Errorf("expected length 3, got %d", length)
		}
		// prefix
		if prefix := rt.Prefix(); prefix != "*" {
			t.Errorf("expected prefix '*', got '%s'", prefix)
		}
		// method
		if method := rt.Method(); method != http.MethodGet {
			t.Errorf("expected method '%s', got '%s'", http.MethodGet, method)
		}
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/test1/test2/test3", nil)
		req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		}
		if p1 := rctx.GetParam(req.Context(), "param1"); p1 == "" || p1 != "test1" {
			t.Errorf("expected route param param1=test1; got %s", p1)
		}
		if p2 := rctx.GetParam(req.Context(), "param2"); p2 == "" || p2 != "test2" {
			t.Errorf("expected route param param2=test2; got %s", p2)
		}
		if p3 := rctx.GetParam(req.Context(), "param3"); p3 == "" || p3 != "test3" {
			t.Errorf("expected route param param3=test3; got %s", p3)
		}
	})
}
func TestWildcardRouteDeclare(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		rt := Declare(http.MethodGet, "/[param1]/[param2]/[param3]")
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/test1/test2/test3", nil)
		req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		}
		if p1 := rctx.GetParam(req.Context(), "param1"); p1 == "" || p1 != "test1" {
			t.Errorf("expected route param param1=test1; got %s", p1)
		}
		if p2 := rctx.GetParam(req.Context(), "param2"); p2 == "" || p2 != "test2" {
			t.Errorf("expected route param param2=test2; got %s", p2)
		}
		if p3 := rctx.GetParam(req.Context(), "param3"); p3 == "" || p3 != "test3" {
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
		// hash
		if hash := rt.Hash(); hash != "GET /{[a-zA-Z]{4}}" {
			t.Errorf("expected hash '/test', got %s", hash)
		}
		// length
		if length := rt.Length(); length != 1 || length != len(rt.Parts()) {
			t.Errorf("expected length 1, got %d", length)
		}
		// prefix
		if prefix := rt.Prefix(); prefix != "*" {
			t.Errorf("expected prefix '*', got '%s'", prefix)
		}
		// method
		if method := rt.Method(); method != http.MethodGet {
			t.Errorf("expected method '%s', got '%s'", http.MethodGet, method)
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
		// hash
		if hash := rt.Hash(); hash != "GET /partial/+" {
			t.Errorf("expected hash '/partial/+', got %s", hash)
		}
		// length (partial routes do *not* include the extension in their length!)
		if length := rt.Length(); length != 1 || length != len(rt.Parts())-1 {
			t.Errorf("expected length 1, got %d", length)
		}
		// prefix
		if prefix := rt.Prefix(); prefix != "/partial" {
			t.Errorf("expected prefix '/partial', got '%s'", prefix)
		}
		// method
		if method := rt.Method(); method != http.MethodGet {
			t.Errorf("expected method '%s', got '%s'", http.MethodGet, method)
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
		req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		} else {
			param := rctx.GetParam(req.Context(), "filename")
			if param == "" {
				t.Errorf("expected a filename param")
			} else {
				if param != "/README.md" {
					t.Errorf("expected filename param %s, got %s", "/README.md", param)
				}
			}
		}
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/file/complex/path/file.txt", nil)
		req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		} else {
			param := rctx.GetParam(req.Context(), "filename")
			if param == "" {
				t.Errorf("expected a filename param")
			} else {
				if param != "/complex/path/file.txt" {
					t.Errorf("expected filename param %s, got %s", "/complex/path/file.txt", param)
				}
			}
		}
		// invalid method
		req, _ = http.NewRequest(http.MethodPut, "http://url.com/file/complex/path/file.txt", nil)
		if req = rt.MatchAndUpdateContext(req); req != nil {
			t.Error("request shouldn't match with incorrect method")
		}
		// invalid name (regex validation failed)
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/file/invalid/name.txt.bck", nil)
		if req = rt.MatchAndUpdateContext(req); req != nil {
			t.Errorf("Expected route to fail when partial part doesn't match")
		}
	})
	t.Run("valid-nested", func(t *testing.T) {
		rt, err := New(http.MethodGet, "/nested/partial/route/[proxy]+")
		if err != nil {
			t.Fatal(err)
		}
		req, _ := http.NewRequest(http.MethodGet, "http://url.com/nested", nil)
		if req = rt.MatchAndUpdateContext(req); req != nil {
			t.Errorf("Expected route to fail when too short")
		}
	})
	t.Run("valid-root", func(t *testing.T) {
		rt, err := New(http.MethodGet, "/[rt]+")
		if err != nil {
			t.Error(err)
		}
		if prefix := rt.Prefix(); prefix != "*" {
			t.Errorf("expected prefix '*', got '%s'", prefix)
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
		req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		} else {
			param := rctx.GetParam(req.Context(), "filename")
			if param == "" {
				t.Errorf("expected a filename param")
			} else {
				if param != "/README.md" {
					t.Errorf("expected filename param %s, got %s", "/README.md", param)
				}
			}
		}
		req, _ = http.NewRequest(http.MethodGet, "http://url.com/file/complex/path/file.txt", nil)
		req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
		if req = rt.MatchAndUpdateContext(req); req == nil {
			t.Errorf("expected route to match")
		} else {
			param := rctx.GetParam(req.Context(), "filename")
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

func TestRouteEdgeCases(t *testing.T) {
	rt, err := New(http.MethodGet, "/consec///slash///route")
	if err != nil {
		t.Fatal(err)
	}
	if h := rt.Hash(); h != "GET /consec/slash/route" {
		t.Errorf("expected hash 'GET /consec/slash/route', got %s", h)
	}
}

func TestInvalidConfig(t *testing.T) {
	rt, err := New(http.MethodGet, "/static/path", invalidConfigFunc)
	if err == nil || rt != nil {
		t.Errorf("expected New to fail if ConfigFunc returns error")
	}
}

func TestCORS(t *testing.T) {
	// Basic
	var aco = &cors.AccessControlOptions{
		AllowOrigin:      []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		MaxAge:           1000,
		AllowCredentials: false,
	}
	rt, err := New(http.MethodGet, "/static/path", CORSHeaders(aco))
	if err != nil || len(rt.Middleware()) != 1 {
		t.Fatal(err)
	}
	u, _ := url.Parse("http://test.com/static/path")
	req := &http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: http.Header{
			"Origin": {"origin.com"},
		},
	}
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	req = rt.MatchAndUpdateContext(req)
	headers := req.Header
	if headers.Get("Origin") != "origin.com" {
		t.Errorf("Expected origin origin.com, got %s", headers.Get("Origin"))
	}
	// Partial/Non-builtin conf
	rt, err = New(http.MethodGet, "/static/path/[add]+", WithMiddleware(cors.CORSMiddleware(aco)))
	if err != nil || len(rt.Middleware()) != 1 {
		t.Fatal(err)
	}
	u, _ = url.Parse("http://test.com/static/path/with/addition")
	req = &http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: http.Header{
			"Origin": {"origin.com"},
		},
	}
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	req = rt.MatchAndUpdateContext(req)
	headers = req.Header
	if headers.Get("Origin") != "origin.com" {
		t.Errorf("Expected origin origin.com, got %s", headers.Get("Origin"))
	}
	if rctx.GetParam(req.Context(), "add") != "/with/addition" {
		t.Errorf("Expected param /with/addition, got %s", rctx.GetParam(req.Context(), "add"))
	}
}

func TestRequire(t *testing.T) {
	webhost := require.Hosts("cloudretic.com", "www.cloudretic.com")
	apihost := require.Hosts("api.cloudretic.com")
	webr := Declare(http.MethodGet, "/", Require(webhost))
	apir := Declare(http.MethodGet, "/")
	apir.Require(apihost)

	req := httptest.NewRequest(http.MethodGet, "https://www.cloudretic.com", nil)
	if !require.Execute(req, webr.Required()) {
		t.Error("expected match")
	}
	if require.Execute(req, apir.Required()) {
		t.Error("expected no match")
	}

	req = httptest.NewRequest(http.MethodGet, "https://api.cloudretic.com", nil)
	if !require.Execute(req, apir.Required()) {
		t.Error("expected match")
	}
	if require.Execute(req, webr.Required()) {
		t.Error("expected no match")
	}

	// Repeat for partial routes
	webr = Declare(http.MethodGet, "/+", Require(webhost))
	apir = Declare(http.MethodGet, "/+")
	apir.Require(apihost)

	req = httptest.NewRequest(http.MethodGet, "https://www.cloudretic.com", nil)
	if !require.Execute(req, webr.Required()) {
		t.Error("expected match")
	}
	if require.Execute(req, apir.Required()) {
		t.Error("expected no match")
	}

	req = httptest.NewRequest(http.MethodGet, "https://api.cloudretic.com", nil)
	if !require.Execute(req, apir.Required()) {
		t.Error("expected match")
	}
	if require.Execute(req, webr.Required()) {
		t.Error("expected no match")
	}
}
