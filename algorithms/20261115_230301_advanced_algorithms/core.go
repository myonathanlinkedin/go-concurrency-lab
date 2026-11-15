package main

import "fmt"

type SegmentTree struct {
	n    int
	tree []int64
	lazy []int64
}

// NewSegmentTree builds a segment tree for the given array.
func NewSegmentTree(arr []int64) *SegmentTree {
	n := len(arr)
	size := 4 * n // sufficient for segment tree storage
	st := &SegmentTree{
		n:    n,
		tree: make([]int64, size),
		lazy: make([]int64, size),
	}
	st.build(1, 0, n-1, arr)
	return st
}

// build recursively constructs the tree.
func (st *SegmentTree) build(node, l, r int, arr []int64) {
	if l == r {
		st.tree[node] = arr[l]
		return
	}
	mid := (l + r) / 2
	st.build(node*2, l, mid, arr)
	st.build(node*2+1, mid+1, r, arr)
	st.tree[node] = st.tree[node*2] + st.tree[node*2+1]
}

// push propagates pending lazy updates to children.
func (st *SegmentTree) push(node, l, r int) {
	if st.lazy[node] != 0 {
		add := st.lazy[node]
		st.tree[node] += add * int64(r-l+1)
		if l != r { // not a leaf
			st.lazy[node*2] += add
			st.lazy[node*2+1] += add
		}
		st.lazy[node] = 0
	}
}

// UpdateRange adds delta to all elements in [ql,qr] (inclusive).
func (st *SegmentTree) UpdateRange(ql, qr int, delta int64) {
	if ql < 0 || qr >= st.n || ql > qr {
		panic("UpdateRange: invalid range")
	}
	st.update(1, 0, st.n-1, ql, qr, delta)
}

// internal recursive update.
func (st *SegmentTree) update(node, l, r, ql, qr int, delta int64) {
	st.push(node, l, r)
	if r < ql || l > qr {
		return
	}
	if ql <= l && r <= qr {
		st.lazy[node] += delta
		st.push(node, l, r)
		return
	}
	mid := (l + r) / 2
	st.update(node*2, l, mid, ql, qr, delta)
	st.update(node*2+1, mid+1, r, ql, qr, delta)
	st.tree[node] = st.tree[node*2] + st.tree[node*2+1]
}

// QueryRange returns the sum of elements in [ql,qr] (inclusive).
func (st *SegmentTree) QueryRange(ql, qr int) int64 {
	if ql < 0 || qr >= st.n || ql > qr {
		panic("QueryRange: invalid range")
	}
	return st.query(1, 0, st.n-1, ql, qr)
}

// internal recursive query.
func (st *SegmentTree) query(node, l, r, ql, qr int) int64 {
	st.push(node, l, r)
	if r < ql || l > qr {
		return 0
	}
	if ql <= l && r <= qr {
		return st.tree[node]
	}
	mid := (l + r) / 2
	leftSum := st.query(node*2, l, mid, ql, qr)
	rightSum := st.query(node*2+1, mid+1, r, ql, qr)
	return leftSum + rightSum
}

// DebugString returns a string representation of the underlying array (for debugging).
func (st *SegmentTree) DebugString() string {
	// Reconstruct the array by querying each position.
	result := make([]int64, st.n)
	for i := 0; i < st.n; i++ {
		result[i] = st.QueryRange(i, i)
	}
	return fmt.Sprintf("%v", result)
}
