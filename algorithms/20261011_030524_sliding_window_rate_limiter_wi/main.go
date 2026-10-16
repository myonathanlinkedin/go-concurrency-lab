package main

import (
	"fmt"
	"time"
)

func assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: " + msg)
	}
}

func testTokenBucketBasic() {
	rl := NewRateLimiter(5, 1, 0) // capacity 5, 1 token/sec, no sliding window
	// Consume initial capacity
	for i := 0; i < 5; i++ {
		assert(rl.Allow(time.Now()), fmt.Sprintf("initial allow %d should succeed", i+1))
	}
	// Sixth request should fail
	assert(!rl.Allow(time.Now()), "sixth request should be rejected")

	// Advance time by 2 seconds
	future := time.Now().Add(2 * time.Second)
	assert(rl.Allow(future), "request after 2s should succeed (1st token)")
	assert(rl.Allow(future), "request after 2s should succeed (2nd token)")
	assert(!rl.Allow(future), "third request after 2s should fail")
}

func testSlidingWindow() {
	window := 1 * time.Second
	rl := NewRateLimiter(3, 0, window) // capacity 3, no refill, sliding window enabled

	base := time.Unix(0, 0) // deterministic start
	// First three requests at t=0
	for i := 0; i < 3; i++ {
		assert(rl.Allow(base), fmt.Sprintf("sliding allow %d at t=0", i+1))
	}
	// Fourth request at t=0.5s should be rejected
	assert(!rl.Allow(base.Add(500*time.Millisecond)), "fourth request within window should be rejected")
	// After 1.1s, the first request falls out of the window
	assert(rl.Allow(base.Add(1100*time.Millisecond)), "request after window expiry should succeed")
	// Now we have 3 timestamps again; another immediate request should fail
	assert(!rl.Allow(base.Add(1100*time.Millisecond)), "exceeding capacity after refill should fail")
}

func testMerge() {
	window := 2 * time.Second
	rlA := NewRateLimiter(4, 0, window)
	rlB := NewRateLimiter(4, 0, window)

	base := time.Unix(0, 0)

	// rlA records two requests at t=0 and t=0.5s
	assert(rlA.Allow(base), "rlA first request")
	assert(rlA.Allow(base.Add(500*time.Millisecond)), "rlA second request")

	// rlB records one request at t=1s
	assert(rlB.Allow(base.Add(1*time.Second)), "rlB first request")

	// Merge B into A
	rlA.Merge(rlB)

	// After merge, timestamps should be three distinct entries
	assert(rlA.TimestampCount() == 3, fmt.Sprintf("expected 3 timestamps after merge, got %d", rlA.TimestampCount()))

	// A request at t=1.5s should be rejected because window still holds 3 tokens (capacity 4, but only 3 used)
	assert(rlA.Allow(base.Add(1500*time.Millisecond)), "request after merge should succeed (4th token)")

	// Now capacity is reached; next request should fail
	assert(!rlA.Allow(base.Add(1500*time.Millisecond)), "subsequent request should fail")
}

func main() {
	fmt.Println("Running Token Bucket Basic test...")
	testTokenBucketBasic()
	fmt.Println("Token Bucket Basic test passed.")

	fmt.Println("Running Sliding Window test...")
	testSlidingWindow()
	fmt.Println("Sliding Window test passed.")

	fmt.Println("Running Merge test...")
	testMerge()
	fmt.Println("Merge test passed.")

	fmt.Println("All tests passed.")
}
