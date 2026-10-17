package main

// Node represents an entry in the doubly linked list used by LRUCache.
type Node struct {
	key  int
	prev *Node
	next *Node
}

// LRUCache implements a fixed‑size Least Recently Used cache.
// It stores integer keys (e.g., page numbers) and tracks hits/misses.
type LRUCache struct {
	capacity int
	size     int
	lookup   map[int]*Node
	head     *Node // most recently used
	tail     *Node // least recently used
	hits     int
	misses   int
}

// NewLRUCache creates a new LRUCache with the given capacity.
// Capacity must be non‑negative; a capacity of 0 results in a cache that always misses.
func NewLRUCache(capacity int) *LRUCache {
	if capacity < 0 {
		capacity = 0
	}
	return &LRUCache{
		capacity: capacity,
		lookup:   make(map[int]*Node),
	}
}
