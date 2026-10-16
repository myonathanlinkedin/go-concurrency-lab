package main

import (
	"fmt"
	"sort"
)

type VectorClock map[int]int

func NewVectorClock() VectorClock {
	return make(VectorClock)
}

func (vc VectorClock) Clone() VectorClock {
	clone := make(VectorClock, len(vc))
	for k, v := range vc {
		clone[k] = v
	}
	return clone
}

func (vc VectorClock) Inc(node int) {
	vc[node] = vc[node] + 1
}

func (vc VectorClock) Merge(other VectorClock) {
	for k, v := range other {
		if cur, ok := vc[k]; !ok || v > cur {
			vc[k] = v
		}
	}
}

// Relation represents the ordering relationship between two vector clocks.
type Relation int

const (
	Equal Relation = iota
	HappensBefore
	HappensAfter
	Concurrent
)

func (vc VectorClock) Compare(other VectorClock) Relation {
	less, greater := false, false
	keys := make(map[int]struct{})
	for k := range vc {
		keys[k] = struct{}{}
	}
	for k := range other {
		keys[k] = struct{}{}
	}
	for k := range keys {
		a := vc[k]
		b := other[k]
		if a < b {
			less = true
		} else if a > b {
			greater = true
		}
		if less && greater {
			return Concurrent
		}
	}
	if !less && !greater {
		return Equal
	}
	if less && !greater {
		return HappensBefore
	}
	return HappensAfter
}

func (r Relation) String() string {
	switch r {
	case Equal:
		return "Equal"
	case HappensBefore:
		return "HappensBefore"
	case HappensAfter:
		return "HappensAfter"
	case Concurrent:
		return "Concurrent"
	default:
		return "Unknown"
	}
}

func (vc VectorClock) String() string {
	if len(vc) == 0 {
		return "{}"
	}
	keys := make([]int, 0, len(vc))
	for k := range vc {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	s := "{"
	for i, k := range keys {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%d:%d", k, vc[k])
	}
	s += "}"
	return s
}
