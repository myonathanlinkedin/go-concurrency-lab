package main

import (
	"sort"
	"sync"
)

// LamportClock implements a thread‑safe logical clock.
type LamportClock struct {
	mu      sync.Mutex
	counter uint64
}

// NewLamportClock creates a clock starting at zero.
func NewLamportClock() *LamportClock {
	return &LamportClock{counter: 0}
}

// Tick increments the local counter and returns the new timestamp.
func (lc *LamportClock) Tick() uint64 {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.counter++
	return lc.counter
}

// Update merges a received timestamp into the local clock.
// It sets counter = max(counter, received) + 1 and returns the new timestamp.
func (lc *LamportClock) Update(received uint64) uint64 {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	if received > lc.counter {
		lc.counter = received
	}
	lc.counter++
	return lc.counter
}

// Value returns the current counter without modification.
func (lc *LamportClock) Value() uint64 {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return lc.counter
}

// Compare returns -1 if a<b, 0 if a==b, 1 if a>b.
func Compare(a, b uint64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// Event represents an action annotated with a Lamport timestamp.
type Event struct {
	ID        string
	Timestamp uint64
}

// Engine maintains a Lamport clock and an ordered log of events.
type Engine struct {
	clock *LamportClock
	mu    sync.Mutex
	log   []Event
}

// NewEngine creates a fresh Engine with its own clock.
func NewEngine() *Engine {
	return &Engine{
		clock: NewLamportClock(),
		log:   make([]Event, 0, 16),
	}
}

// LocalEvent records a new local event, advances the clock, and stores it.
func (e *Engine) LocalEvent(id string) Event {
	ts := e.clock.Tick()
	ev := Event{ID: id, Timestamp: ts}
	e.mu.Lock()
	e.log = append(e.log, ev)
	e.mu.Unlock()
	return ev
}

// RemoteEvent processes an event received from another process.
// It updates the clock with the remote timestamp and records the event.
func (e *Engine) RemoteEvent(id string, remoteTS uint64) Event {
	ts := e.clock.Update(remoteTS)
	ev := Event{ID: id, Timestamp: ts}
	e.mu.Lock()
	e.log = append(e.log, ev)
	e.mu.Unlock()
	return ev
}

// Log returns a copy of the event log sorted by timestamp.
func (e *Engine) Log() []Event {
	e.mu.Lock()
	defer e.mu.Unlock()
	cpy := make([]Event, len(e.log))
	copy(cpy, e.log)
	sort.Slice(cpy, func(i, j int) bool {
		return cpy[i].Timestamp < cpy[j].Timestamp
	})
	return cpy
}

// Clear empties the event log.
func (e *Engine) Clear() {
	e.mu.Lock()
	e.log = e.log[:0]
	e.mu.Unlock()
}
