package main

import (
	"fmt"
)

func assert(condition bool, msg string) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: %s", msg))
	}
}

// TestIncrement verifies that Increment updates only the targeted entry.
func TestIncrement() {
	vc := NewVectorClock()
	vc.Increment(1)
	vc.Increment(1)
	vc.Increment(2)

	assert(vc.clock[1] == 2, "pid 1 should be 2")
	assert(vc.clock[2] == 1, "pid 2 should be 1")
	assert(len(vc.clock) == 2, "exactly two entries expected")
}

// TestMerge checks element‑wise maximum semantics.
func TestMerge() {
	a := NewVectorClock()
	b := NewVectorClock()

	a.Increment(1) // a: {1:1}
	a.Increment(2) // a: {1:1,2:1}
	b.Increment(2) // b: {2:1}
	b.Increment(2) // b: {2:2}
	b.Increment(3) // b: {2:2,3:1}

	a.Merge(b) // expected {1:1,2:2,3:1}
	assert(a.clock[1] == 1, "pid 1 should remain 1")
	assert(a.clock[2] == 2, "pid 2 should become 2")
	assert(a.clock[3] == 1, "pid 3 should be introduced as 1")
}

// TestCompare covers all four relation outcomes.
func TestCompare() {
	// Equal clocks
	c1 := NewVectorClock()
	c2 := NewVectorClock()
	assert(c1.Compare(c2) == RelationEqual, "empty clocks should be equal")

	// Happens‑before
	c1.Increment(1) // {1:1}
	assert(c1.Compare(c2) == RelationHappensAfter, "c1 after empty")
	assert(c2.Compare(c1) == RelationHappensBefore, "empty before c1")

	// Concurrent scenario
	c2.Increment(2) // {2:1}
	rel := c1.Compare(c2)
	assert(rel == RelationConcurrent, fmt.Sprintf("expected concurrent, got %v", rel))

	// Merge makes them comparable
	c1.Merge(c2) // {1:1,2:1}
	c2.Merge(c1) // {1:1,2:1}
	assert(c1.Compare(c2) == RelationEqual, "after merge should be equal")
}

// TestString ensures deterministic formatting.
func TestString() {
	vc := NewVectorClock()
	vc.Increment(5)
	vc.Increment(2)
	vc.Increment(5)
	// Expected order: 2 then 5
	expected := "{2:1,5:2}"
	assert(vc.String() == expected, fmt.Sprintf("string mismatch: got %s", vc.String()))
}

// main runs all self‑contained tests.
func main() {
	fmt.Println("Running VectorClock tests...")
	TestIncrement()
	TestMerge()
	TestCompare()
	TestString()
	fmt.Println("All tests passed.")
}
