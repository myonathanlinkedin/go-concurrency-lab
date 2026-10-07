package main

import (
	"fmt"
)

func main() {
	// Create two replicas
	a := NewPNCounter("A")
	b := NewPNCounter("B")

	// Perform operations on replica A
	a.Increment(5)
	a.Decrement(2)

	// Perform operations on replica B
	b.Increment(3)
	b.Decrement(4)

	// Verify individual values before merge
	if a.Value() != 3 {
		panic(fmt.Sprintf("expected A value 3, got %d", a.Value()))
	}
	if b.Value() != -1 {
		panic(fmt.Sprintf("expected B value -1, got %d", b.Value()))
	}

	// Merge counters
	a.Merge(b)
	b.Merge(a)

	// Expected merged value: (5+3) - (2+4) = 2
	expected := int64(2)

	if a.Value() != expected {
		panic(fmt.Sprintf("expected merged A value %d, got %d", expected, a.Value()))
	}
	if b.Value() != expected {
		panic(fmt.Sprintf("expected merged B value %d, got %d", expected, b.Value()))
	}

	// Edge case: merging with an empty counter
	empty := NewPNCounter("C")
	a.Merge(empty)
	if a.Value() != expected {
		panic(fmt.Sprintf("merge with empty altered value, expected %d, got %d", expected, a.Value()))
	}

	// Edge case: large increments to test uint64 handling
	large := NewPNCounter("D")
	const big uint64 = uint64(0x9e3779b97f4a7c15)
	large.Increment(big)
	large.Decrement(1)
	large.Merge(a)

	// Compute expected after merge
	var sumP uint64 = 5 + 3 + big
	var sumN uint64 = 2 + 4 + 1
	expectedLarge := int64(sumP) - int64(sumN)

	if large.Value() != expectedLarge {
		panic(fmt.Sprintf("expected large merged value %d, got %d", expectedLarge, large.Value()))
	}

	fmt.Println("All assertions passed.")
}
