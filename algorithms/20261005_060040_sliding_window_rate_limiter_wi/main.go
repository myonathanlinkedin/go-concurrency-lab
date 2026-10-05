package main

import (
	"fmt"
	"sync"
	"time"
)

// Clock abstracts time source for production and testing.
type Clock interface {
	Now() time.Time
}

// realClock uses the actual system time.
type realClock struct{}

func (rc *realClock) Now() time.Time { return time.Now() }

// mockClock allows manual time manipulation in tests.
type mockClock struct {
	mu  sync.Mutex
	now time.Time
}

func newMockClock(start time.Time) *mockClock {
	return &mockClock{now: start}
}

func (mc *mockClock) Now() time.Time {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	return mc.now
}

func (mc *mockClock) Add(d time.Duration) {
	mc.mu.Lock()
	mc.now = mc.now.Add(d)
	mc.mu.Unlock()
}

// RateLimiter implements a token bucket with sliding window semantics.
type RateLimiter struct {
	capacity   float64       // maximum tokens
	refillRate float64       // tokens per second
	clock      Clock         // time source
	mu         sync.Mutex    // protects state
	tokens     float64       // current token count
	lastRefill time.Time     // last refill timestamp
}

// NewRateLimiter creates a new limiter.
func NewRateLimiter(capacity int, refillPerSec float64, clk Clock) *RateLimiter {
	if capacity <= 0 {
		panic("capacity must be positive")
	}
	if refillPerSec <= 0 {
		panic("refill rate must be positive")
	}
	if clk == nil {
		clk = &realClock{}
	}
	now := clk.Now()
	return &RateLimiter{
		capacity:   float64(capacity),
		refillRate: refillPerSec,
		clock:      clk,
		tokens:     float64(capacity),
		lastRefill: now,
	}
}

// refill updates token count based on elapsed time.
func (rl *RateLimiter) refill() {
	now := rl.clock.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	if elapsed <= 0 {
		return
	}
	rl.tokens += elapsed * rl.refillRate
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}
	rl.lastRefill = now
}

// Allow attempts to consume n tokens. Returns true if allowed.
func (rl *RateLimiter) Allow(n int) bool {
	if n <= 0 {
		return false
	}
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.refill()
	if rl.tokens >= float64(n) {
		rl.tokens -= float64(n)
		return true
	}
	return false
}

// ---------- Unit Tests ----------
type testCase struct {
	name string
	run  func() error
}

func assert(cond bool, msg string) error {
	if !cond {
		return fmt.Errorf("assertion failed: %s", msg)
	}
	return nil
}

func testBasicAllowance() error {
	clk := newMockClock(time.Unix(0, 0))
	rl := NewRateLimiter(5, 1, clk) // 5 tokens capacity, 1 token/sec refill
	if err := assert(rl.Allow(3), "first allow 3"); err != nil {
		return err
	}
	if err := assert(!rl.Allow(3), "second allow 3 should fail"); err != nil {
		return err
	}
	clk.Add(2 * time.Second) // refill 2 tokens
	if err := assert(rl.Allow(2), "allow after refill 2"); err != nil {
		return err
	}
	if err := assert(!rl.Allow(1), "no tokens left"); err != nil {
		return err
	}
	return nil
}

func testRefillPrecision() error {
	clk := newMockClock(time.Unix(0, 0))
	rl := NewRateLimiter(10, 2.5, clk) // 2.5 tokens/sec
	if err := assert(rl.Allow(10), "consume all"); err != nil {
		return err
	}
	clk.Add(400 * time.Millisecond) // 0.4s * 2.5 = 1 token
	if err := assert(rl.Allow(1), "allow 1 after 0.4s"); err != nil {
		return err
	}
	if err := assert(!rl.Allow(1), "second token not yet available"); err != nil {
		return err
	}
	clk.Add(200 * time.Millisecond) // total 0.6s => 1.5 tokens, 0.5 leftover
	if err := assert(rl.Allow(1), "allow after additional 0.2s"); err != nil {
		return err
	}
	return nil
}

func testConcurrentAccess() error {
	clk := newMockClock(time.Unix(0, 0))
	rl := NewRateLimiter(100, 50, clk) // high capacity for test
	var wg sync.WaitGroup
	errCh := make(chan error, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !rl.Allow(5) {
				errCh <- fmt.Errorf("concurrent allow failed")
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for e := range errCh {
		return e
	}
	return nil
}

func runTests() error {
	tests := []testCase{
		{"BasicAllowance", testBasicAllowance},
		{"RefillPrecision", testRefillPrecision},
		{"ConcurrentAccess", testConcurrentAccess},
	}
	for _, tc := range tests {
		if err := tc.run(); err != nil {
			return fmt.Errorf("%s: %w", tc.name, err)
		}
	}
	return nil
}

// ---------- Main ----------
func main() {
	if err := runTests(); err != nil {
		panic(err)
	}
	fmt.Println("All tests passed")
}
