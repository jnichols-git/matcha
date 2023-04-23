package tree

import (
	"net/http"
	"testing"

	"github.com/cloudretic/matcha/pkg/rctx"
	"github.com/cloudretic/matcha/pkg/route"
)

func TestTree(t *testing.T) {
	rtree := New()
	rtree.Add(route.Declare(http.MethodGet, "/test/route/[p1]{[a-z]*}"))
	rtree.Add(route.Declare(http.MethodGet, "/test/route/[p2]{[a-zA-Z]*}"))
	rtree.Add(route.Declare(http.MethodGet, "/test/route/[ext]+"))
	rtree.Add(route.Declare(http.MethodGet, "/test"))
	req, _ := http.NewRequest(http.MethodGet, "http://test.com/test/route/lowercase", nil)
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	if leaf_id := rtree.Match(req); leaf_id != 1 {
		t.Errorf("wrong match: %d", leaf_id)
	}
	req, _ = http.NewRequest(http.MethodGet, "http://test.com/test/route/Uppercase", nil)
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	if leaf_id := rtree.Match(req); leaf_id != 2 {
		t.Errorf("wrong match: %d", leaf_id)
	}
	req, _ = http.NewRequest(http.MethodGet, "http://test.com/test/route/longer/request", nil)
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	if leaf_id := rtree.Match(req); leaf_id != 3 {
		t.Errorf("wrong match: %d", leaf_id)
	}
	req, _ = http.NewRequest(http.MethodGet, "http://test.com/test", nil)
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	if leaf_id := rtree.Match(req); leaf_id != 4 {
		t.Errorf("wrong match: %d", leaf_id)
	}
	req, _ = http.NewRequest(http.MethodGet, "http://test.com/notfound", nil)
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	if leaf_id := rtree.Match(req); leaf_id != NO_LEAF_ID {
		t.Error("route shouldn't exist")
	}
	req, _ = http.NewRequest(http.MethodPost, "http://test.com/test/route/lowercase", nil)
	req = rctx.PrepareRequestContext(req, rctx.DefaultMaxParams)
	if leaf_id := rtree.Match(req); leaf_id != NO_LEAF_ID {
		t.Error("route shouldn't exist")
	}
}
