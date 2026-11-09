package main

import (
	"fmt"
	"math"
	"sync"
)

// Node represents an element in the optimistic lock‑coupled list.
type Node struct {
	key    int
	value  interface{}
	next   *Node
	marked bool          // logical deletion flag
	lock   sync.RWMutex // protects next, marked, and value
}

// List is a sorted singly linked list supporting safe optimistic lock coupling.
type List struct {
	head *Node // sentinel node with key = math.MinInt
}

// NewList creates an empty list with a sentinel head node.
func NewList() *List {
	sentinel := &Node{key: math.MinInt}
	return &List{head: sentinel}
}

// find traverses the list without acquiring any locks and returns the predecessor
// and current node where current.key >= key (or nil if end of list).
func (l *List) find(key int) (pred, curr *Node) {
	pred = l.head
	curr = pred.next
	for curr != nil && curr.key < key {
		pred = curr
		curr = curr.next
	}
	return
}

// Get returns the value associated with key, or (nil, false) if absent.
// It uses optimistic reads followed by validation under read locks.
func (l *List) Get(key int) (interface{}, bool) {
	for {
		pred, curr := l.find(key)
		if curr == nil || curr.key != key {
			return nil, false
		}
		// Acquire read locks for validation.
		pred.lock.RLock()
		curr.lock.RLock()
		if pred.marked || curr.marked || pred.next != curr {
			// Validation failed; retry.
			curr.lock.RUnlock()
			pred.lock.RUnlock()
			continue
		}
		val := curr.value
		curr.lock.RUnlock()
		pred.lock.RUnlock()
		return val, true
	}
}

// Insert adds a key/value pair to the list if the key is not already present.
// Returns true if insertion succeeded, false if the key already existed (value updated).
func (l *List) Insert(key int, value interface{}) bool {
	for {
		pred, curr := l.find(key)

		// Acquire exclusive locks in order: predecessor then current (if not nil).
		pred.lock.Lock()
		if curr != nil {
			curr.lock.Lock()
		}

		// Validate adjacency and that neither node is marked.
		if pred.marked || (curr != nil && curr.marked) || pred.next != curr {
			if curr != nil {
				curr.lock.Unlock()
			}
			pred.lock.Unlock()
			continue // retry
		}

		if curr != nil && curr.key == key {
			// Key exists: update value.
			curr.value = value
			if curr != nil {
				curr.lock.Unlock()
			}
			pred.lock.Unlock()
			return false
		}

		// Insert new node between pred and curr.
		newNode := &Node{
			key:   key,
			value: value,
			next:  curr,
		}
		pred.next = newNode

		if curr != nil {
			curr.lock.Unlock()
		}
		pred.lock.Unlock()
		return true
	}
}

// Delete removes the node with the given key if present.
// Returns true if deletion succeeded, false otherwise.
func (l *List) Delete(key int) bool {
	for {
		pred, curr := l.find(key)
		if curr == nil || curr.key != key {
			return false // not found
		}

		// Acquire exclusive locks.
		pred.lock.Lock()
		curr.lock.Lock()

		// Validate adjacency and that nodes are not already marked.
		if pred.marked || curr.marked || pred.next != curr {
			curr.lock.Unlock()
			pred.lock.Unlock()
			continue // retry
		}

		// Logical deletion.
		curr.marked = true
		// Physical removal.
		pred.next = curr.next

		curr.lock.Unlock()
		pred.lock.Unlock()
		return true
	}
}

// DebugString returns a string representation of the list (for debugging).
func (l *List) DebugString() string {
	var s string
	for n := l.head.next; n != nil; n = n.next {
		s += fmt.Sprintf("%d ", n.key)
	}
	return s
}
