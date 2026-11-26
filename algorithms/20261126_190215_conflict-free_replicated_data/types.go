package main

// PN-Counter is a Conflict-free Replicated Data Type that supports both
// increment and decrement operations. It is composed of two grow-only
// counters: one for positive increments (P) and one for negative increments (N).
type PNCounter struct {
	// replicaID uniquely identifies the replica that owns this instance.
	replicaID string
	// p stores the positive (increment) component: map[replicaID]counter.
	p map[string]uint64
	// n stores the negative (decrement) component: map[replicaID]counter.
	n map[string]uint64
}

// NewPNCounter creates a new PN-Counter for the given replica identifier.
func NewPNCounter(replicaID string) *PNCounter {
	return &PNCounter{
		replicaID: replicaID,
		p:         make(map[string]uint64),
		n:         make(map[string]uint64),
	}
}

// Increment adds delta to the positive component for this replica.
func (c *PNCounter) Increment(delta uint64) {
	if delta == 0 {
		return
	}
	c.p[c.replicaID] = c.p[c.replicaID] + delta
}

// Decrement adds delta to the negative component for this replica.
func (c *PNCounter) Decrement(delta uint64) {
	if delta == 0 {
		return
	}
	c.n[c.replicaID] = c.n[c.replicaID] + delta
}

// Value returns the current logical value of the counter as a signed integer.
func (c *PNCounter) Value() int64 {
	var sumP uint64
	for _, v := range c.p {
		sumP += v
	}
	var sumN uint64
	for _, v := range c.n {
		sumN += v
	}
	// Convert to signed int64; the subtraction may be negative.
	return int64(sumP) - int64(sumN)
}

// Merge incorporates the state of another PN-Counter into this one.
// The operation is commutative, associative, and idempotent.
func (c *PNCounter) Merge(other *PNCounter) {
	// Merge positive component.
	for id, val := range other.p {
		if cur, ok := c.p[id]; !ok || val > cur {
			c.p[id] = val
		}
	}
	// Merge negative component.
	for id, val := range other.n {
		if cur, ok := c.n[id]; !ok || val > cur {
			c.n[id] = val
		}
	}
}
