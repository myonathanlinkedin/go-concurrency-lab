package main

import (
	"fmt"
	"sort"
	"strings"
)

// VectorClock implements a classic vector clock for distributed event ordering.
type VectorClock struct {
	// counters maps a process identifier to its logical counter.
	counters map[int]int
}

// NewVectorClock creates a new empty vector clock.
// The optional initialSize hints the initial map capacity.
func NewVectorClock(initialSize int) *VectorClock {
	if initialSize < 0 {
		initialSize = 0
	}
	return &VectorClock{
		counters: make(map[int]int, initialSize),
	}
}

// Clone returns a deep copy of the vector clock.
func (vc *VectorClock) Clone() *VectorClock {
	clone := NewVectorClock(len(vc.counters))
	for pid, cnt := range vc.counters {
		clone.counters[pid] = cnt
	}
	return clone
}

// Increment advances the logical time of the given process identifier.
// If the pid is unseen, it is added with a counter of 1.
func (vc *VectorClock) Increment(pid int) {
	vc.counters[pid] = vc.counters[pid] + 1
}

// Update merges another vector clock into this one, taking the element‑wise maximum.
// The method mutates the receiver.
func (vc *VectorClock) Update(other *VectorClock) {
	for pid, cnt := range other.counters {
		if current, ok := vc.counters[pid]; !ok || cnt > current {
			vc.counters[pid] = cnt
		}
	}
}

// Merge returns a new vector clock that is the element‑wise maximum of the two inputs.
// Neither input is mutated.
func (vc *VectorClock) Merge(other *VectorClock) *VectorClock {
	merged := vc.Clone()
	merged.Update(other)
	return merged
}

// Compare returns a relation code between this clock and another:
//
//   -1 : this happens‑before other (strictly less)
//    0 : clocks are equal
//    1 : this happens‑after other (strictly greater)
//    2 : concurrent (incomparable)
func (vc *VectorClock) Compare(other *VectorClock) int {
	less, greater := false, false

	// Examine all keys present in either clock.
	seen := make(map[int]struct{}, len(vc.counters)+len(other.counters))
	for pid := range vc.counters {
		seen[pid] = struct{}{}
	}
	for pid := range other.counters {
		seen[pid] = struct{}{}
	}

	for pid := range seen {
		a := vc.counters[pid]
		b := other.counters[pid]
		if a < b {
			less = true
		} else if a > b {
			greater = true
		}
		if less && greater {
			return 2 // concurrent
		}
	}
	switch {
	case less && !greater:
		return -1
	case greater && !less:
		return 1
	default:
		return 0
	}
}

// String returns a deterministic textual representation of the clock.
func (vc *VectorClock) String() string {
	if len(vc.counters) == 0 {
		return "{}"
	}
	keys := make([]int, 0, len(vc.counters))
	for pid := range vc.counters {
		keys = append(keys, pid)
	}
	sort.Ints(keys)
	parts := make([]string, 0, len(keys))
	for _, pid := range keys {
		parts = append(parts, fmt.Sprintf("%d:%d", pid, vc.counters[pid]))
	}
	return "{" + strings.Join(parts, ",") + "}"
}
