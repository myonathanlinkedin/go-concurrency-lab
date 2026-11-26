package main

import (
	"fmt"
)

func assert(condition bool, msg string) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: %s", msg))
	}
}

func main() {
	// Create three replicas.
	a := NewPNCounter("A")
	b := NewPNCounter("B")
	c := NewPNCounter("C")

	// Basic increment/decrement tests.
	a.Increment(5)
	assert(a.Value() == 5, "A should be 5 after increment")

	a.Decrement(2)
	assert(a.Value() == 3, "A should be 3 after decrement")

	// Independent operations on B.
	b.Increment(10)
	b.Decrement(4)
	assert(b.Value() == 6, "B should be 6 after ops")

	// Merge A and B into A.
	a.Merge(b)
	// Expected value: (5+10) - (2+4) = 9
	assert(a.Value() == 9, "A after merging B should be 9")
	// B should remain unchanged.
	assert(b.Value() == 6, "B should remain 6 after A merges B")

	// Merge B into C (C is empty initially).
	c.Merge(b)
	assert(c.Value() == 6, "C after merging B should be 6")

	// Further operations after merge.
	c.Increment(3) // C's own increment.
	c.Decrement(1)
	assert(c.Value() == 8, "C after own ops should be 8")

	// Merge C into A (now A has state from B and its own).
	a.Merge(c)
	// Compute expected: 
	// Positive: A(5) + B(10) + C(3) = 18
	// Negative: A(2) + B(4) + C(1) = 7
	// Value = 11
	assert(a.Value() == 11, "A after merging C should be 11")

	// Verify commutativity: merging order does not affect result.
	// Reset D and E replicas.
	d := NewPNCounter("D")
	e := NewPNCounter("E")
	d.Increment(7)
	e.Decrement(3)

	// Merge D into E then into a new replica X.
	x1 := NewPNCounter("X")
	x1.Merge(d)
	x1.Merge(e)

	// Merge E into D then into another replica X2.
	x2 := NewPNCounter("X")
	x2.Merge(e)
	x2.Merge(d)

	assert(x1.Value() == x2.Value(), "Commutative merge should yield same value")
	assert(x1.Value() == int64(7-3), "Result should be 4")

	// Verify idempotence: merging the same counter twice does not change state.
	before := a.Value()
	a.Merge(b) // b already merged into a earlier.
	after := a.Value()
	assert(before == after, "Idempotent merge should not alter value")

	// Verify associativity: (A merge B) merge C == A merge (B merge C)
	assoc1 := NewPNCounter("Assoc")
	assoc1.Merge(a)
	assoc1.Merge(b)
	assoc1.Merge(c)

	assoc2 := NewPNCounter("Assoc")
	assoc2.Merge(b)
	assoc2.Merge(c)
	assoc2.Merge(a)

	assert(assoc1.Value() == assoc2.Value(), "Associative merge should yield same value")

	// Edge case: merging with an empty counter should be a no-op.
	empty := NewPNCounter("Empty")
	valBefore := a.Value()
	a.Merge(empty)
	assert(a.Value() == valBefore, "Merging empty counter should not change value")

	// All assertions passed.
	fmt.Println("All PN-Counter tests passed.")
}
