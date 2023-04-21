package tree

import (
	"net/http"

	"github.com/cloudretic/router/pkg/path"
	"github.com/cloudretic/router/pkg/route"
)

const NO_LEAF_ID = int(0)

type node struct {
	p        route.Part
	children []*node
	leaf_id  int
}

func (n *node) isRoot() bool {
	return n.p == nil
}

func (n *node) isLeaf() bool {
	return n.leaf_id != NO_LEAF_ID
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
	if !n.isLeaf() {
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

// match is a routine that traverses a tree of nodes to find the first matching route.
func (n *node) match(expr string, last int) int {
	// If we've reached the end of the expression, return the leaf_id of the current node.
	// This encapsulates several edge cases where it's difficult to know if the routine should return early or not,
	// like with partial leaves.
	if last == -1 {
		return n.leaf_id
	}
	// Get the next token from the path and match it against the route.Part of the current node.
	token, next := path.Next(expr, last)
	ok := n.p.Match(nil, token)
	if !ok {
		// If the part doesn't match, return NO_LEAF_ID.
		return NO_LEAF_ID
	} else if n.isLeaf() {
		// If the part matches...
		if route.IsPartialEndPart(n.p) {
			// ...and the leaf is partial, return the result of recursively matching until termination.
			return n.match(expr, next)
		} else if next == -1 {
			// ...and the route has been exhausted, return the id of the leaf as a successful match.
			return n.leaf_id
		} else {
			// ...and the route has not been exhausted, return NO_LEAF_ID.
			return NO_LEAF_ID
		}
	}
	// Iterate through the children of this node.
	for _, child := range n.children {
		match_leaf_id := child.match(expr, next)
		if match_leaf_id != 0 {
			// If a child matches the entire remaining route, return its leaf_id.
			return match_leaf_id
		}
	}
	// If we reach this point, the entire subtree from this node has been traversed with no match.
	return NO_LEAF_ID
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
		root = createNode(nil)
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
		match_leaf_id := r.match(expr, 0)
		if match_leaf_id != 0 {
			return match_leaf_id
		}
	}
	return 0
}
