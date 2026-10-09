package main

import (
	"sync"
	"time"
)

type RateLimiter interface {
	Allow() bool
	AllowN(n int) bool
}

type SlidingWindowLimiter struct {
	mu         sync.Mutex
	timestamps []time.Time
	capacity   int
	window     time.Duration
}

func NewSlidingWindowLimiter(capacity int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		capacity: capacity,
		window:   window,
	}
}

func (l *SlidingWindowLimiter) prune(now time.Time) {
	cutoff := now.Add(-l.window)
	i := 0
	for _, t := range l.timestamps {
		if t.After(cutoff) {
			break
		}
		i++
	}
	l.timestamps = l.timestamps[i:]
}

func (l *SlidingWindowLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	l.prune(now)
	if len(l.timestamps) < l.capacity {
		l.timestamps = append(l.timestamps, now)
		return true
	}
	return false
}

func (l *SlidingWindowLimiter) AllowN(n int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	l.prune(now)
	if len(l.timestamps)+n <= l.capacity {
		for i := 0; i < n; i++ {
			l.timestamps = append(l.timestamps, now)
		}
		return true
	}
	return false
}

type TokenBucketLimiter struct {
	mu         sync.Mutex
	capacity   int
	tokens     int
	fillRate   float64 // tokens per second
	lastRefill time.Time
}

func NewTokenBucketLimiter(capacity int, fillRate float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity:   capacity,
		tokens:     capacity,
		fillRate:   fillRate,
		lastRefill: time.Now(),
	}
}

func (l *TokenBucketLimiter) refill(now time.Time) {
	elapsed := now.Sub(l.lastRefill).Seconds()
	added := int(elapsed * l.fillRate)
	if added > 0 {
		l.tokens += added
		if l.tokens > l.capacity {
			l.tokens = l.capacity
		}
		l.lastRefill = now
	}
}

func (l *TokenBucketLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	l.refill(now)
	if l.tokens > 0 {
		l.tokens--
		return true
	}
	return false
}

func (l *TokenBucketLimiter) AllowN(n int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	l.refill(now)
	if l.tokens >= n {
		l.tokens -= n
		return true
	}
	return false
}

type DistributedLimiter struct {
	tb *TokenBucketLimiter
	sw *SlidingWindowLimiter
}

func NewDistributedLimiter(capacity int, window time.Duration, fillRate float64) *DistributedLimiter {
	return &DistributedLimiter{
		tb: NewTokenBucketLimiter(capacity, fillRate),
		sw: NewSlidingWindowLimiter(capacity, window),
	}
}

func (l *DistributedLimiter) Allow() bool {
	if !l.tb.Allow() {
		return false
	}
	return l.sw.Allow()
}

func (l *DistributedLimiter) AllowN(n int) bool {
	if !l.tb.AllowN(n) {
		return false
	}
	return l.sw.AllowN(n)
}
