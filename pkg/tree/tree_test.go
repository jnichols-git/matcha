package tree

import (
	"net/http"
	"testing"

	"github.com/cloudretic/router/pkg/rctx"
	"github.com/cloudretic/router/pkg/route"
)

func TestTree(t *testing.T) {
	rtree := New()
	rtree.Add(route.Declare(http.MethodGet, "/test/route/[p1]{[a-z]*}"))
	rtree.Add(route.Declare(http.MethodGet, "/test/route/[p2]{[a-zA-Z]*}"))
	req, _ := http.NewRequest(http.MethodGet, "http://test.com/test/route/lowercase", nil)
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	if leaf_id := rtree.Match(req); leaf_id == 0 {
		t.Error("nil request")
	}
	req, _ = http.NewRequest(http.MethodGet, "http://test.com/test/route/Uppercase", nil)
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	if leaf_id := rtree.Match(req); leaf_id == 0 {
		t.Error("nil request")
	}
}
