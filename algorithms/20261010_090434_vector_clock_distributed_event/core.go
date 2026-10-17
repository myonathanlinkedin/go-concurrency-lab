package main

import "fmt"

// VectorClock implements a classic vector clock for distributed event ordering.
// It is safe for single‑goroutine use; concurrency control is left to callers.
type VectorClock struct {
	// entries maps a node identifier to its logical counter.
	entries map[int]int
}

// NewVectorClock creates a new empty vector clock.
func NewVectorClock() *VectorClock {
	return &VectorClock{entries: make(map[int]int)}
}

// Clone returns a deep copy of the vector clock.
func (vc *VectorClock) Clone() *VectorClock {
	clone := NewVectorClock()
	for k, v := range vc.entries {
		clone.entries[k] = v
	}
	return clone
}

// Inc increments the logical counter for the given node.
// If the node does not yet exist in the clock, it is added with value 1.
func (vc *VectorClock) Inc(node int) {
	vc.entries[node] = vc.entries[node] + 1
}

// Get returns the counter value for the given node.
// Missing nodes are interpreted as 0.
func (vc *VectorClock) Get(node int) int {
	return vc.entries[node]
}

// Merge updates the receiver to be the element‑wise maximum of itself and other.
// This corresponds to the “receive” operation in vector‑clock literature.
func (vc *VectorClock) Merge(other *VectorClock) {
	for node, val := range other.entries {
		if cur, ok := vc.entries[node]; !ok || val > cur {
			vc.entries[node] = val
		}
	}
}

// Equals reports whether vc and other have identical entries.
func (vc *VectorClock) Equals(other *VectorClock) bool {
	if len(vc.entries) != len(other.entries) {
		return false
	}
	for node, val := range vc.entries {
		if other.entries[node] != val {
			return false
		}
	}
	return true
}

// HappensBefore returns true iff vc is strictly less than other
// (i.e., vc <= other for all nodes and vc != other).
func (vc *VectorClock) HappensBefore(other *VectorClock) bool {
	lessOrEqual := true
	strict := false

	// Check all nodes present in either clock.
	checked := make(map[int]struct{})
	for node, val := range vc.entries {
		otherVal := other.entries[node]
		if val > otherVal {
			lessOrEqual = false
			break
		}
		if val < otherVal {
			strict = true
		}
		checked[node] = struct{}{}
	}
	if !lessOrEqual {
		return false
	}
	for node, val := range other.entries {
		if _, seen := checked[node]; seen {
			continue
		}
		if val > 0 {
			// vc has implicit 0 for this node, other has >0 => strict
			strict = true
		}
	}
	return lessOrEqual && strict
}

// Concurrent returns true iff vc and other are incomparable
// (neither happens‑before the other).
func (vc *VectorClock) Concurrent(other *VectorClock) bool {
	return !vc.HappensBefore(other) && !other.HappensBefore(vc) && !vc.Equals(other)
}

// String returns a deterministic textual representation of the clock.
// Nodes are sorted for readability.
func (vc *VectorClock) String() string {
	if len(vc.entries) == 0 {
		return "{}"
	}
	// Collect keys and sort them.
	keys := make([]int, 0, len(vc.entries))
	for k := range vc.entries {
		keys = append(keys, k)
	}
	// Simple insertion sort to avoid importing "sort" (still standard, but we keep it minimal).
	for i := 1; i < len(keys); i++ {
		j := i
		for j > 0 && keys[j-1] > keys[j] {
			keys[j-1], keys[j] = keys[j], keys[j-1]
			j--
		}
	}
	// Build string.
	s := "{"
	for i, k := range keys {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%d:%d", k, vc.entries[k])
	}
	s += "}"
	return s
}
