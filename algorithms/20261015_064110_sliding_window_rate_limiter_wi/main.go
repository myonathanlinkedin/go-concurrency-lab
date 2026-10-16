package main

import (
	"fmt"
	"time"
)

func main() {
	// ---------- TokenBucketLimiter deterministic tests ----------
	start := time.Unix(0, 0)
	mock := NewMockClock(start)

	tb := NewTokenBucketLimiter(3, 2.0, mock) // capacity 3, 2 tokens/sec

	// Initially, bucket is full.
	if !tb.Allow() {
		panic("expected first request to be allowed")
	}
	if !tb.Allow() {
		panic("expected second request to be allowed")
	}
	if !tb.Allow() {
		panic("expected third request to be allowed")
	}
	// Bucket exhausted.
	if tb.Allow() {
		panic("expected fourth request to be rejected")
	}

	// Advance 0.5 seconds: should add 1 token (2 * 0.5)
	mock.Advance(500 * time.Millisecond)
	if !tb.Allow() {
		panic("expected request after 0.5s refill to be allowed")
	}
	// Exhaust again.
	if tb.Allow() {
		panic("expected request after immediate refill to be rejected")
	}

	// Advance 1 second: should refill to capacity (adds 2 tokens, but max 3)
	mock.Advance(1 * time.Second)
	if !tb.Allow() || !tb.Allow() || !tb.Allow() {
		panic("expected three requests after 1s refill to be allowed")
	}
	if tb.Allow() {
		panic("expected request beyond capacity after refill to be rejected")
	}

	// ---------- SlidingWindowLimiter deterministic tests ----------
	sw := NewSlidingWindowLimiter(4*time.Second, 4, 3, mock) // 3 req per 4s, 4 buckets of 1s each

	// All three requests should pass.
	if !sw.Allow() || !sw.Allow() || !sw.Allow() {
		panic("expected first three requests to be allowed")
	}
	// Fourth request within same window should be rejected.
	if sw.Allow() {
		panic("expected fourth request within window to be rejected")
	}

	// Advance 1 second: oldest bucket expires, freeing one slot.
	mock.Advance(1 * time.Second)
	if !sw.Allow() {
		panic("expected request after sliding window advanced to be allowed")
	}
	// Now at limit again.
	if sw.Allow() {
		panic("expected request exceeding limit after advance to be rejected")
	}

	// Advance 3 seconds: all buckets expire, should reset.
	mock.Advance(3 * time.Second)
	if !sw.Allow() || !sw.Allow() || !sw.Allow() {
		panic("expected three requests after full window reset to be allowed")
	}
	if sw.Allow() {
		panic("expected request beyond limit after full reset to be rejected")
	}

	// ---------- Engine combined limiter test ----------
	// Use fresh mock clock for clarity.
	mock2 := NewMockClock(start)
	tb2 := NewTokenBucketLimiter(2, 1.0, mock2)               // 2 tokens, 1 token/sec
	sw2 := NewSlidingWindowLimiter(5*time.Second, 5, 4, mock2) // 4 req/5s
	eng := NewEngine(tb2, sw2)

	// First two requests: token bucket allows, sliding window allows.
	if !eng.Allow() || !eng.Allow() {
		panic("engine: first two requests should be allowed")
	}
	// Third request: token bucket exhausted, should be rejected.
	if eng.Allow() {
		panic("engine: third request should be rejected due to token bucket")
	}
	// Advance time to refill one token.
	mock2.Advance(1 * time.Second)
	if !eng.Allow() {
		panic("engine: request after refill should be allowed")
	}
	// Now sliding window has 3 requests in 5s (limit 4), token bucket has 0 tokens.
	if eng.Allow() {
		panic("engine: request should be rejected due to token bucket")
	}
	// Advance time to refill token and expire oldest sliding window bucket.
	mock2.Advance(2 * time.Second)
	if !eng.Allow() {
		panic("engine: request after both limits refreshed should be allowed")
	}

	fmt.Println("All assertions passed.")
	// Optional demonstration (non‑deterministic, uses real time).
	DemoEngine()
}
