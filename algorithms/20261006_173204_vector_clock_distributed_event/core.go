package main

import "fmt"

type Relation int

const (
	Equal Relation = iota
	HappensBefore
	HappensAfter
	Concurrent
)

type VectorClock struct {
	clock map[int]int
}

// NewVectorClock creates a vector clock with the given node IDs initialized to zero.
func NewVectorClock(nodes []int) *VectorClock {
	vc := &VectorClock{clock: make(map[int]int, len(nodes))}
	for _, n := range nodes {
		vc.clock[n] = 0
	}
	return vc
}

// Clone returns a deep copy of the vector clock.
func (vc *VectorClock) Clone() *VectorClock {
	newClock := make(map[int]int, len(vc.clock))
	for k, v := range vc.clock {
		newClock[k] = v
	}
	return &VectorClock{clock: newClock}
}

// Increment increases the counter for the specified node.
func (vc *VectorClock) Increment(node int) {
	if _, ok := vc.clock[node]; !ok {
		vc.clock[node] = 0
	}
	vc.clock[node]++
}

// Merge updates the clock to be the element‑wise maximum of itself and another clock.
func (vc *VectorClock) Merge(other *VectorClock) {
	for node, ts := range other.clock {
		if cur, ok := vc.clock[node]; !ok || ts > cur {
			vc.clock[node] = ts
		}
	}
}

// Compare determines the causal relationship between two vector clocks.
// Returns Equal, HappensBefore, HappensAfter, or Concurrent.
func (vc *VectorClock) Compare(other *VectorClock) Relation {
	less, greater := false, false
	// Ensure we consider all nodes present in either clock.
	seen := make(map[int]struct{})
	for node := range vc.clock {
		seen[node] = struct{}{}
	}
	for node := range other.clock {
		seen[node] = struct{}{}
	}
	for node := range seen {
		a := vc.clock[node]
		b := other.clock[node]
		if a < b {
			less = true
		} else if a > b {
			greater = true
		}
	}
	switch {
	case !less && !greater:
		return Equal
	case less && !greater:
		return HappensBefore
	case greater && !less:
		return HappensAfter
	default:
		return Concurrent
	}
}

// String returns a deterministic string representation of the vector clock.
func (vc *VectorClock) String() string {
	// Simple format: {node1:counter1, node2:counter2}
	s := "{"
	first := true
	for node, ts := range vc.clock {
		if !first {
			s += ", "
		}
		s += fmt.Sprintf("%d:%d", node, ts)
		first = false
	}
	s += "}"
	return s
}
