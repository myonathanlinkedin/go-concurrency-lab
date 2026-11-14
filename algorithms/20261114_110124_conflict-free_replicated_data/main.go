package main

import (
	"fmt"
)

func assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: " + msg)
	}
}

func main() {
	// Basic operations on two replicas.
	a := NewPNCounter("A")
	b := NewPNCounter("B")

	a.Increment(5)
	a.Decrement(2) // net +3
	b.Increment(3)
	b.Decrement(1) // net +2

	assert(a.Value() == 3, fmt.Sprintf("A value expected 3 got %d", a.Value()))
	assert(b.Value() == 2, fmt.Sprintf("B value expected 2 got %d", b.Value()))

	// Merge A into B; total should be 5.
	b.Merge(a)
	assert(b.Value() == 5, fmt.Sprintf("B after merge expected 5 got %d", b.Value()))

	// Merge back into A; both should see 5.
	a.Merge(b)
	assert(a.Value() == 5, fmt.Sprintf("A after merge expected 5 got %d", a.Value()))

	// Idempotence: merging again does not change value.
	prev := a.Value()
	a.Merge(b)
	assert(a.Value() == prev, "idempotent merge altered value")

	// Simulate many increments on a third replica.
	c := NewPNCounter("C")
	for i := 0; i < 1000; i++ {
		c.Increment(1)
	}
	assert(c.Value() == 1000, fmt.Sprintf("C expected 1000 got %d", c.Value()))

	// Merge with an empty replica; value unchanged.
	empty := NewPNCounter("D")
	c.Merge(empty)
	assert(c.Value() == 1000, "merge with empty changed value")

	// Verify Clone isolates state.
	cloneA := a.Clone()
	cloneB := b.Clone()
	cloneA.Increment(10)
	assert(a.Value() == 5, "original A changed after clone increment")
	assert(cloneA.Value() == 15, fmt.Sprintf("cloneA expected 15 got %d", cloneA.Value()))

	// Merging cloneA into B should raise B's value.
	b.Merge(cloneA)
	assert(b.Value() == 15, fmt.Sprintf("B after merging cloneA expected 15 got %d", b.Value()))

	// Ensure original cloneB unchanged.
	assert(cloneB.Value() == 5, "cloneB changed unexpectedly")

	fmt.Println("All PN-Counter tests passed")
}
