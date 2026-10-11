package main

import "fmt"

type SegmentTree struct {
	n    int           // size of the underlying array (power of two)
	tree []int64       // segment tree nodes storing sums
	lazy []int64       // lazy propagation values
}

// NewSegmentTree constructs a segment tree from the given slice.
// It supports range addition and range sum queries.
func NewSegmentTree(arr []int64) *SegmentTree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	tree := make([]int64, 2*n)
	lazy := make([]int64, 2*n)

	// Initialize leaves
	for i, v := range arr {
		tree[n+i] = v
	}
	// Build internal nodes
	for i := n - 1; i > 0; i-- {
		tree[i] = tree[i<<1] + tree[i<<1|1]
	}
	return &SegmentTree{n: n, tree: tree, lazy: lazy}
}

// apply adds value to node idx covering segment of length segLen.
func (st *SegmentTree) apply(idx int, value int64, segLen int) {
	st.tree[idx] += value * int64(segLen)
	if idx < st.n { // not a leaf
		st.lazy[idx] += value
	}
}

// push propagates pending updates from node idx down to its children.
func (st *SegmentTree) push(idx int, segLen int) {
	if st.lazy[idx] != 0 {
		left := idx << 1
		right := left | 1
		half := segLen >> 1
		st.apply(left, st.lazy[idx], half)
		st.apply(right, st.lazy[idx], half)
		st.lazy[idx] = 0
	}
}

// rangeAdd adds value to all elements in [l, r).
func (st *SegmentTree) rangeAdd(l, r int, value int64) {
	st.rangeAddRec(l, r, value, 1, 0, st.n)
}

// internal recursive helper for range addition.
func (st *SegmentTree) rangeAddRec(l, r int, value int64, idx, nodeL, nodeR int) {
	if l >= nodeR || r <= nodeL { // no overlap
		return
	}
	if l <= nodeL && nodeR <= r { // total overlap
		st.apply(idx, value, nodeR-nodeL)
		return
	}
	// partial overlap
	st.push(idx, nodeR-nodeL)
	mid := (nodeL + nodeR) >> 1
	st.rangeAddRec(l, r, value, idx<<1, nodeL, mid)
	st.rangeAddRec(l, r, value, idx<<1|1, mid, nodeR)
	st.tree[idx] = st.tree[idx<<1] + st.tree[idx<<1|1]
}

// rangeSum returns the sum of elements in [l, r).
func (st *SegmentTree) rangeSum(l, r int) int64 {
	return st.rangeSumRec(l, r, 1, 0, st.n)
}

// internal recursive helper for range sum query.
func (st *SegmentTree) rangeSumRec(l, r int, idx, nodeL, nodeR int) int64 {
	if l >= nodeR || r <= nodeL { // no overlap
		return 0
	}
	if l <= nodeL && nodeR <= r { // total overlap
		return st.tree[idx]
	}
	// partial overlap
	st.push(idx, nodeR-nodeL)
	mid := (nodeL + nodeR) >> 1
	leftSum := st.rangeSumRec(l, r, idx<<1, nodeL, mid)
	rightSum := st.rangeSumRec(l, r, idx<<1|1, mid, nodeR)
	return leftSum + rightSum
}

// Update adds value to the half‑open interval [l, r).
func (st *SegmentTree) Update(l, r int, value int64) {
	if l < 0 || r > st.n || l > r {
		panic(fmt.Sprintf("Update: invalid range [%d,%d)", l, r))
	}
	st.rangeAdd(l, r, value)
}

// Query returns the sum over the half‑open interval [l, r).
func (st *SegmentTree) Query(l, r int) int64 {
	if l < 0 || r > st.n || l > r {
		panic(fmt.Sprintf("Query: invalid range [%d,%d)", l, r))
	}
	return st.rangeSum(l, r)
}
