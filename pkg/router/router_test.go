package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudretic/router/pkg/route"
	"github.com/cloudretic/router/pkg/router/params"
)

// Return a handler that writes OK to all requests
func okHandler(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}
}

// Return a handler that writes 404 with body "not found"
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
		p, ok := params.Get(r, rp)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("router param %s not found", rp)))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(p))
		}
	}
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
	_, err := New(
		WithRoute(route.ForceNew("/"), okHandler("root")),
		WithNotFound(nfHandler()),
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBasicRoutes(t *testing.T) {
	r, err := New(
		WithRoute(route.ForceNew("/"), okHandler("root")),
		WithRoute(route.ForceNew("/[wildcard]"), rpHandler("wildcard")),
		WithRoute(route.ForceNew(`/route/{[a-zA-Z]+}`), okHandler("letters")),
		WithRoute(route.ForceNew(`/route/[id]{[\w]{4}}`), rpHandler("id")),
	)
	if err != nil {
		t.Fatal(err)
	}
	s := httptest.NewServer(r)
	runEvalRequest(t, s, "/", reqGen(http.MethodGet), map[string]any{
		"code": 200,
		"body": "root",
	})
	runEvalRequest(t, s, "/test", reqGen(http.MethodGet), map[string]any{
		"code": 200,
		"body": "test",
	})
	runEvalRequest(t, s, "/route/word", reqGen(http.MethodGet), map[string]any{
		"code": 200,
		"body": "letters",
	})
	runEvalRequest(t, s, "/route/id01", reqGen(http.MethodGet), map[string]any{
		"code": 200,
		"body": "id01",
	})
	runEvalRequest(t, s, "/route/n0tID", reqGen(http.MethodGet), map[string]any{
		"code": 404,
	})
}
