package main

import (
	"fmt"
	"sync"
	"time"
)

type mockClock struct {
	now time.Time
}

func (mc *mockClock) Now() time.Time { return mc.now }

func (mc *mockClock) Advance(d time.Duration) { mc.now = mc.now.Add(d) }

func assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: " + msg)
	}
}

func approxEqual(a, b, eps float64) bool {
	if a > b {
		return a-b < eps
	}
	return b-a < eps
}

func main() {
	// Basic functionality test
	mc := &mockClock{now: time.Unix(0, 0)}
	limiter := NewTokenBucket(10, 5, mc) // capacity 10, refill 5 tokens/sec

	assert(limiter.GetTokens() == 10, "initial tokens")
	assert(limiter.Allow(3), "allow 3")
	assert(limiter.GetTokens() == 7, "tokens after 3")
	assert(!limiter.Allow(8), "reject 8")
	assert(limiter.GetTokens() == 7, "tokens unchanged after reject")

	mc.Advance(time.Second) // +5 tokens, capped at 10
	assert(limiter.GetTokens() == 10, "refill to capacity")
	assert(limiter.Allow(10), "allow full")
	assert(limiter.GetTokens() == 0, "empty after full consume")

	mc.Advance(500 * time.Millisecond) // +2.5 tokens
	toks := limiter.GetTokens()
	assert(approxEqual(toks, 2.5, 1e-6), fmt.Sprintf("half second refill got %f", toks))

	// Concurrency test
	mc2 := &mockClock{now: time.Unix(0, 0)}
	concLimiter := NewTokenBucket(100, 20, mc2)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for i := 0; i < 50; i++ {
			_ = concLimiter.Allow(1)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < 50; i++ {
			_ = concLimiter.Allow(1)
		}
		wg.Done()
	}()
	wg.Wait()
	remaining := concLimiter.GetTokens()
	assert(remaining >= 0 && remaining <= 100, fmt.Sprintf("concurrency remaining %f", remaining))

	fmt.Println("All assertions passed")
}
