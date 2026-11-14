package main

import "sync"

type PNMap map[string]int64

type PNCounter struct {
	nodeID string
	pos    PNMap // increments per node
	neg    PNMap // decrements per node
	mu     sync.RWMutex
}

// NewPNCounter creates a new PN-Counter for the given replica identifier.
func NewPNCounter(nodeID string) *PNCounter {
	c := &PNCounter{
		nodeID: nodeID,
		pos:    make(PNMap),
		neg:    make(PNMap),
	}
	c.pos[nodeID] = 0
	c.neg[nodeID] = 0
	return c
}

// Increment adds a positive delta to this replica's component.
func (c *PNCounter) Increment(delta int64) {
	if delta <= 0 {
		return
	}
	c.mu.Lock()
	c.pos[c.nodeID] += delta
	c.mu.Unlock()
}

// Decrement adds a positive delta to this replica's negative component.
func (c *PNCounter) Decrement(delta int64) {
	if delta <= 0 {
		return
	}
	c.mu.Lock()
	c.neg[c.nodeID] += delta
	c.mu.Unlock()
}

// Value returns the current logical value of the counter (Σpos – Σneg).
func (c *PNCounter) Value() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var sumPos, sumNeg int64
	for _, v := range c.pos {
		sumPos += v
	}
	for _, v := range c.neg {
		sumNeg += v
	}
	return sumPos - sumNeg
}

// Merge integrates another replica's state using element‑wise maximum.
func (c *PNCounter) Merge(other *PNCounter) {
	if other == nil {
		return
	}
	c.mu.Lock()
	other.mu.RLock()
	for id, v := range other.pos {
		if cur, ok := c.pos[id]; !ok || v > cur {
			c.pos[id] = v
		}
	}
	for id, v := range other.neg {
		if cur, ok := c.neg[id]; !ok || v > cur {
			c.neg[id] = v
		}
	}
	other.mu.RUnlock()
	c.mu.Unlock()
}

// Clone produces a deep copy of the counter (useful for isolation in tests).
func (c *PNCounter) Clone() *PNCounter {
	c.mu.RLock()
	defer c.mu.RUnlock()
	clone := &PNCounter{
		nodeID: c.nodeID,
		pos:    make(PNMap),
		neg:    make(PNMap),
	}
	for id, v := range c.pos {
		clone.pos[id] = v
	}
	for id, v := range c.neg {
		clone.neg[id] = v
	}
	return clone
}
