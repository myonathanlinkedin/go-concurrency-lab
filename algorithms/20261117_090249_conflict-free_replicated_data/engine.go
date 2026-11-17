package main

import (
	"fmt"
)

// Increment adds delta (must be >0) to the positive G-Counter for this replica.
func (c *PNCounter) Increment(delta int64) {
	if delta <= 0 {
		panic("increment delta must be positive")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.inc[c.id] = c.inc[c.id] + uint64(delta)
}

// Decrement adds delta (must be >0) to the negative G-Counter for this replica.
func (c *PNCounter) Decrement(delta int64) {
	if delta <= 0 {
		panic("decrement delta must be positive")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dec[c.id] = c.dec[c.id] + uint64(delta)
}

// Value returns the current logical value of the PN-Counter.
func (c *PNCounter) Value() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var posSum uint64
	var negSum uint64
	for _, v := range c.inc {
		posSum += v
	}
	for _, v := range c.dec {
		negSum += v
	}
	// Convert to signed int64; overflow is undefined behavior for extremely large values,
	// but for typical usage it stays within int64 range.
	if posSum >= negSum {
		return int64(posSum - negSum)
	}
	return -int64(negSum - posSum)
}

// Merge integrates another PN-Counter's state into this one.
// The operation is commutative, associative, and idempotent.
func (c *PNCounter) Merge(other *PNCounter) {
	if other == nil {
		panic("cannot merge nil PN-Counter")
	}
	// Acquire locks in a deterministic order to avoid deadlock.
	// Use address ordering.
	var first, second *PNCounter
	if fmt.Sprintf("%p", c) < fmt.Sprintf("%p", other) {
		first, second = c, other
	} else {
		first, second = other, c
	}
	first.mu.Lock()
	second.mu.Lock()
	defer second.mu.Unlock()
	defer first.mu.Unlock()

	// Merge positive map.
	for id, val := range other.inc {
		if cur, ok := c.inc[id]; !ok || val > cur {
			c.inc[id] = val
		}
	}
	// Merge negative map.
	for id, val := range other.dec {
		if cur, ok := c.dec[id]; !ok || val > cur {
			c.dec[id] = val
		}
	}
}

// Clone creates a deep copy of the PN-Counter.
func (c *PNCounter) Clone() *PNCounter {
	c.mu.RLock()
	defer c.mu.RUnlock()
	clone := NewPNCounter(c.id)
	for id, v := range c.inc {
		clone.inc[id] = v
	}
	for id, v := range c.dec {
		clone.dec[id] = v
	}
	return clone
}
