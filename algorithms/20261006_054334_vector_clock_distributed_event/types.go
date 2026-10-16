package main

import "fmt"

type Comparison int

const (
	Before Comparison = -1
	Concurrent Comparison = 0
	After Comparison = 1
	Equal Comparison = 2
)

type VectorClock map[string]int

// Clone creates a deep copy of the vector clock.
func (vc VectorClock) Clone() VectorClock {
	clone := make(VectorClock, len(vc))
	for k, v := range vc {
		clone[k] = v
	}
	return clone
}

// Increment increases the counter for the given node.
func (vc VectorClock) Increment(node string) {
	vc[node] = vc[node] + 1
}

// Merge combines another vector clock into this one, taking the element‑wise maximum.
func (vc VectorClock) Merge(other VectorClock) {
	for node, ts := range other {
		if cur, ok := vc[node]; !ok || ts > cur {
			vc[node] = ts
		}
	}
}

// Compare returns the causal relationship between two vector clocks.
func (vc VectorClock) Compare(other VectorClock) Comparison {
	aBeforeB := false
	bBeforeA := false

	// Union of keys
	keys := make(map[string]struct{})
	for k := range vc {
		keys[k] = struct{}{}
	}
	for k := range other {
		keys[k] = struct{}{}
	}

	for node := range keys {
		va := vc[node]
		vb := other[node]
		if va < vb {
			aBeforeB = true
		} else if va > vb {
			bBeforeA = true
		}
		if aBeforeB && bBeforeA {
			return Concurrent
		}
	}

	switch {
	case !aBeforeB && !bBeforeA:
		return Equal
	case aBeforeB && !bBeforeA:
		return Before
	case !aBeforeB && bBeforeA:
		return After
	default:
		return Concurrent
	}
}

// String returns a deterministic string representation of the vector clock.
func (vc VectorClock) String() string {
	if len(vc) == 0 {
		return "{}"
	}
	// Build deterministic order
	type kv struct {
		k string
		v int
	}
	kvs := make([]kv, 0, len(vc))
	for k, v := range vc {
		kvs = append(kvs, kv{k, v})
	}
	// Simple insertion sort (few entries)
	for i := 1; i < len(kvs); i++ {
		j := i
		for j > 0 && kvs[j-1].k > kvs[j].k {
			kvs[j-1], kvs[j] = kvs[j], kvs[j-1]
			j--
		}
	}
	s := "{"
	for i, kv := range kvs {
		s += fmt.Sprintf("%s:%d", kv.k, kv.v)
		if i < len(kvs)-1 {
			s += " "
		}
	}
	s += "}"
	return s
}

// Event represents a single occurrence in the system with its associated vector clock.
type Event struct {
	ID    string
	Node  string
	Clock VectorClock
}
