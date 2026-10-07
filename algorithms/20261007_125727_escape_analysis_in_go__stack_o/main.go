package main

import (
	"fmt"
	"time"
)

func main() {
	// Allocate many escaping nodes; their finalizers should run.
	for i := 0; i < 1000; i++ {
		_ = newEscapingNode(i)
	}
	triggerGC()
	time.Sleep(100 * time.Millisecond)

	count := getFinalizerCount()
	if count == 0 {
		panic(fmt.Sprintf("expected finalizers to run, got %d", count))
	}

	// Allocate many non‑escaping nodes; finalizer count must stay unchanged.
	prev := count
	for i := 0; i < 1000; i++ {
		_ = newNonEscapingNode(i)
	}
	triggerGC()
	time.Sleep(100 * time.Millisecond)

	count2 := getFinalizerCount()
	if count2 != prev {
		panic(fmt.Sprintf("finalizer count changed for non‑escaping nodes: before %d after %d", prev, count2))
	}

	// Additional sanity checks.
	n1 := newEscapingNode(1)
	n2 := newEscapingNode(2)
	if n1 == n2 {
		panic("expected distinct pointers for separate allocations")
	}

	v1 := newNonEscapingNode(1)
	v2 := newNonEscapingNode(2)
	if &v1 == &v2 {
		panic("unexpected identical addresses for stack values")
	}

	fmt.Println("All escape analysis assertions passed")
}
