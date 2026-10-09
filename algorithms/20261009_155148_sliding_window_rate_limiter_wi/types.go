package main

import (
	"sync"
	"time"
)

// Config holds the immutable configuration for a RateLimiter.
type Config struct {
	Capacity    uint64        // Maximum number of tokens in the bucket.
	RefillRate  float64       // Tokens added per second (may be fractional).
	WindowSize  time.Duration // Sliding window duration for request counting.
	MaxRequests uint64        // Maximum allowed requests within WindowSize.
}

// RateLimiter implements a thread‑safe sliding‑window rate limiter backed by a
// token bucket. It is suitable for distributed scenarios where multiple
// instances can merge their state deterministically.
type RateLimiter struct {
	mu sync.Mutex

	cfg Config

	// Token bucket state.
	tokens       float64   // Current token count (may be fractional).
	lastRefill   time.Time // Time of the last token refill.

	// Sliding window state.
	requests []time.Time // Timestamps of recent allowed requests (ordered).
}
