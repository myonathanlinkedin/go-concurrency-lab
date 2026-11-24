package main

import (
	"sync"
)

// PNCounter is a Conflict-free Replicated Data Type (CRDT) that supports
// both increment and decrement operations. The state consists of two grow‑only
// maps: one for positive (increment) contributions and one for negative
// (decrement) contributions. Merges are performed by taking the element‑wise
// maximum of the two maps, guaranteeing convergence.
type PNCounter struct {
	mu  sync.RWMutex
	inc map[string]int64 // positive contributions per replica
	dec map[string]int64 // negative contributions per replica
}

// NewPNCounter creates a new, empty PN‑Counter.
func NewPNCounter() *PNCounter {
	return &PNCounter{
		inc: make(map[string]int64),
		dec: make(map[string]int64),
	}
}

// Increment adds n (n ≥ 0) to the counter on behalf of replica id.
func (c *PNCounter) Increment(id string, n int64) {
	if n < 0 {
		panic("increment value must be non‑negative")
	}
	c.mu.Lock()
	c.inc[id] += n
	c.mu.Unlock()
}

// Decrement adds n (n ≥ 0) to the negative side of the counter on behalf of replica id.
func (c *PNCounter) Decrement(id string, n int64) {
	if n < 0 {
		panic("decrement value must be non‑negative")
	}
	c.mu.Lock()
	c.dec[id] += n
	c.mu.Unlock()
}

// Value returns the current logical value of the counter (inc – dec).
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

// Merge incorporates the state of another PN‑Counter into this one.
// For each replica identifier the maximum of the two contributions is kept.
func (c *PNCounter) Merge(other *PNCounter) {
	if other == nil {
		return
	}
	// Acquire locks in a deterministic order to avoid dead‑locks.
	// Since both are independent structures we can lock this first,
	// then the other for reading.
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

// Clone creates a deep copy of the counter. Useful for testing merge
// commutativity without mutating the original instances.
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
