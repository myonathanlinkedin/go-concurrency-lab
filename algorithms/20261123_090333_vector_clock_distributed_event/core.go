package main

import (
	"fmt"
	"strings"
)

// Relation describes the partial order relationship between two vector clocks.
type Relation int

const (
	RelationConcurrent Relation = iota // Neither clock happened before the other.
	RelationHappensBefore               // The receiver happened before the argument.
	RelationHappensAfter                // The receiver happened after the argument.
	RelationEqual                       // Both clocks are identical.
)

// VectorClock implements a classic vector clock for distributed event ordering.
type VectorClock struct {
	clock map[int]int // map[processID]counter
}

// NewVectorClock creates an empty vector clock.
func NewVectorClock() *VectorClock {
	return &VectorClock{clock: make(map[int]int)}
}

// Clone returns a deep copy of the vector clock.
func (vc *VectorClock) Clone() *VectorClock {
	clone := NewVectorClock()
	for pid, ts := range vc.clock {
		clone.clock[pid] = ts
	}
	return clone
}

// Increment advances the counter for the given process identifier.
func (vc *VectorClock) Increment(pid int) {
	vc.clock[pid] = vc.clock[pid] + 1
}

// Merge incorporates another clock into the receiver, taking the element‑wise maximum.
func (vc *VectorClock) Merge(other *VectorClock) {
	for pid, ts := range other.clock {
		if cur, ok := vc.clock[pid]; !ok || ts > cur {
			vc.clock[pid] = ts
		}
	}
}

// Compare determines the partial‑order relationship between the receiver and another clock.
func (vc *VectorClock) Compare(other *VectorClock) Relation {
	less, greater := false, false

	// Examine all process identifiers present in either clock.
	seen := make(map[int]struct{})
	for pid := range vc.clock {
		seen[pid] = struct{}{}
	}
	for pid := range other.clock {
		seen[pid] = struct{}{}
	}

	for pid := range seen {
		a := vc.clock[pid]
		b := other.clock[pid]
		if a < b {
			less = true
		} else if a > b {
			greater = true
		}
		if less && greater {
			return RelationConcurrent
		}
	}

	switch {
	case less && !greater:
		return RelationHappensBefore
	case greater && !less:
		return RelationHappensAfter
	case !less && !greater:
		return RelationEqual
	default:
		return RelationConcurrent
	}
}

// String returns a deterministic textual representation of the clock.
func (vc *VectorClock) String() string {
	if len(vc.clock) == 0 {
		return "{}"
	}
	parts := make([]string, 0, len(vc.clock))
	for pid, ts := range vc.clock {
		parts = append(parts, fmt.Sprintf("%d:%d", pid, ts))
	}
	// Sorting ensures stable output across runs.
	// Simple bubble sort to avoid importing "sort" (still standard library, but kept minimal).
	for i := 0; i < len(parts)-1; i++ {
		for j := i + 1; j < len(parts); j++ {
			if parts[i] > parts[j] {
				parts[i], parts[j] = parts[j], parts[i]
			}
		}
	}
	return "{" + strings.Join(parts, ",") + "}"
}
