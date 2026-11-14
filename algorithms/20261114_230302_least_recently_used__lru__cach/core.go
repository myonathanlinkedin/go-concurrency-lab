package main

type node struct {
	key   int
	value int
	prev  *node
	next  *node
}

// LRUCache implements a fixed‑size least‑recently‑used cache.
type LRUCache struct {
	capacity int
	cache    map[int]*node
	head     *node // dummy head (most recent side)
	tail     *node // dummy tail (least recent side)
	size     int
}

// NewLRUCache creates a new LRU cache with the given capacity.
func NewLRUCache(capacity int) *LRUCache {
	head := &node{}
	tail := &node{}
	head.next = tail
	tail.prev = head
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*node),
		head:     head,
		tail:     tail,
		size:     0,
	}
}

// Get returns the value for the given key and true if the key exists.
// It also marks the key as most recently used.
func (c *LRUCache) Get(key int) (int, bool) {
	if n, ok := c.cache[key]; ok {
		c.moveToFront(n)
		return n.value, true
	}
	return 0, false
}

// Put inserts or updates the key with the supplied value.
// If the cache exceeds its capacity, the least‑recently‑used entry is evicted.
func (c *LRUCache) Put(key int, value int) {
	if c.capacity <= 0 {
		return
	}
	if n, ok := c.cache[key]; ok {
		n.value = value
		c.moveToFront(n)
		return
	}
	n := &node{key: key, value: value}
	c.cache[key] = n
	c.addToFront(n)
	c.size++
	if c.size > c.capacity {
		// Evict least recently used node (node before dummy tail)
		lru := c.tail.prev
		if lru != c.head {
			c.removeNode(lru)
			delete(c.cache, lru.key)
			c.size--
		}
	}
}

// addToFront inserts node n right after the dummy head.
func (c *LRUCache) addToFront(n *node) {
	n.prev = c.head
	n.next = c.head.next
	c.head.next.prev = n
	c.head.next = n
}

// removeNode detaches node n from the linked list.
func (c *LRUCache) removeNode(n *node) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
}

// moveToFront moves an existing node to the most‑recent position.
func (c *LRUCache) moveToFront(n *node) {
	c.removeNode(n)
	c.addToFront(n)
}
