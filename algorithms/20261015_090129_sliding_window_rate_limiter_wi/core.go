package main

import (
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (rc *realClock) Now() time.Time { return time.Now() }

type TokenBucket struct {
	capacity   float64
	refillRate float64 // tokens per second
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
	clock      Clock
}

// NewTokenBucket creates a token bucket with given capacity and refill rate (tokens per second).
// If clk is nil, the real system clock is used.
func NewTokenBucket(capacity int, refillRate float64, clk Clock) *TokenBucket {
	if clk == nil {
		clk = &realClock{}
	}
	now := clk.Now()
	return &TokenBucket{
		capacity:   float64(capacity),
		refillRate: refillRate,
		tokens:     float64(capacity),
		lastRefill: now,
		clock:      clk,
	}
}

// refill updates the token count based on elapsed time.
func (tb *TokenBucket) refill() {
	now := tb.clock.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	if elapsed <= 0 {
		return
	}
	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastRefill = now
}

// Allow attempts to consume n tokens. Returns true on success.
func (tb *TokenBucket) Allow(n int) bool {
	if n <= 0 {
		return true
	}
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	if tb.tokens >= float64(n) {
		tb.tokens -= float64(n)
		return true
	}
	return false
}

// GetTokens returns the current token count (rounded to 6 decimal places for testing).
func (tb *TokenBucket) GetTokens() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	return tb.tokens
}
