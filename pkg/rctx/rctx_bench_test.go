package rctx

import (
	"net/http"
	"net/url"
	"testing"
)

func use(any) {}

func BenchmarkNewParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		params := newParams(DefaultMaxParams)
		use(params)
	}
}

func BenchmarkPrepare(b *testing.B) {
	url, _ := url.Parse("/")
	req := &http.Request{
		URL: url,
	}
	for i := 0; i < b.N; i++ {
		req = PrepareRequestContext(req, DefaultMaxParams)
		//req = req.WithContext(context.Background())
		// req = req.WithContext(context.Background())
		req = &http.Request{
			URL: url,
		}
	}
}

func BenchmarkSetGetSingleParam(b *testing.B) {
	req := &http.Request{}
	req = PrepareRequestContext(req, DefaultMaxParams)
	for i := 0; i < b.N; i++ {
		SetParam(req.Context(), "paramKey", "paramVal")
		GetParam(req.Context(), "paramKey")
		ResetRequestContext(req)
	}
}
