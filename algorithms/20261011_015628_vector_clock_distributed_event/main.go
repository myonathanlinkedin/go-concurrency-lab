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
	// Test 1: Simple increment and comparison (happens‑after).
	c1 := NewVectorClock(0)
	c1.Increment(0) // pid 0 -> 1
	c1.Increment(0) // pid 0 -> 2
	c1.Increment(1) // pid 1 -> 1

	c0 := NewVectorClock(0) // empty clock

	rel := c0.Compare(c1)
	assert(rel == -1, "c0 should happen‑before c1")
	rel = c1.Compare(c0)
	assert(rel == 1, "c1 should happen‑after c0")

	// Test 2: Concurrent clocks.
	a := NewVectorClock(0)
	b := NewVectorClock(0)

	a.Increment(0) // a: {0:1}
	b.Increment(1) // b: {1:1}

	rel = a.Compare(b)
	assert(rel == 2, "a and b should be concurrent")
	rel = b.Compare(a)
	assert(rel == 2, "b and a should be concurrent")

	merged := a.Merge(b)
	expected := "{0:1,1:1}"
	assert(merged.String() == expected, fmt.Sprintf("merged clock expected %s got %s", expected, merged.String()))

	// Test 3: Transitive ordering.
	c := a.Clone()
	c.Update(merged) // c now {0:1,1:1}
	c.Increment(2)   // c now {0:1,1:1,2:1}

	rel = c.Compare(merged)
	assert(rel == 1, "c should happen‑after merged")
	rel = merged.Compare(c)
	assert(rel == -1, "merged should happen‑before c")

	// Test 4: Equality after identical operations.
	d1 := NewVectorClock(0)
	d2 := NewVectorClock(0)

	for i := 0; i < 3; i++ {
		d1.Increment(0)
		d2.Increment(0)
	}
	for i := 0; i < 2; i++ {
		d1.Increment(1)
		d2.Increment(1)
	}
	assert(d1.Compare(d2) == 0, "d1 and d2 should be equal")
	assert(d2.Compare(d1) == 0, "d2 and d1 should be equal")
	assert(d1.String() == d2.String(), "string representations should match")

	// Test 5: Merge does not mutate originals.
	preA := a.Clone()
	preB := b.Clone()
	_ = a.Merge(b) // result ignored
	assert(a.Compare(preA) == 0, "a must remain unchanged after Merge")
	assert(b.Compare(preB) == 0, "b must remain unchanged after Merge")

	fmt.Println("All vector clock tests passed.")
}
