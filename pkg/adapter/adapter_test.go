package adapter

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/jnichols-git/matcha/v2/internal/route"
	"github.com/jnichols-git/matcha/v2/internal/router"
)

/*
This test file is a baseline example of the usage of Adapter. Testing implementations is up to the person
doing the implementing.
*/

type testAdaptable struct {
	Method  string
	Headers map[string][]string
	Body    []byte
}

type testAdapted struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

type testWriter struct {
	Resp *testAdapted
}

func (tw *testWriter) Header() http.Header {
	return tw.Resp.Headers
}

func (tw *testWriter) Write(b []byte) (int, error) {
	if tw.Resp.StatusCode == 0 {
		tw.WriteHeader(200)
	}
	tw.Resp.Body = append(tw.Resp.Body, b...)
	return len(b), nil
}

func (tw *testWriter) WriteHeader(statusCode int) {
	if tw.Resp.StatusCode < 200 {
		tw.Resp.StatusCode = statusCode
	}
}

type testAdapter struct {
}

func (adapter *testAdapter) Adapt(in *testAdaptable) (http.ResponseWriter, *http.Request, *testAdapted, error) {
	if in.Method == "" || in.Headers == nil {
		return nil, nil, nil, errors.New("invalid input testAdaptable")
	}
	target := &testAdapted{}
	req := &http.Request{
		URL:    &url.URL{Path: "/"},
		Method: in.Method,
		Header: in.Headers,
		Body:   io.NopCloser(bytes.NewReader(in.Body)),
	}
	w := &testWriter{
		Resp: target,
	}
	return w, req, target, nil
}

func handleTest(w http.ResponseWriter, req *http.Request) {
	bodyRaw, _ := io.ReadAll(req.Body)
	body := string(bodyRaw)
	if body == "say_ok" {
		w.Write([]byte("ok"))
	} else {
		w.WriteHeader(400)
		w.Write([]byte("not ok"))
	}
}

func testHandler(ctx context.Context, in *testAdaptable) (*testAdapted, error) {
	adapter := &testAdapter{}
	rt := router.Declare(
		router.Default(),
		router.HandleRoute(route.Declare(http.MethodGet, "/"), http.HandlerFunc(handleTest)),
	)
	w, req, out, err := adapter.Adapt(in)
	if err != nil {
		return nil, err
	}
	rt.ServeHTTP(w, req)
	return out, nil
}

func TestAdapter(t *testing.T) {
	in1 := &testAdaptable{
		Method: http.MethodGet,
		Headers: map[string][]string{
			"X-Test": {"test-value"},
		},
		Body: []byte("say_ok"),
	}
	out1, err := testHandler(context.Background(), in1)
	if err != nil {
		t.Error(err)
	}
	if string(out1.Body) != "ok" || out1.StatusCode != 200 {
		t.Errorf("Expected out body 'ok' and status 200, got body '%s' and status %d", out1.Body, out1.StatusCode)
	}
	in2 := &testAdaptable{
		Method: http.MethodGet,
		Headers: map[string][]string{
			"X-Test": {"test-value"},
		},
		Body: []byte("some_other_body"),
	}
	out2, err := testHandler(context.Background(), in2)
	if err != nil {
		t.Error(err)
	}
	if string(out2.Body) != "not ok" || out2.StatusCode != 400 {
		t.Errorf("Expected out body 'not ok' and status 400, got body '%s' and status %d", out2.Body, out2.StatusCode)
	}
}
