package main

import (
	"sync"
	"time"
)

// RateLimiter implements a token bucket with an optional sliding‑window
// enforcement. It is safe for concurrent use.
type RateLimiter struct {
	capacity   uint64        // maximum number of tokens
	tokens     uint64        // current token count
	refillRate uint64        // tokens added per second
	lastRefill int64         // nanoseconds since epoch of last refill
	windowSize time.Duration // size of sliding window (0 disables sliding window)
	mu         sync.Mutex
	// timestamps of successful Allow calls; used only when windowSize > 0
	timestamps []int64 // nanoseconds since epoch, ordered by insertion
}

// NewRateLimiter creates a new RateLimiter.
//   capacity   – maximum tokens that can be stored.
//   refillRate – tokens added each second (integer granularity).
//   windowSize – if >0, sliding‑window enforcement is active; otherwise only token bucket logic applies.
func NewRateLimiter(capacity, refillRate uint64, windowSize time.Duration) *RateLimiter {
	now := time.Now().UnixNano()
	return &RateLimiter{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: now,
		windowSize: windowSize,
	}
}

// refill adds tokens according to elapsed time since the last refill.
func (rl *RateLimiter) refill(now int64) {
	if rl.refillRate == 0 {
		rl.lastRefill = now
		return
	}
	elapsedSec := uint64((now - rl.lastRefill) / int64(time.Second))
	if elapsedSec == 0 {
		return
	}
	added := elapsedSec * rl.refillRate
	if added > 0 {
		if rl.tokens+added > rl.capacity {
			rl.tokens = rl.capacity
		} else {
			rl.tokens += added
		}
		rl.lastRefill += int64(elapsedSec) * int64(time.Second)
	}
}

// pruneOld removes timestamps that are older than the sliding window.
func (rl *RateLimiter) pruneOld(now int64) {
	if rl.windowSize == 0 {
		return
	}
	cutoff := now - rl.windowSize.Nanoseconds()
	idx := 0
	for _, ts := range rl.timestamps {
		if ts > cutoff {
			break
		}
		idx++
	}
	if idx > 0 {
		rl.timestamps = rl.timestamps[idx:]
	}
}

// Allow attempts to consume a single token at the supplied time.
// It returns true if the request is permitted, false otherwise.
// The method respects both token‑bucket limits and, when windowSize > 0,
// the sliding‑window request count.
func (rl *RateLimiter) Allow(at time.Time) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := at.UnixNano()
	rl.refill(now)
	rl.pruneOld(now)

	// Sliding‑window check (if enabled)
	if rl.windowSize > 0 {
		if uint64(len(rl.timestamps)) >= rl.capacity {
			return false
		}
		// Record the timestamp for the successful request
		rl.timestamps = append(rl.timestamps, now)
		rl.tokens-- // also consume a token for consistency
		return true
	}

	// Pure token‑bucket check
	if rl.tokens == 0 {
		return false
	}
	rl.tokens--
	return true
}

// Merge combines the state of another RateLimiter into this one.
// It is intended for a distributed scenario where each node reports its
// local limiter state. The merge respects the smallest capacity among the
// two limiters and unifies timestamps, then prunes according to the window.
func (rl *RateLimiter) Merge(other *RateLimiter) {
	// Ensure both limiters use the same configuration; otherwise ignore merge.
	if rl.capacity != other.capacity || rl.refillRate != other.refillRate || rl.windowSize != other.windowSize {
		return
	}
	rl.mu.Lock()
	other.mu.Lock()
	defer rl.mu.Unlock()
	defer other.mu.Unlock()

	// Merge token counts (take the minimum to avoid over‑granting)
	if other.tokens < rl.tokens {
		rl.tokens = other.tokens
	}
	// Merge timestamps (union)
	merged := make([]int64, 0, len(rl.timestamps)+len(other.timestamps))
	i, j := 0, 0
	for i < len(rl.timestamps) && j < len(other.timestamps) {
		if rl.timestamps[i] < other.timestamps[j] {
			merged = append(merged, rl.timestamps[i])
			i++
		} else if rl.timestamps[i] > other.timestamps[j] {
			merged = append(merged, other.timestamps[j])
			j++
		} else {
			merged = append(merged, rl.timestamps[i])
			i++
			j++
		}
	}
	for ; i < len(rl.timestamps); i++ {
		merged = append(merged, rl.timestamps[i])
	}
	for ; j < len(other.timestamps); j++ {
		merged = append(merged, other.timestamps[j])
	}
	rl.timestamps = merged

	// Update lastRefill to the most recent of the two
	if other.lastRefill > rl.lastRefill {
		rl.lastRefill = other.lastRefill
	}
	// Prune any timestamps that are now outside the window
	rl.pruneOld(time.Now().UnixNano())
}

// Tokens returns the current token count (for testing/observability).
func (rl *RateLimiter) Tokens() uint64 {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.refill(time.Now().UnixNano())
	return rl.tokens
}

// TimestampCount returns the number of timestamps stored (for testing).
func (rl *RateLimiter) TimestampCount() int {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	return len(rl.timestamps)
}
