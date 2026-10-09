package main

import (
	"fmt"
	"time"
)

func main() {
	// Configuration: 5 tokens capacity, 2 tokens per second refill,
	// sliding window of 3 seconds allowing up to 3 requests.
	cfg := Config{
		Capacity:    5,
		RefillRate:  2.0,
		WindowSize:  3 * time.Second,
		MaxRequests: 3,
	}
	rl := NewRateLimiter(cfg)

	// Helper to assert conditions.
	assert := func(cond bool, msg string) {
		if !cond {
			panic(fmt.Sprintf("assertion failed: %s", msg))
		}
	}

	// 1. Initial burst: should allow up to Capacity (5) requests instantly,
	// but sliding window caps at MaxRequests (3) within 3 s.
	for i := 0; i < 3; i++ {
		assert(rl.Allow(), fmt.Sprintf("initial request %d should be allowed", i+1))
	}
	assert(!rl.Allow(), "fourth request should be blocked by sliding window")

	// 2. Wait for sliding window to expire (3 s) and for tokens to refill.
	time.Sleep(3 * time.Second)

	// After waiting, the sliding window is empty and tokens have refilled:
	// 3 seconds elapsed → 6 tokens added, but capacity caps at 5.
	// So we should have full capacity again.
	assert(rl.Allow(), "request after window expiry should be allowed")
	assert(rl.Allow(), "second request after window expiry should be allowed")
	assert(rl.Allow(), "third request after window expiry should be allowed")
	assert(!rl.Allow(), "fourth request after window expiry should be blocked by sliding window")

	// 3. Test token depletion without sliding‑window pressure.
	// Create a limiter without a sliding window.
	cfgNoWindow := Config{
		Capacity:    3,
		RefillRate:  1.0,
		WindowSize:  0,
		MaxRequests: 0,
	}
	rlNoWindow := NewRateLimiter(cfgNoWindow)

	// Consume all tokens.
	for i := 0; i < 3; i++ {
		assert(rlNoWindow.Allow(), fmt.Sprintf("token bucket request %d should be allowed", i+1))
	}
	assert(!rlNoWindow.Allow(), "fourth request should be blocked (no tokens)")

	// Wait half a second → 0.5 token should be added (fractional tokens allowed internally).
	time.Sleep(500 * time.Millisecond)
	assert(!rlNoWindow.Allow(), "still blocked because only 0.5 token accumulated")

	// Wait another 600 ms → total >1 token accumulated.
	time.Sleep(600 * time.Millisecond)
	assert(rlNoWindow.Allow(), "request should be allowed after enough tokens refilled")

	// 4. Distributed merge test.
	// Two limiters with identical config.
	rlA := NewRateLimiter(cfg)
	rlB := NewRateLimiter(cfg)

	// Consume 2 tokens on A, 1 token on B.
	assert(rlA.Allow(), "A first request")
	assert(rlA.Allow(), "A second request")
	assert(rlB.Allow(), "B first request")

	// Merge B into A.
	rlA.Merge(rlB)

	// After merge, token count should be the minimum of the two.
	// A had 5‑2=3 tokens left, B had 5‑1=4 tokens left → A should now have 3 tokens.
	// Verify by consuming remaining tokens.
	for i := 0; i < 3; i++ {
		assert(rlA.Allow(), fmt.Sprintf("post‑merge request %d should be allowed", i+1))
	}
	assert(!rlA.Allow(), "post‑merge fourth request should be blocked (tokens exhausted)")

	// 5. DebugString sanity check (should not panic).
	_ = rlA.DebugString()

	fmt.Println("All assertions passed.")
}
