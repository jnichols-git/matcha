package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CloudRETIC/router/pkg/router/params"
)

type routeTestCase struct {
	name       string
	routeExpr  string
	req        func(url string) *http.Request
	handle     http.Handler
	expectBody string
}

func setupRouter(rt Router, tcs []routeTestCase) {
	for _, tc := range tcs {
		Handle(rt, tc.routeExpr, tc.handle)
	}
}

func body(a *http.Response) string {
	abody, _ := io.ReadAll(a.Body)
	return string(abody)
}

// Run a test case.
// Runs tc.req() through s, then compares its result against ts.handle directly.
// Returns the HTTP response time from the test.
func runTestCase(t *testing.T, s *httptest.Server, tc routeTestCase) time.Duration {
	ts := time.Now()
	resp, err := http.DefaultClient.Do(tc.req(s.URL))
	te := time.Now()
	if err != nil {
		t.Fatal(err)
	}
	rb := body(resp)
	if rb != tc.expectBody {
		t.Errorf("Expected response %s, got %s", tc.expectBody, rb)
	}
	return te.Sub(ts)
}

func runTestBench(b *testing.B, s *httptest.Server, tc routeTestCase) {
	resp, err := http.DefaultClient.Do(tc.req(s.URL))
	if err != nil {
		b.Fatal(err)
	}
	rb := body(resp)
	if rb != tc.expectBody {
		b.Errorf("Expected response %s, got %s", tc.expectBody, rb)
	}
}

func req(method, path string) func(url string) *http.Request {
	return func(url string) *http.Request {
		req, _ := http.NewRequest(method, url+path, nil)
		return req
	}
}

func testHandlerOKEmptyBody() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func testHandlerOKWildcardBody(wc string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, _ := params.Get(r, wc)
		w.Write([]byte(p))
	})
}

var basicRoutes []routeTestCase = []routeTestCase{
	{"root", "/", req(http.MethodGet, "/"), testHandlerOKEmptyBody(), ""},
	{"string literal", "/test", req(http.MethodGet, "/test"), testHandlerOKEmptyBody(), ""},
	{"empty wildcard", "/test/[wildcard]", req(http.MethodGet, "/test/testEmptyWildcard"), testHandlerOKWildcardBody("wildcard"), "testEmptyWildcard"},
	{"regex wildcard", `/rg/[rgexpr]{\w+}`, req(http.MethodGet, "/test/testWord"), testHandlerOKWildcardBody("rgexpr"), "testWord"},
}

var complexRoutes []routeTestCase = []routeTestCase{
	{
		"root",
		"/",
		req(http.MethodGet, "/"), testHandlerOKEmptyBody(),
		"",
	},
	{
		"device/[uuid]{...}/data",
		`/device/[uuid]{^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$}/data`,
		req(http.MethodGet, "/device/123e4567-e89b-12d3-a456-426655440000/data"), testHandlerOKWildcardBody("uuid"),
		"123e4567-e89b-12d3-a456-426655440000",
	},
	{
		"device/[uuid]{...}/register",
		`/device/[uuid]{^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$}/register`,
		req(http.MethodGet, "/device/123e4567-e89b-12d3-a456-426655440000/register"), testHandlerOKWildcardBody("uuid"),
		"123e4567-e89b-12d3-a456-426655440000",
	},
	{
		"user/[uuid]{...}/dashboard",
		`/user/[uuid]{^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$}/dashboard`,
		req(http.MethodGet, "/user/123e4567-e89b-12d3-a456-426655440000/dashboard"), testHandlerOKWildcardBody("uuid"),
		"123e4567-e89b-12d3-a456-426655440000",
	},
}

func TestBasicRoutesl(t *testing.T) {
	// Setup router
	r := Default()
	setupRouter(r, basicRoutes)
	// Run cases
	s := httptest.NewServer(r)
	for _, tc := range basicRoutes {
		runTestCase(t, s, tc)
	}
}

func TestComplexRoutes(t *testing.T) {
	// Setup router
	r := Default()
	setupRouter(r, complexRoutes)
	// Run cases
	s := httptest.NewServer(r)
	for _, tc := range complexRoutes {
		runTestCase(t, s, tc)
	}
}

func BenchmarkBasicRoutes(b *testing.B) {
	// Setup router
	r := Default()
	setupRouter(r, basicRoutes)
	// Run cases
	s := httptest.NewServer(r)
	for i := 0; i < b.N; i++ {
		tc := basicRoutes[i%len(basicRoutes)]
		runTestBench(b, s, tc)
	}
	b.StopTimer()
}

func BenchmarkComplexRoutes(b *testing.B) {
	// Setup router
	r := Default()
	setupRouter(r, basicRoutes)
	// Run cases
	s := httptest.NewServer(r)
	for i := 0; i < b.N; i++ {
		tc := basicRoutes[i%len(basicRoutes)]
		runTestBench(b, s, tc)
	}
}
