package router

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sync"

	"github.com/jnichols-git/matcha/v2/pkg/cors"
	"github.com/jnichols-git/matcha/v2/pkg/require"

	"github.com/jnichols-git/matcha/v2/internal/rctx"
	"github.com/jnichols-git/matcha/v2/internal/route"
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

func reqGenHeaders(method string, headers http.Header) func(url, path string) *http.Request {
	return func(url, path string) *http.Request {
		req, _ := http.NewRequest(method, url+path, nil)
		req.Header = headers
		return req
	}
}

func runEvalRequest(t *testing.T,
	s *httptest.Server, path string, genReqTo func(string, string) *http.Request, expect map[string]any) {
	var name string
	if len(path) > 0 {
		name = strings.ReplaceAll(path, "/", "-")[1:]
	} else {
		name = "empty"
	}
	t.Run(name, func(t *testing.T) {
		req := genReqTo(s.URL, path)
		req.Close = true
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
			case "header":
				hd := v.(http.Header)
				for h, vsAgainst := range hd {
					vsCheck := res.Header.Values(h)
					if len(vsCheck) != len(vsAgainst) {
						t.Errorf("expected headers len %d, got %d", len(vsAgainst), len(vsCheck))
					} else {
						for i := range vsCheck {
							if vsCheck[i] != vsAgainst[i] {
								t.Errorf("expected header value '%s' at %d, got '%s'", vsAgainst[i], i, vsAgainst[i])
							}
						}
					}
				}
			}
		}
	})
}

// Test routes; method, route, positive test, negative test, expected response
var testRoutes = [][5]string{
	{http.MethodGet, `/`, `/`, `/a`, "Got root"},
	{http.MethodGet, `/hello/world`, "/hello/world", "/hello", "Hello, World!"},
	{http.MethodGet, `/collection/`, "/collection/", "/collection", "Got root of /collection/"},
	{http.MethodGet, `/collection/:item`, "/collection/testItem", "/collection/testItem/more", "Got collection item testItem"},
	{http.MethodGet, `/users/:id[\d{4}]`, "/users/1234", "/users/abcd", "Got user 1234"},
	{http.MethodGet, `/file.txt`, "/file.txt", "/file", "Got file.txt"},
	{http.MethodGet, `/files/`, "/files/", "/files", "Got root of /files/"},
	{http.MethodGet, `/files/:filename+`, "/files/test/file.txt", "/files/test/other.txt", "Got file /test/file.txt"},
}

var testHandlers = []http.HandlerFunc{
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Got root"))
	},
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	},
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Got root of /collection/"))
	},
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Got collection item " + rctx.GetParam(r.Context(), "item")))
	},
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Got user " + rctx.GetParam(r.Context(), "id")))
	},
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Got file.txt"))
	},
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Got root of /files/"))
	},
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Got file " + rctx.GetParam(r.Context(), "filename")))
	},
}

func TestBasicRoutes(t *testing.T) {
	rt := Default()
	for i, tr := range testRoutes {
		rt.HandleFunc(tr[0], tr[1], testHandlers[i])
	}

	for _, tr := range testRoutes {
		t.Run(tr[1], func(t *testing.T) {
			// Test positive case
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tr[0], tr[2], nil)
			rt.ServeHTTP(w, req)
			res := w.Body.String()
			if w.Result().StatusCode != http.StatusOK || res != tr[4] {
				t.Error(tr[4], res)
			}
			// Test negative case
			w = httptest.NewRecorder()
			req = httptest.NewRequest(tr[0], tr[3], nil)
			rt.ServeHTTP(w, req)
			res = w.Body.String()
			if w.Result().StatusCode == http.StatusOK && res == tr[4] {
				t.Error(tr[3], "should have failed")
			}
		})
	}
}

func TestEdgeCaseRoutes(t *testing.T) {
	r := Default()
	r.HandleRoute(route.Declare(http.MethodGet, "/odd///path"), okHandler("odd"))
	r.HandleRoute(route.Declare(http.MethodGet, "/reject").Use(reject), okHandler("never"))
	r.HandleFunc(http.MethodGet, "/not/implemented/handler", nil)
	r.Handle(http.MethodGet, "/not/implemented/func", nil)
	r.HandleRoute(route.Declare(http.MethodGet, "/not/implemented/routehandler"), nil)
	r.HandleRouteFunc(route.Declare(http.MethodGet, "/not/implemented/routefunc"), nil)
	s := httptest.NewServer(r)
	runEvalRequest(t, s, "/odd/path", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusOK,
		"body": "odd",
	})
	runEvalRequest(t, s, "/odd///path", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusOK,
		"body": "odd",
	})
	runEvalRequest(t, s, "/reject", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusForbidden,
	})
	runEvalRequest(t, s, "/not/implemented/handler", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusNotImplemented,
	})
	runEvalRequest(t, s, "/not/implemented/func", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusNotImplemented,
	})
	runEvalRequest(t, s, "/not/implemented/routehandler", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusNotImplemented,
	})
	runEvalRequest(t, s, "/not/implemented/routefunc", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusNotImplemented,
	})

}

func TestConcurrent(t *testing.T) {
	r := Default()
	r.Handle(http.MethodGet, "/", okHandler("root"))
	r.HandleFunc(http.MethodGet, "/middlewareTest", genericValueHandler("mwkey"))
	r.HandleRoute(route.Declare(http.MethodGet, "/:wildcard"), rpHandler("wildcard"))
	r.HandleRouteFunc(route.Declare(http.MethodGet, `/route/[[a-zA-Z]+]`), okHandler("letters"))
	r.Handle(http.MethodGet, `/route/:id[[\w]{4}]`, rpHandler("id"))
	r.HandleFunc(http.MethodGet, `/static/file/:filename[\w+(?:\.\w+)?]+`, rpHandler("filename"))
	r.Use(testMiddleware)
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
	r := Default()
	r.Use(reject)
	s := httptest.NewServer(r)
	runEvalRequest(t, s, "/", reqGen(http.MethodGet), map[string]any{
		"code": http.StatusForbidden,
	})
}

var aco = &cors.AccessControlOptions{
	AllowOrigin:      []string{"*"},
	AllowMethods:     []string{"*"},
	AllowHeaders:     []string{"*"},
	ExposeHeaders:    []string{"*"},
	MaxAge:           1000,
	AllowCredentials: false,
}

func TestCORS(t *testing.T) {
	r := Default()
	r.Use(cors.CORSMiddleware(aco))
	r.HandleFunc(http.MethodOptions, "/", func(w http.ResponseWriter, r *http.Request) {
		cors.SetCORSResponseHeaders(w, r, aco)
		w.WriteHeader(http.StatusNoContent)
	})
	r.Handle(http.MethodGet, "/", okHandler("ok"))
	r.AddNotFound(nfHandler())
	s := httptest.NewServer(r)

	runEvalRequest(t, s, "/", reqGenHeaders(http.MethodGet, http.Header{"Origin": {"test-origin"}}), map[string]any{
		"code":   http.StatusOK,
		"body":   "ok",
		"header": http.Header{},
	})
	runEvalRequest(t, s, "/", reqGenHeaders(
		http.MethodOptions, http.Header{"Origin": {"test-origin"}, "Access-Control-Request-Headers": {"x-header-1"}},
	), map[string]any{
		"code":   http.StatusNoContent,
		"header": http.Header{"Access-Control-Allow-Headers": {"x-header-1"}},
	})
}

func TestDuplicate(t *testing.T) {
	h1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	r := Default()
	r.HandleRoute(route.Declare(http.MethodGet, "/duplicate/route"), h1)
	r.HandleRoute(route.Declare(http.MethodGet, "/duplicate/route"), h2)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/duplicate/route", nil)
	r.ServeHTTP(w, req)
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("re-declaration should be noop; expected %d, got %d", http.StatusOK, w.Result().StatusCode)
	}
}

func TestValidatedDuplicate(t *testing.T) {
	h1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	rt := Default()
	r_origin := route.Declare(http.MethodGet, "/")
	r_origin.Require(require.Hosts("origin.com"))
	rt.HandleRoute(r_origin, h1)
	rt.HandleRoute(route.Declare(http.MethodGet, "/"), h2)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://origin.com/", nil)
	rt.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("origin.com request should be OK; got %d", w.Code)
	}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rt.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("example.com request should fall through to BadRequest; got %d", w.Code)
	}
}

func TestComposition(t *testing.T) {
	// Set up a an API for composition
	h1 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	api1 := Default()
	api1.HandleFunc(http.MethodGet, "/hello", h1)
	api1.HandleFunc(http.MethodPost, "/hello", h1)
	api1.HandleFunc(http.MethodGet, "/hello/:name", h1)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	api1.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error(200, w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/hello", nil)
	api1.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error(200, w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/hello/jakenichols2719", nil)
	api1.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error(200, w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/hello/jakenichols2719.github.com", nil)
	api1.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error(200, w.Code)
	}

	// Pass through ONLY get
	api2 := Default()
	api2.Mount("/api", api1, http.MethodGet)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	api2.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error(200, w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/hello", nil)
	api2.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Error(404, w.Code)
	}

	// Pass through all methods
	api2 = Default()
	api2.Mount("/api", api1)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	api2.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error(200, w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/hello", nil)
	api2.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error(200, w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/hello/jakenichols2719.github.com", nil)
	api2.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error(200, w.Code)
	}

	// Invalid path test
	api2 = Default()
	err := api2.Mount("/{", api1)
	if err == nil {
		t.Error("expected failure due to route formatting")
	}
}
