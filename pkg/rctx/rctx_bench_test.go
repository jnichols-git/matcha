package rctx

import (
	"net/http"
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
	req := &http.Request{}
	for i := 0; i < b.N; i++ {
		req = PrepareRequestContext(req, DefaultMaxParams)
		// req = req.WithContext(context.Background())
		req = &http.Request{}
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
