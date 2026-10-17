package main

import (
	"time"
)

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func testSlidingWindow() {
	limiter := NewSlidingWindowLimiter(5, 1*time.Second)
	for i := 0; i < 5; i++ {
		assert(limiter.Allow(), "should allow within capacity")
	}
	assert(!limiter.Allow(), "should reject after capacity")
	time.Sleep(1100 * time.Millisecond)
	assert(limiter.Allow(), "should allow after window")
}

func testTokenBucket() {
	limiter := NewTokenBucketLimiter(3, 1.0)
	assert(limiter.Allow(), "token 1")
	assert(limiter.Allow(), "token 2")
	assert(limiter.Allow(), "token 3")
	assert(!limiter.Allow(), "no tokens left")
	time.Sleep(1500 * time.Millisecond)
	assert(limiter.Allow(), "refilled token")
	assert(limiter.Allow(), "token")
	assert(limiter.Allow(), "token")
	assert(!limiter.Allow(), "none left")
}

func testDistributedLimiter() {
	limiter := NewDistributedLimiter(5, 2*time.Second, 2.0)
	for i := 0; i < 5; i++ {
		assert(limiter.Allow(), "distributed allow")
	}
	assert(!limiter.Allow(), "distributed reject after capacity")
	time.Sleep(2100 * time.Millisecond)
	assert(limiter.Allow(), "distributed allow after window")
	time.Sleep(1000 * time.Millisecond)
	assert(limiter.Allow(), "distributed allow after refill")
	assert(limiter.Allow(), "distributed allow after refill")
	assert(!limiter.Allow(), "distributed reject at capacity")
}

func main() {
	testSlidingWindow()
	testTokenBucket()
	testDistributedLimiter()
}
