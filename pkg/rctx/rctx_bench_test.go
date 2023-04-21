package rctx

import (
	"context"
	"net/http"
	"testing"
)

func BenchmarkPrepare(b *testing.B) {
	req := &http.Request{}
	for i := 0; i < b.N; i++ {
		req = PrepareRequestContext(req, DefaultMaxParams)
		req = req.WithContext(context.Background())
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
