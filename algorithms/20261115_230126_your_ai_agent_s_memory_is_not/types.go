package main

import (
	"sync"
)

// valueEntry holds a stored value together with its version.
type valueEntry struct {
	value   string
	version uint64
}

// Store provides a thread‑safe key‑value store with version tracking.
// Each write increments a global monotonically increasing version counter.
// The version associated with a key can be used to detect stale reads.
type Store struct {
	mu            sync.RWMutex
	data          map[string]valueEntry
	globalVersion uint64
}

// NewStore creates an empty Store ready for use.
func NewStore() *Store {
	return &Store{
		data: make(map[string]valueEntry),
	}
}
