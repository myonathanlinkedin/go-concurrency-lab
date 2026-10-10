package main

import "fmt"

func main() {
	// Basic creation and increment tests.
	vcA := NewVectorClock()
	vcA.Inc(1) // node 1 -> 1
	vcA.Inc(2) // node 2 -> 1
	vcA.Inc(1) // node 1 -> 2

	if vcA.Get(1) != 2 {
		panic(fmt.Sprintf("expected node 1 counter 2, got %d", vcA.Get(1)))
	}
	if vcA.Get(2) != 1 {
		panic(fmt.Sprintf("expected node 2 counter 1, got %d", vcA.Get(2)))
	}
	if vcA.Get(3) != 0 {
		panic("expected missing node to have counter 0")
	}

	// Clone and equality test.
	vcB := vcA.Clone()
	if !vcA.Equals(vcB) {
		panic("clone should be equal to original")
	}
	vcB.Inc(3)
	if vcA.Equals(vcB) {
		panic("cloned clock after mutation must differ")
	}

	// Merge test.
	vcC := NewVectorClock()
	vcC.Inc(2) // 2:1
	vcC.Inc(3) // 3:1
	vcC.Inc(3) // 3:2

	vcA.Merge(vcC) // vcA now has max per node.
	if vcA.Get(1) != 2 {
		panic("node 1 should remain 2 after merge")
	}
	if vcA.Get(2) != 1 {
		panic("node 2 should remain 1 after merge")
	}
	if vcA.Get(3) != 2 {
		panic("node 3 should be 2 after merge")
	}

	// Happens‑before relationship tests.
	vcX := NewVectorClock()
	vcX.Inc(1) // {1:1}
	vcY := NewVectorClock()
	vcY.Inc(1) // {1:1}
	vcY.Inc(2) // {1:1,2:1}

	if !vcX.HappensBefore(vcY) {
		panic("vcX should happen before vcY")
	}
	if vcY.HappensBefore(vcX) {
		panic("vcY should not happen before vcX")
	}
	if vcX.Concurrent(vcY) {
		panic("vcX and vcY are not concurrent")
	}
	if vcX.Equals(vcY) {
		panic("vcX and vcY are not equal")
	}

	// Concurrent clocks test.
	vcP := NewVectorClock()
	vcP.Inc(1) // {1:1}
	vcQ := NewVectorClock()
	vcQ.Inc(2) // {2:1}
	if !vcP.Concurrent(vcQ) {
		panic("vcP and vcQ should be concurrent")
	}
	if vcP.HappensBefore(vcQ) || vcQ.HappensBefore(vcP) {
		panic("concurrent clocks must not have happens‑before relation")
	}

	// Equality after identical operations.
	vcR := NewVectorClock()
	vcR.Inc(5)
	vcR.Inc(5)
	vcS := NewVectorClock()
	vcS.Inc(5)
	vcS.Inc(5)
	if !vcR.Equals(vcS) {
		panic("identical sequences must yield equal clocks")
	}
	if vcR.Concurrent(vcS) {
		panic("equal clocks are not concurrent")
	}

	// String representation sanity check (deterministic ordering).
	expectedStr := "{1:2, 2:1, 3:2}"
	if vcA.String() != expectedStr {
		panic(fmt.Sprintf("unexpected string: got %s, want %s", vcA.String(), expectedStr))
	}

	// Edge case: empty clocks.
	empty1 := NewVectorClock()
	empty2 := NewVectorClock()
	if !empty1.Equals(empty2) {
		panic("two empty clocks must be equal")
	}
	if empty1.Concurrent(empty2) {
		panic("empty clocks are not concurrent")
	}
	if empty1.HappensBefore(empty2) {
		panic("empty clocks do not happen before each other")
	}

	// All assertions passed.
	fmt.Println("All vector clock tests passed.")
}
