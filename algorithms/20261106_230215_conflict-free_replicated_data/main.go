package main

import "fmt"

func main() {
	// Test 1: basic increment/decrement.
	c := NewPNCounter()
	c.Increment("A", 5)
	c.Decrement("A", 2)
	if c.Value() != 3 {
		panic(fmt.Sprintf("Test 1 failed: expected 3, got %d", c.Value()))
	}

	// Test 2: merge two independent replicas.
	a := NewPNCounter()
	b := NewPNCounter()
	a.Increment("A", 10)
	a.Decrement("A", 3) // value 7
	b.Increment("B", 4)
	b.Decrement("B", 1) // value 3

	b.Merge(a) // combine states
	// Expected: (10+4) - (3+1) = 10
	if b.Value() != 10 {
		panic(fmt.Sprintf("Test 2 failed: expected 10, got %d", b.Value()))
	}

	// Test 2b: idempotence of merge.
	before := b.Value()
	b.Merge(a)
	if b.Value() != before {
		panic("Test 2b failed: merge not idempotent")
	}

	// Test 3: commutativity of merge.
	c1 := NewPNCounter()
	c2 := NewPNCounter()
	c1.Increment("A", 2)
	c2.Increment("B", 3)

	c1Copy := c1.Clone()
	c1Copy.Merge(c2)

	c2Copy := c2.Clone()
	c2Copy.Merge(c1)

	if c1Copy.Value() != c2Copy.Value() {
		panic(fmt.Sprintf("Test 3 failed: commutativity broken (%d vs %d)", c1Copy.Value(), c2Copy.Value()))
	}

	fmt.Println("All PN-Counter tests passed")
}
