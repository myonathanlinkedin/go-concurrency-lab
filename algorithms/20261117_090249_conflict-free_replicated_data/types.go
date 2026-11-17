package main

import "sync"

// PNMap is a grow-only map from replica IDs to unsigned counters.
type PNMap map[string]uint64

// PNCounter implements a Conflict-free Replicated Data Type (PN-Counter).
// It maintains two G-Counters: one for increments (P) and one for decrements (N).
type PNCounter struct {
	mu   sync.RWMutex
	id   string
	inc  PNMap // positive increments
	dec  PNMap // negative decrements
}

// NewPNCounter creates a new PN-Counter for the given replica identifier.
func NewPNCounter(id string) *PNCounter {
	return &PNCounter{
		id:  id,
		inc: make(PNMap),
		dec: make(PNMap),
	}
}
