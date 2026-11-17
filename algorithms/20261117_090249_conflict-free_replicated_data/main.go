package main

import (
	"fmt"
)

func main() {
	// Create two replicas.
	a := NewPNCounter("A")
	b := NewPNCounter("B")

	// Perform operations on replica A.
	a.Increment(5)   // +5
	a.Decrement(2)   // -2 => net +3
	a.Increment(3)   // +3 => net +6

	// Perform operations on replica B.
	b.Increment(4)   // +4
	b.Decrement(1)   // -1 => net +3
	b.Decrement(2)   // -2 => net +1

	// Verify individual values.
	if a.Value() != 6 {
		panic(fmt.Sprintf("replica A expected value 6, got %d", a.Value()))
	}
	if b.Value() != 1 {
		panic(fmt.Sprintf("replica B expected value 1, got %d", b.Value()))
	}

	// Merge B into A and verify convergence.
	a.Merge(b)
	if a.Value() != 7 {
		panic(fmt.Sprintf("after merge, replica A expected value 7, got %d", a.Value()))
	}
	// Merge A into B (should be symmetric).
	b.Merge(a)
	if b.Value() != 7 {
		panic(fmt.Sprintf("after merge, replica B expected value 7, got %d", b.Value()))
	}

	// Idempotence: merging again should not change value.
	aPrev := a.Value()
	a.Merge(b)
	if a.Value() != aPrev {
		panic("idempotence violated on merge")
	}

	// Commutativity test with a third replica.
	c := NewPNCounter("C")
	c.Increment(10) // +10
	c.Decrement(5) // -5 => net +5

	// Merge order 1: ((a merge b) merge c)
	ab := a.Clone()
	ab.Merge(b)
	ab.Merge(c)
	val1 := ab.Value()

	// Merge order 2: (a merge (b merge c))
	bc := b.Clone()
	bc.Merge(c)
	abc := a.Clone()
	abc.Merge(bc)
	val2 := abc.Value()

	if val1 != val2 {
		panic(fmt.Sprintf("commutativity failed: %d vs %d", val1, val2))
	}
	if val1 != 12 {
		panic(fmt.Sprintf("final converged value expected 12, got %d", val1))
	}

	// Edge case: merging with self should be a no-op.
	selfPrev := a.Value()
	a.Merge(a)
	if a.Value() != selfPrev {
		panic("merge with self altered state")
	}

	// Demonstrate that negative deltas are rejected.
	func() {
		defer func() {
			if r := recover(); r == nil {
				panic("expected panic on negative increment")
			}
		}()
		a.Increment(-1)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				panic("expected panic on zero decrement")
			}
		}()
		b.Decrement(0)
	}()

	fmt.Println("All PN-Counter assertions passed.")
}
