package main

import "sync"

type PNCounter struct {
	mu  sync.RWMutex
	inc map[string]int64 // positive increments per replica
	dec map[string]int64 // decrements per replica
}

// NewPNCounter creates an empty PN-Counter.
func NewPNCounter() *PNCounter {
	return &PNCounter{
		inc: make(map[string]int64),
		dec: make(map[string]int64),
	}
}

// Increment adds n (>=0) to the counter for the given replica identifier.
func (c *PNCounter) Increment(id string, n int64) {
	if n < 0 {
		panic("increment delta must be non‑negative")
	}
	c.mu.Lock()
	c.inc[id] += n
	c.mu.Unlock()
}

// Decrement adds n (>=0) to the decrement component for the given replica identifier.
func (c *PNCounter) Decrement(id string, n int64) {
	if n < 0 {
		panic("decrement delta must be non‑negative")
	}
	c.mu.Lock()
	c.dec[id] += n
	c.mu.Unlock()
}

// Value returns the current logical value of the counter.
func (c *PNCounter) Value() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var sum int64
	for _, v := range c.inc {
		sum += v
	}
	for _, v := range c.dec {
		sum -= v
	}
	return sum
}

// Merge integrates another PN‑Counter's state, keeping the maximum per‑replica component.
func (c *PNCounter) Merge(other *PNCounter) {
	c.mu.Lock()
	other.mu.RLock()
	for id, v := range other.inc {
		if v > c.inc[id] {
			c.inc[id] = v
		}
	}
	for id, v := range other.dec {
		if v > c.dec[id] {
			c.dec[id] = v
		}
	}
	other.mu.RUnlock()
	c.mu.Unlock()
}

// Clone produces a deep copy of the counter.
func (c *PNCounter) Clone() *PNCounter {
	c.mu.RLock()
	defer c.mu.RUnlock()
	clone := NewPNCounter()
	for id, v := range c.inc {
		clone.inc[id] = v
	}
	for id, v := range c.dec {
		clone.dec[id] = v
	}
	return clone
}
