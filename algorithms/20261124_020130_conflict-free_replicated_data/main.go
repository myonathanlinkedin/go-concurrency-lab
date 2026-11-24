package main

import (
	"fmt"
)

func assert(cond bool, msg string) {
	if !cond {
		panic(fmt.Sprintf("assertion failed: %s", msg))
	}
}

func main() {
	// Test 1: Basic increment / decrement behavior.
	c1 := NewPNCounter()
	c1.Increment("A", 5)
	c1.Decrement("A", 2)
	assert(c1.Value() == 3, "basic inc/dec value should be 3")

	// Test 2: Merge two counters with disjoint replica sets.
	c2 := NewPNCounter()
	c2.Increment("B", 7)
	c2.Decrement("B", 1)

	c1.Merge(c2)
	assert(c1.Value() == 9, "after merge with B, value should be 9 (3 + 6)")

	// Test 3: Merge with overlapping replica identifiers – max rule.
	c3 := NewPNCounter()
	c3.Increment("A", 10) // larger than previous A contribution (5)
	c3.Decrement("B", 4) // larger than previous B contribution (1)

	c1.Merge(c3)
	// Expected: inc[A]=10, inc[B]=7, dec[A]=2, dec[B]=4
	expected := int64(10 + 7 - 2 - 4) // =11
	assert(c1.Value() == expected, "merge with overlapping ids should respect max")

	// Test 4: Idempotent merge (merging same state twice does not change value).
	before := c1.Value()
	c1.Merge(c3) // merge again
	assert(c1.Value() == before, "idempotent merge should not alter value")

	// Test 5: Commutativity – A.Merge(B) == B.Merge(A)
	x := NewPNCounter()
	x.Increment("X", 3)
	y := NewPNCounter()
	y.Increment("Y", 4)

	// Clone for independent merges
	x1 := x.Clone()
	y1 := y.Clone()

	x1.Merge(y1)
	y1.Merge(x)

	assert(x1.Value() == y1.Value(), "commutative merge should yield same value")
	assert(x1.Value() == 7, "commutative merge result should be 7")

	// Test 6: Concurrency safety – perform concurrent ops.
	c4 := NewPNCounter()
	done := make(chan struct{})
	go func() {
		for i := 0; i < 1000; i++ {
			c4.Increment("C", 1)
		}
		done <- struct{}{}
	}()
	go func() {
		for i := 0; i < 600; i++ {
			c4.Decrement("C", 1)
		}
		done <- struct{}{}
	}()
	<-done
	<-done
	assert(c4.Value() == 400, "concurrent ops should result in 400")

	// All tests passed.
	fmt.Println("All PN‑Counter tests passed.")
}
