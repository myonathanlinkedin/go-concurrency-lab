package main

import (
	"sync"
	"time"
)

// Limiter defines the common interface for rate limiting algorithms.
type Limiter interface {
	// Allow attempts to consume a single token/request.
	// Returns true if the request is allowed, false otherwise.
	Allow() bool
}

// Clock abstracts time source for deterministic testing.
type Clock interface {
	Now() time.Time
}

// RealClock implements Clock using the actual system time.
type RealClock struct{}

// Now returns the current local time.
func (rc *RealClock) Now() time.Time { return time.Now() }

// MockClock implements Clock with manually controllable time.
type MockClock struct {
	mu      sync.Mutex
	current time.Time
}

// NewMockClock creates a MockClock initialized at the provided start time.
func NewMockClock(start time.Time) *MockClock {
	return &MockClock{current: start}
}

// Now returns the mocked current time.
func (mc *MockClock) Now() time.Time {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	return mc.current
}

// Advance moves the mocked time forward by d.
func (mc *MockClock) Advance(d time.Duration) {
	mc.mu.Lock()
	mc.current = mc.current.Add(d)
	mc.mu.Unlock()
}

// TokenBucketLimiter implements a classic token bucket algorithm.
type TokenBucketLimiter struct {
	capacity    uint64        // maximum number of tokens
	refillRate  float64       // tokens added per second
	tokens      float64       // current token count (may be fractional)
	lastRefill  time.Time     // last time tokens were refilled
	clock       Clock         // time source
	mu          sync.Mutex    // protects fields above
}

// NewTokenBucketLimiter constructs a token bucket limiter.
// capacity must be >0, refillRate must be >=0.
func NewTokenBucketLimiter(capacity uint64, refillRate float64, clk Clock) *TokenBucketLimiter {
	if capacity == 0 {
		panic("capacity must be greater than zero")
	}
	if refillRate < 0 {
		panic("refillRate cannot be negative")
	}
	if clk == nil {
		clk = &RealClock{}
	}
	return &TokenBucketLimiter{
		capacity:   capacity,
		refillRate: refillRate,
		tokens:     float64(capacity),
		lastRefill: clk.Now(),
		clock:      clk,
	}
}

// refill updates the token count based on elapsed time.
func (tb *TokenBucketLimiter) refill() {
	now := tb.clock.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	if elapsed <= 0 {
		return
	}
	added := elapsed * tb.refillRate
	tb.tokens += added
	if tb.tokens > float64(tb.capacity) {
		tb.tokens = float64(tb.capacity)
	}
	tb.lastRefill = now
}

// Allow consumes a single token if available.
func (tb *TokenBucketLimiter) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	if tb.tokens >= 1 {
		tb.tokens -= 1
		return true
	}
	return false
}

// SlidingWindowLimiter implements a fixed-size sliding window counter.
// It divides the total window into equal sub-windows (buckets) and
// tracks request counts per bucket.
type SlidingWindowLimiter struct {
	windowSize      time.Duration // total window length (e.g., 1 minute)
	subWindowCount  int           // number of buckets (e.g., 10)
	maxRequests     int           // allowed requests per window
	buckets         []int         // request count per bucket
	bucketStart     time.Time     // start time of the first bucket
	bucketDuration  time.Duration // duration of each sub-window
	clock           Clock
	mu              sync.Mutex
}

// NewSlidingWindowLimiter constructs a sliding window limiter.
// windowSize must be >0, subWindowCount >0, maxRequests >=0.
func NewSlidingWindowLimiter(windowSize time.Duration, subWindowCount int, maxRequests int, clk Clock) *SlidingWindowLimiter {
	if windowSize <= 0 {
		panic("windowSize must be positive")
	}
	if subWindowCount <= 0 {
		panic("subWindowCount must be positive")
	}
	if maxRequests < 0 {
		panic("maxRequests cannot be negative")
	}
	if clk == nil {
		clk = &RealClock{}
	}
	buckets := make([]int, subWindowCount)
	return &SlidingWindowLimiter{
		windowSize:     windowSize,
		subWindowCount: subWindowCount,
		maxRequests:    maxRequests,
		buckets:        buckets,
		bucketStart:    clk.Now(),
		bucketDuration: windowSize / time.Duration(subWindowCount),
		clock:          clk,
	}
}

// advanceWindow moves the sliding window forward, clearing stale buckets.
func (sw *SlidingWindowLimiter) advanceWindow() {
	now := sw.clock.Now()
	elapsed := now.Sub(sw.bucketStart)
	if elapsed < sw.bucketDuration {
		return
	}
	steps := int(elapsed / sw.bucketDuration)
	if steps >= sw.subWindowCount {
		// All buckets are stale.
		for i := range sw.buckets {
			sw.buckets[i] = 0
		}
		sw.bucketStart = now
		return
	}
	// Shift buckets left by steps, zeroing the vacated slots.
	copy(sw.buckets, sw.buckets[steps:])
	for i := sw.subWindowCount - steps; i < sw.subWindowCount; i++ {
		sw.buckets[i] = 0
	}
	sw.bucketStart = sw.bucketStart.Add(time.Duration(steps) * sw.bucketDuration)
}

// totalCount returns the sum of all bucket counts.
func (sw *SlidingWindowLimiter) totalCount() int {
	sum := 0
	for _, v := range sw.buckets {
		sum += v
	}
	return sum
}

// Allow registers a request and returns true if within limit.
func (sw *SlidingWindowLimiter) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.advanceWindow()
	if sw.totalCount() < sw.maxRequests {
		// Increment the current bucket (the last one).
		sw.buckets[sw.subWindowCount-1]++
		return true
	}
	return false
}
