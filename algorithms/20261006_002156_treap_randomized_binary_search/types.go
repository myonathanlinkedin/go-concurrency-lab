package main

import (
	"math/rand"
	"time"
)

// Node represents a single node in the Treap.
type Node struct {
	key      int
	priority int
	left     *Node
	right    *Node
}

// Treap represents the randomized binary search tree.
type Treap struct {
	root *Node
	rnd  *rand.Rand
}

// NewTreap creates a new empty Treap with a seeded random source.
func NewTreap() *Treap {
	return &Treap{
		root: nil,
		rnd:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Insert adds a key to the Treap. Duplicate keys are ignored.
func (t *Treap) Insert(key int) {
	if t.Search(key) {
		return
	}
	newNode := &Node{
		key:      key,
		priority: t.rnd.Int(),
	}
	t.root = merge(t.root, newNode)
}

// Delete removes a key from the Treap. If the key does not exist, the Treap remains unchanged.
func (t *Treap) Delete(key int) {
	t.root = deleteNode(t.root, key)
}

// Search returns true if the key exists in the Treap.
func (t *Treap) Search(key int) bool {
	return searchNode(t.root, key)
}

// InOrder returns a slice of all keys in ascending order.
func (t *Treap) InOrder() []int {
	var res []int
	inOrder(t.root, &res)
	return res
}
