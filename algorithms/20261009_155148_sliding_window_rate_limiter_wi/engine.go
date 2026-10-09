package main

import (
	"fmt"
	"time"
)

// NewRateLimiter creates a RateLimiter with the supplied configuration.
// The initial token count equals the capacity, and the request history is empty.
func NewRateLimiter(cfg Config) *RateLimiter {
	if cfg.Capacity == 0 {
		panic("capacity must be > 0")
	}
	if cfg.RefillRate <= 0 {
		panic("refill rate must be > 0")
	}
	if cfg.WindowSize < 0 {
		panic("window size must be >= 0")
	}
	if cfg.MaxRequests == 0 && cfg.WindowSize > 0 {
		panic("max requests must be > 0 when window size is positive")
	}
	rl := &RateLimiter{
		cfg:        cfg,
		tokens:     float64(cfg.Capacity),
		lastRefill: time.Now(),
		requests:   make([]time.Time, 0, cfg.MaxRequests),
	}
	return rl
}

// refill updates the token count based on elapsed time since the last refill.
func (rl *RateLimiter) refill(now time.Time) {
	elapsed := now.Sub(rl.lastRefill).Seconds()
	if elapsed <= 0 {
		return
	}
	added := elapsed * rl.cfg.RefillRate
	rl.tokens += added
	if rl.tokens > float64(rl.cfg.Capacity) {
		rl.tokens = float64(rl.cfg.Capacity)
	}
	rl.lastRefill = now
}

// pruneOldRequests discards timestamps that fall outside the sliding window.
func (rl *RateLimiter) pruneOldRequests(now time.Time) {
	if rl.cfg.WindowSize == 0 {
		// No sliding‑window constraint.
		rl.requests = rl.requests[:0]
		return
	}
	cutoff := now.Add(-rl.cfg.WindowSize)
	idx := 0
	for _, ts := range rl.requests {
		if ts.After(cutoff) {
			break
		}
		idx++
	}
	rl.requests = rl.requests[idx:]
}

// Allow determines whether a request may proceed.
// It returns true if both the token bucket and sliding‑window constraints
// are satisfied; otherwise false.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	rl.refill(now)
	rl.pruneOldRequests(now)

	// Token bucket check.
	if rl.tokens < 1.0 {
		return false
	}
	// Sliding window check.
	if rl.cfg.WindowSize > 0 && uint64(len(rl.requests)) >= rl.cfg.MaxRequests {
		return false
	}

	// Consume a token and record the request.
	rl.tokens -= 1.0
	rl.requests = append(rl.requests, now)
	return true
}

// Merge combines the state of another RateLimiter into this one.
// The merge is deterministic and safe for distributed deployments:
//   * Tokens are set to the minimum of the two token counts (to avoid exceeding capacity).
//   * Request histories are merged and then trimmed to respect the sliding window.
// The method assumes both limiters share the same configuration.
func (rl *RateLimiter) Merge(other *RateLimiter) {
	if rl == nil || other == nil {
		return
	}
	if rl.cfg != other.cfg {
		panic("cannot merge limiters with different configurations")
	}

	rl.mu.Lock()
	other.mu.Lock()
	defer rl.mu.Unlock()
	defer other.mu.Unlock()

	now := time.Now()
	rl.refill(now)
	other.refill(now)

	// Deterministic token merge: keep the smaller token count.
	if other.tokens < rl.tokens {
		rl.tokens = other.tokens
	}
	// Merge request timestamps.
	merged := make([]time.Time, 0, len(rl.requests)+len(other.requests))
	merged = append(merged, rl.requests...)
	merged = append(merged, other.requests...)
	// Sort merged slice (simple insertion sort because slice is small).
	for i := 1; i < len(merged); i++ {
		j := i
		for j > 0 && merged[j-1].After(merged[j]) {
			merged[j-1], merged[j] = merged[j], merged[j-1]
			j--
		}
	}
	rl.requests = merged
	rl.pruneOldRequests(now)
}

// DebugString returns a human‑readable snapshot of the limiter's state.
func (rl *RateLimiter) DebugString() string {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	rl.refill(now)
	rl.pruneOldRequests(now)
	return fmt.Sprintf("tokens=%.2f, requests=%d, window=%s",
		rl.tokens, len(rl.requests), rl.cfg.WindowSize)
}
