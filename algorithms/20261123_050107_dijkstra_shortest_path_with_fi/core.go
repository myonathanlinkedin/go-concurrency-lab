package main

import (
	"math"
)

// Edge represents a directed weighted edge.
type Edge struct {
	to     int
	weight int64
}

// Graph is an adjacency list where each vertex has a slice of outgoing edges.
type Graph [][]Edge

// ---------- Fibonacci Heap Implementation ----------

// FibNode is a node in the Fibonacci heap.
type FibNode struct {
	key    int64
	value  int
	degree int
	mark   bool
	parent *FibNode
	child  *FibNode
	left   *FibNode
	right  *FibNode
}

// FibHeap represents a Fibonacci heap.
type FibHeap struct {
	min *FibNode
	n   int
}

// NewFibHeap creates an empty Fibonacci heap.
func NewFibHeap() *FibHeap {
	return &FibHeap{}
}

// makeNode creates a singleton circular doubly linked list node.
func makeNode(key int64, value int) *FibNode {
	node := &FibNode{
		key:   key,
		value: value,
	}
	node.left = node
	node.right = node
	return node
}

// Insert adds a new key/value pair to the heap and returns the created node.
func (h *FibHeap) Insert(key int64, value int) *FibNode {
	node := makeNode(key, value)
	h.mergeRootList(node)
	if h.min == nil || node.key < h.min.key {
		h.min = node
	}
	h.n++
	return node
}

// mergeRootList merges a node into the root list.
func (h *FibHeap) mergeRootList(node *FibNode) {
	if h.min == nil {
		h.min = node
		return
	}
	// Insert node to the left of min.
	node.left = h.min.left
	node.right = h.min
	h.min.left.right = node
	h.min.left = node
}

// ExtractMin removes and returns the node with the smallest key.
func (h *FibHeap) ExtractMin() (*FibNode, bool) {
	z := h.min
	if z == nil {
		return nil, false
	}
	// For each child of z, add it to the root list.
	if z.child != nil {
		children := []*FibNode{}
		c := z.child
		for {
			children = append(children, c)
			c = c.right
			if c == z.child {
				break
			}
		}
		for _, child := range children {
			h.mergeRootList(child)
			child.parent = nil
		}
	}
	// Remove z from root list.
	h.removeFromRootList(z)
	if z == z.right { // heap is now empty
		h.min = nil
	} else {
		h.min = z.right
		h.consolidate()
	}
	h.n--
	return z, true
}

// removeFromRootList detaches a node from its circular list.
func (h *FibHeap) removeFromRootList(node *FibNode) {
	node.left.right = node.right
	node.right.left = node.left
	node.left = node
	node.right = node
}

// consolidate merges trees of equal degree.
func (h *FibHeap) consolidate() {
	maxDegree := int(math.Floor(math.Log2(float64(h.n)))) + 1
	A := make([]*FibNode, maxDegree+1)

	roots := []*FibNode{}
	x := h.min
	if x != nil {
		for {
			roots = append(roots, x)
			x = x.right
			if x == h.min {
				break
			}
		}
	}
	for _, w := range roots {
		x = w
		d := x.degree
		for A[d] != nil {
			y := A[d]
			if x.key > y.key {
				x, y = y, x
			}
			h.heapLink(y, x)
			A[d] = nil
			d++
		}
		A[d] = x
	}
	h.min = nil
	for _, node := range A {
		if node != nil {
			if h.min == nil {
				node.left = node
				node.right = node
				h.min = node
			} else {
				h.mergeRootList(node)
				if node.key < h.min.key {
					h.min = node
				}
			}
		}
	}
}

// heapLink makes y a child of x.
func (h *FibHeap) heapLink(y, x *FibNode) {
	// Remove y from root list.
	h.removeFromRootList(y)
	// Make y a child of x.
	y.parent = x
	if x.child == nil {
		y.left = y
		y.right = y
		x.child = y
	} else {
		// Insert y into x's child list.
		y.left = x.child.left
		y.right = x.child
		x.child.left.right = y
		x.child.left = y
	}
	x.degree++
	y.mark = false
}

// DecreaseKey decreases the key of a given node.
func (h *FibHeap) DecreaseKey(x *FibNode, k int64) {
	if k > x.key {
		panic("new key is greater than current key")
	}
	x.key = k
	y := x.parent
	if y != nil && x.key < y.key {
		h.cut(x, y)
		h.cascadingCut(y)
	}
	if h.min == nil || x.key < h.min.key {
		h.min = x
	}
}

// cut removes x from its parent y's child list and adds it to the root list.
func (h *FibHeap) cut(x, y *FibNode) {
	// Remove x from child list of y.
	if x.right == x {
		y.child = nil
	} else {
		if y.child == x {
			y.child = x.right
		}
		x.left.right = x.right
		x.right.left = x.left
	}
	y.degree--
	// Add x to root list.
	x.parent = nil
	x.left = x
	x.right = x
	h.mergeRootList(x)
	x.mark = false
}

// cascadingCut performs a cascading cut operation.
func (h *FibHeap) cascadingCut(y *FibNode) {
	z := y.parent
	if z != nil {
		if !y.mark {
			y.mark = true
		} else {
			h.cut(y, z)
			h.cascadingCut(z)
		}
	}
}

// ---------- Dijkstra's Algorithm ----------

// Dijkstra computes shortest path distances from src to all vertices.
// Returns a slice where dist[v] is the minimum distance or math.MaxInt64 if unreachable.
func Dijkstra(g Graph, src int) []int64 {
	n := len(g)
	const INF int64 = math.MaxInt64
	dist := make([]int64, n)
	for i := range dist {
		dist[i] = INF
	}
	dist[src] = 0

	heap := NewFibHeap()
	nodes := make([]*FibNode, n)
	for v := 0; v < n; v++ {
		nodes[v] = heap.Insert(dist[v], v)
	}

	for heap.n > 0 {
		minNode, ok := heap.ExtractMin()
		if !ok {
			break
		}
		u := minNode.value
		du := minNode.key
		if du == INF {
			break
		}
		for _, e := range g[u] {
			alt := du + e.weight
			if alt < dist[e.to] {
				dist[e.to] = alt
				heap.DecreaseKey(nodes[e.to], alt)
			}
		}
	}
	return dist
}
