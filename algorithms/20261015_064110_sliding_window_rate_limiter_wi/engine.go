package main

import (
	"fmt"
	"time"
)

// Engine provides higher‑level utilities that combine multiple limiters.
type Engine struct {
	limiters []Limiter
}

// NewEngine creates an Engine with the supplied limiters.
func NewEngine(limiters ...Limiter) *Engine {
	return &Engine{limiters: limiters}
}

// Allow checks all configured limiters; the request is allowed only if
// every limiter permits it.
func (e *Engine) Allow() bool {
	for _, l := range e.limiters {
		if !l.Allow() {
			return false
		}
	}
	return true
}

// DemoEngine runs a short demonstration of combined limiters.
func DemoEngine() {
	clk := &RealClock{}
	tb := NewTokenBucketLimiter(5, 1.0, clk)               // 5 tokens, 1 token/sec
	sw := NewSlidingWindowLimiter(10*time.Second, 5, 8, clk) // 8 req/10s, 5 buckets
	eng := NewEngine(tb, sw)

	fmt.Println("Demo: combined limiter (TokenBucket + SlidingWindow)")
	for i := 1; i <= 12; i++ {
		allowed := eng.Allow()
		fmt.Printf("Request %2d: allowed=%v\n", i, allowed)
		time.Sleep(500 * time.Millisecond)
	}
}
