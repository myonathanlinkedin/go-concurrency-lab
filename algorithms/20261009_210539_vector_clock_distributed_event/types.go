package main

// Relation describes the partial order relationship between two vector clocks.
type Relation int

const (
	Concurrent Relation = iota // Neither clock happens before the other.
	HappensBefore              // The left clock happens before the right clock.
	HappensAfter               // The left clock happens after the right clock.
	Equal                      // Both clocks are identical.
)

// VectorClock implements a classic vector clock for distributed event ordering.
type VectorClock struct {
	clock map[int]int // map from process identifier to logical timestamp.
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

// Increment advances the logical time for the given process identifier.
func (vc *VectorClock) Increment(pid int) {
	vc.clock[pid] = vc.clock[pid] + 1
}

// Merge incorporates the maximum timestamps from another vector clock.
func (vc *VectorClock) Merge(other *VectorClock) {
	for pid, ts := range other.clock {
		if cur, ok := vc.clock[pid]; !ok || ts > cur {
			vc.clock[pid] = ts
		}
	}
}

// Compare determines the partial order relationship between this clock and another.
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
			return Concurrent
		}
	}

	switch {
	case less && !greater:
		return HappensBefore
	case greater && !less:
		return HappensAfter
	case !less && !greater:
		return Equal
	default:
		return Concurrent
	}
}
