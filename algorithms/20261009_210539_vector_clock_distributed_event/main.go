package main

import "fmt"

func main() {
	// Basic creation and increment tests.
	vcA := NewVectorClock()
	vcA.Increment(1) // pid 1 => 1
	vcA.Increment(2) // pid 2 => 1
	vcA.Increment(1) // pid 1 => 2

	if vcA.clock[1] != 2 || vcA.clock[2] != 1 {
		panic("increment failed")
	}

	// Clone should be equal but independent.
	vcB := vcA.Clone()
	if vcB.Compare(vcA) != Equal {
		panic("clone equality failed")
	}
	vcB.Increment(3)
	if vcB.Compare(vcA) != HappensAfter {
		panic("clone independence failed")
	}

	// Merge test.
	vcC := NewVectorClock()
	vcC.Increment(2) // pid 2 => 1
	vcC.Increment(3) // pid 3 => 1
	vcC.Increment(3) // pid 3 => 2

	vcA.Merge(vcC) // vcA should now have max timestamps.
	if vcA.clock[1] != 2 || vcA.clock[2] != 1 || vcA.clock[3] != 2 {
		panic("merge failed")
	}

	// Comparison tests.
	vcX := NewVectorClock()
	vcX.Increment(1) // {1:1}
	vcY := NewVectorClock()
	vcY.Increment(1) // {1:1}
	vcY.Increment(2) // {1:1,2:1}

	if vcX.Compare(vcY) != HappensBefore {
		panic("happens-before detection failed")
	}
	if vcY.Compare(vcX) != HappensAfter {
		panic("happens-after detection failed")
	}
	if vcX.Compare(vcX.Clone()) != Equal {
		panic("self equality failed")
	}

	// Concurrent case.
	vcP := NewVectorClock()
	vcP.Increment(1) // {1:1}
	vcQ := NewVectorClock()
	vcQ.Increment(2) // {2:1}
	if vcP.Compare(vcQ) != Concurrent {
		panic("concurrent detection failed")
	}
	if vcQ.Compare(vcP) != Concurrent {
		panic("concurrent detection failed (reverse)")
	}

	// Edge case: empty clocks are equal.
	empty1 := NewVectorClock()
	empty2 := NewVectorClock()
	if empty1.Compare(empty2) != Equal {
		panic("empty clocks equality failed")
	}

	// Demonstration output (not part of assertions).
	// The following prints are safe and illustrate the final state.
	// They do not affect test outcomes.
	_ = fmt.Sprintf // silence unused import warning if fmt is not otherwise used
}
