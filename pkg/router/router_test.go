package router

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/cloudretic/router/pkg/rctx"
	"github.com/cloudretic/router/pkg/route"
)

// Return a handler that writes OK to all requests
func okHandler(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}
}

// Return a handler that writes http.StatusNotFound with body "not found"
func nfHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

// Return a handler that writes OK and a body containing the requested router param.
// If the param doesn't exist, writes internal error and "router param %s not found".
// This shouldn't happen unless something about router params is failing.
func rpHandler(rp string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := rctx.GetParam(r.Context(), rp)
		if p == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("router param %s not found", rp)))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(p))
		}
	}
}

func genericValueHandler(key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(key).(string)
		if p == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("context value %s not found", key)))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(p))
		}
	}
}

func errConf(rt Router) error {
	return errors.New("this should cause a router to fail")
}

func testMiddleware(w http.ResponseWriter, req *http.Request) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), "mwkey", "mwval"))
}

func reject(w http.ResponseWriter, req *http.Request) *http.Request {
	w.WriteHeader(http.StatusForbidden)
	return nil
}

func reqGen(method string) func(url, path string) *http.Request {
	return func(url, path string) *http.Request {
		req, _ := http.NewRequest(method, url+path, nil)
		return req
	}
}

func runEvalRequest(t *testing.T,
	s *httptest.Server, path string, genReqTo func(string, string) *http.Request, expect map[string]any) {
	t.Run(path, func(t *testing.T) {
		req := genReqTo(s.URL, path)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		for k, v := range expect {
			switch k {
			case "code":
				if res.StatusCode != v.(int) {
					t.Errorf("expected code %d, got %d", v.(int), res.StatusCode)
				}
			case "body":
				if res.Body == nil {
					t.Error("response has no body")
					break
				}
				body, err := io.ReadAll(res.Body)
				if err != nil {
					t.Error(err)
					break
				}
				if string(body) != v.(string) {
					t.Errorf("expected body %s, got %s", v.(string), string(body))
				}
			}
		}
	})
}

// Test all options of New().
// Doesn't check to see if those options work, just that they compile and don't cause errors. Check options individually.
func TestNewRouter(t *testing.T) {
	_, err := New(Default(),
		WithRoute(route.Declare(http.MethodGet, "/"), okHandler("root")),
		WithNotFound(nfHandler()),
		WithMiddleware(testMiddleware),
	)
	if err != nil {
		t.Error(err)
	}
	_, err = New(Default(), errConf)
	if err == nil {
		t.Error("router New should return an error if config returns an error")
	}
}

func TestBasicRoutes(t *testing.T) {
	r, err := New(Default(),
		WithRoute(route.Declare(http.MethodGet, "/"), okHandler("root")),
		WithRoute(route.Declare(http.MethodGet, "/[wildcard]"), rpHandler("wildcard")),
		WithRoute(route.Declare(http.MethodGet, `/route/{[a-zA-Z]+}`), okHandler("letters")),
		WithRoute(route.Declare(http.MethodGet, `/route/[id]{[\w]{4}}`), rpHandler("id")),
		WithRoute(route.Declare(http.MethodGet, `/static/file/[filename]{\w+(?:\.\w+)?}+`), rpHandler("filename")),
		WithMiddleware(testMiddleware),
		WithRoute(route.Declare(http.MethodGet, "/middlewareTest"), genericValueHandler("mwkey")),
	)
	if err != nil {
		t.Fatal(err)
	}
	s := httptest.NewServer(r)
	runEvalRequest(t, s, "/", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusOK,
		"body": "root",
	})
	runEvalRequest(t, s, "/test", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusOK,
		"body": "test",
	})
	runEvalRequest(t, s, "/route/word", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusOK,
		"body": "letters",
	})
	runEvalRequest(t, s, "/route/id01", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusOK,
		"body": "id01",
	})
	runEvalRequest(t, s, "/route/n0tID", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusNotFound,
	})
	runEvalRequest(t, s, "/static/file/docs/README.md", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusOK,
		"body": "/docs/README.md",
	})
	runEvalRequest(t, s, "/static/file", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusInternalServerError,
		"body": "router param filename not found",
	})
	runEvalRequest(t, s, "/middlewareTest", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusOK,
		"body": "mwval",
	})
}

func TestConcurrent(t *testing.T) {
	r := Declare(Default(),
		WithRoute(route.Declare(http.MethodGet, "/"), okHandler("root")),
		WithRoute(route.Declare(http.MethodGet, "/[wildcard]"), rpHandler("wildcard")),
		WithRoute(route.Declare(http.MethodGet, `/route/[id]{[a-zA-Z]+}`), rpHandler("id")),
		WithRoute(route.Declare(http.MethodGet, `/route/[id]{[\w]{4}}`), rpHandler("id")),
		WithRoute(route.Declare(http.MethodGet, `/static/file/[filename]{\w+(?:\.\w+)?}+`), rpHandler("filename")),
	)
	s := httptest.NewServer(r)
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		parseBody := func(res *http.Response) string {
			raw, _ := io.ReadAll(res.Body)
			return string(raw)
		}
		wg.Add(1)
		go func() {
			r1 := reqGen(http.MethodGet)(s.URL, "/")
			r2 := reqGen(http.MethodGet)(s.URL, "/12345")
			r3 := reqGen(http.MethodGet)(s.URL, "/longid")
			r4 := reqGen(http.MethodGet)(s.URL, "/id01")
			r5 := reqGen(http.MethodGet)(s.URL, "/static/file/some/file.txt")
			res1, err := http.DefaultClient.Do(r1)
			if err != nil {
				t.Error(err)
			}
			res2, _ := http.DefaultClient.Do(r2)
			res3, _ := http.DefaultClient.Do(r3)
			res4, _ := http.DefaultClient.Do(r4)
			res5, _ := http.DefaultClient.Do(r5)
			if body := parseBody(res1); body != "root" {
				t.Error("root", body)
			}
			if body := parseBody(res2); body != "12345" {
				t.Error("12345", body)
			}
			if body := parseBody(res3); body != "longid" {
				t.Error("longid", body)
			}
			if body := parseBody(res4); body != "id01" {
				t.Error("id01", body)
			}
			if body := parseBody(res5); body != "/some/file.txt" {
				t.Error("/some/file.txt", body)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestDeclare(t *testing.T) {
	// Test rejection middleware
	r := Declare(Default(), WithMiddleware(reject))
	s := httptest.NewServer(r)
	runEvalRequest(t, s, "/", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusForbidden,
	})
	// Test declaration fail
	var err error
	defer func() {
		err = recover().(error)
	}()
	Declare(Default(), errConf)
	if err == nil {
		t.Error("expected declare to fail and panic")
	}
}
