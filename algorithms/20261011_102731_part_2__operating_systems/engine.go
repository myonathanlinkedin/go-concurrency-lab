package main

// moveToFront moves the given node to the head of the list (most recently used).
func (c *LRUCache) moveToFront(n *Node) {
	if c.head == n {
		return // already at front
	}
	// Detach n
	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}
	// Update tail if needed
	if c.tail == n {
		c.tail = n.prev
	}
	// Insert at front
	n.prev = nil
	n.next = c.head
	if c.head != nil {
		c.head.prev = n
	}
	c.head = n
	if c.tail == nil {
		c.tail = n
	}
}

// evictTail removes the least recently used node from the cache.
func (c *LRUCache) evictTail() {
	if c.tail == nil {
		return
	}
	evicted := c.tail
	if evicted.prev != nil {
		evicted.prev.next = nil
	}
	c.tail = evicted.prev
	if c.tail == nil {
		c.head = nil
	}
	delete(c.lookup, evicted.key)
	c.size--
}

// Get checks whether the key is present in the cache.
// It returns true on a hit (and updates recency) and false on a miss.
func (c *LRUCache) Get(key int) bool {
	if node, ok := c.lookup[key]; ok {
		c.hits++
		c.moveToFront(node)
		return true
	}
	c.misses++
	return false
}

// Put inserts the key into the cache.
// If the key already exists, it is promoted to most recent.
// If insertion would exceed capacity, the least recently used entry is evicted.
func (c *LRUCache) Put(key int) {
	if node, ok := c.lookup[key]; ok {
		// Already present – just promote.
		c.moveToFront(node)
		return
	}
	// New entry.
	if c.capacity == 0 {
		// Cache disabled – count as miss.
		c.misses++
		return
	}
	if c.size == c.capacity {
		c.evictTail()
	}
	newNode := &Node{key: key}
	c.lookup[key] = newNode
	c.moveToFront(newNode)
	c.size++
}

// Stats returns the accumulated hit and miss counts.
func (c *LRUCache) Stats() (hits int, misses int) {
	return c.hits, c.misses
}
