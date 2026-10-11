package main

import (
	"fmt"
)

// NewCluster creates a cluster with the given mode and three nodes.
func NewCluster(mode Mode) *Cluster {
	nodes := []*Node{
		{id: 0, store: make(map[string]string), alive: true},
		{id: 1, store: make(map[string]string), alive: true},
		{id: 2, store: make(map[string]string), alive: true},
	}
	return &Cluster{
		nodes: nodes,
		mode:  mode,
	}
}

// Partition simulates a network partition by marking a subset of nodes as dead.
// For deterministic testing we kill nodes with highest IDs first.
func (c *Cluster) Partition(deadCount int) {
	if deadCount < 0 || deadCount > len(c.nodes) {
		panic("invalid deadCount")
	}
	c.partitioned = true
	// Kill nodes starting from the highest ID
	for i := len(c.nodes) - 1; i >= 0 && deadCount > 0; i-- {
		if c.nodes[i].alive {
			c.nodes[i].alive = false
			deadCount--
		}
	}
}

// Heal restores all nodes to alive state.
func (c *Cluster) Heal() {
	c.partitioned = false
	for _, n := range c.nodes {
		n.alive = true
	}
}

// aliveNodes returns the slice of currently alive nodes.
func (c *Cluster) aliveNodes() []*Node {
	var alive []*Node
	for _, n := range c.nodes {
		if n.alive {
			alive = append(alive, n)
		}
	}
	return alive
}

// majorityCount returns the minimal number of nodes required for a majority.
func (c *Cluster) majorityCount() int {
	return len(c.nodes)/2 + 1
}

// Write stores the key/value according to the cluster's mode and partition state.
func (c *Cluster) Write(key, value string) error {
	alive := c.aliveNodes()
	switch c.mode {
	case ModeCP:
		if len(alive) < c.majorityCount() {
			return ErrUnavailable
		}
		// Write to all alive nodes to keep them consistent.
		for _, n := range alive {
			n.store[key] = value
		}
		return nil
	case ModeAP:
		if len(alive) == 0 {
			return ErrUnavailable
		}
		// Write to every alive node (no guarantee of consistency).
		for _, n := range alive {
			n.store[key] = value
		}
		return nil
	case ModeCA:
		if c.partitioned {
			return ErrInvalidMode // CA cannot tolerate partitions.
		}
		// All nodes are alive; write to all.
		for _, n := range c.nodes {
			n.store[key] = value
		}
		return nil
	default:
		return ErrInvalidMode
	}
}

// Read retrieves the value for a key according to the cluster's mode and partition state.
func (c *Cluster) Read(key string) (string, error) {
	alive := c.aliveNodes()
	switch c.mode {
	case ModeCP:
		if len(alive) < c.majorityCount() {
			return "", ErrUnavailable
		}
		// For CP we require that a majority agree on the value.
		// Since we write consistently in CP, any alive node suffices.
		for _, n := range alive {
			if v, ok := n.store[key]; ok {
				return v, nil
			}
		}
		return "", ErrNotFound
	case ModeAP:
		if len(alive) == 0 {
			return "", ErrUnavailable
		}
		// Return the first value found (may be stale).
		for _, n := range alive {
			if v, ok := n.store[key]; ok {
				return v, nil
			}
		}
		return "", ErrNotFound
	case ModeCA:
		if c.partitioned {
			return "", ErrInvalidMode
		}
		// All nodes are alive; any node suffices.
		for _, n := range c.nodes {
			if v, ok := n.store[key]; ok {
				return v, nil
			}
		}
		return "", ErrNotFound
	default:
		return "", ErrInvalidMode
	}
}

// DebugString returns a human‑readable snapshot of the cluster (used only for debugging).
func (c *Cluster) DebugString() string {
	s := fmt.Sprintf("Mode:%v Partitioned:%v\n", c.mode, c.partitioned)
	for _, n := range c.nodes {
		status := "down"
		if n.alive {
			status = "up"
		}
		s += fmt.Sprintf("Node %d [%s]: %v\n", n.id, status, n.store)
	}
	return s
}
