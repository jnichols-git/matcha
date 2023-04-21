package tree

import (
	"net/http"

	"github.com/cloudretic/router/pkg/path"
	"github.com/cloudretic/router/pkg/route"
)

type node struct {
	p        route.Part
	children []*node
	is_root  bool
	leaf_id  int
}

func createRoot() *node {
	return &node{
		p:        nil,
		children: make([]*node, 0),
		is_root:  true,
	}
}

func createNode(p route.Part) *node {
	return &node{
		p:        p,
		children: make([]*node, 0),
	}
}

func (n *node) propogate(ps []route.Part, leaf_id int) {
	if len(ps) == 0 {
		n.leaf_id = leaf_id
		return
	}
	next := ps[0]
	if n.leaf_id == 0 {
		for _, child := range n.children {
			if child.p.Eq(next) {
				child.propogate(ps[1:], leaf_id)
				return
			}
		}
	}
	child := createNode(next)
	child.propogate(ps[1:], leaf_id)
	n.children = append(n.children, child)
}

func (n *node) match(req *http.Request, expr string, last int) int {
	if last == -1 {
		return n.leaf_id
	}
	token, next := path.Next(expr, last)
	ok := n.p.Match(nil, token)
	// head, err := rctx.Head(req)
	if !ok {
		return 0
	} else if n.leaf_id != 0 {
		if route.IsPartialEndPart(n.p) {
			return n.match(req, expr, next)
		} else if next == -1 {
			return n.leaf_id
		} else {
			return 0
		}
	}
	for _, child := range n.children {
		match_leaf_id := child.match(req, expr, next)
		if match_leaf_id != 0 {
			return match_leaf_id
		}
		/*
			if match_leaf_id == 0 {
				rctx.ResetRequestContextHead(req, head)
			} else {
				return match_leaf_id
			}
		*/
	}
	return 0
}

type RouteTree struct {
	methodRoot map[string]*node
	nextId     int
}

func New() *RouteTree {
	return &RouteTree{
		methodRoot: make(map[string]*node),
		nextId:     0,
	}
}

func (rtree *RouteTree) Add(r route.Route) int {
	root, ok := rtree.methodRoot[r.Method()]
	if !ok || root == nil {
		root = createRoot()
		rtree.methodRoot[r.Method()] = root
	}
	rtree.nextId++
	root.propogate(r.Parts(), rtree.nextId)
	return rtree.nextId
}

func (rtree *RouteTree) Match(req *http.Request) int {
	root, ok := rtree.methodRoot[req.Method]
	if !ok || root == nil {
		return 0
	}
	expr := req.URL.Path
	for _, r := range root.children {
		match_leaf_id := r.match(req, expr, 0)
		if match_leaf_id != 0 {
			return match_leaf_id
		}
	}
	return 0
}
