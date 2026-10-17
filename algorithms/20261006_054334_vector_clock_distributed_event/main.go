package main

import (
	"fmt"
	"log"
)

func assert(cond bool, msg string) {
	if !cond {
		log.Fatalf("Assertion failed: %s", msg)
	}
}

func main() {
	// Define nodes in the system.
	nodes := []string{"A", "B", "C"}
	vcA := NewVectorClock(nodes)
	vcB := NewVectorClock(nodes)
	vcC := NewVectorClock(nodes)

	// Simulate events.
	// Event e1 on A
	vcA.Increment("A")
	e1 := Event{ID: "e1", Node: "A", Clock: vcA.Clone()}

	// Event e2 on B, after receiving e1
	vcB.Merge(vcA)
	vcB.Increment("B")
	e2 := Event{ID: "e2", Node: "B", Clock: vcB.Clone()}

	// Event e3 on C, concurrent with e2 (no knowledge of e2)
	vcC.Increment("C")
	e3 := Event{ID: "e3", Node: "C", Clock: vcC.Clone()}

	// Event e4 on A, after receiving e2
	vcA.Merge(vcB)
	vcA.Increment("A")
	e4 := Event{ID: "e4", Node: "A", Clock: vcA.Clone()}

	// Assertions on pairwise comparisons
	assert(e1.Clock.Compare(e2.Clock) == Before, "e1 should be before e2")
	assert(e2.Clock.Compare(e1.Clock) == After, "e2 should be after e1")
	assert(e2.Clock.Compare(e3.Clock) == Concurrent, "e2 and e3 should be concurrent")
	assert(e4.Clock.Compare(e2.Clock) == After, "e4 should be after e2")
	assert(e4.Clock.Compare(e1.Clock) == After, "e4 should be after e1")
	assert(e4.Clock.Compare(e3.Clock) == Concurrent, "e4 and e3 should be concurrent")

	// Test ordering of a mixed set of events.
	events := []Event{e3, e4, e2, e1}
	ordered := OrderEvents(events)

	// Expected order: e1, e2, e4, e3 (e3 concurrent, placed last due to ID ordering)
	expected := []string{"e1", "e2", "e4", "e3"}
	if len(ordered) != len(expected) {
		log.Fatalf("Ordered length mismatch")
	}
	for i, ev := range ordered {
		assert(ev.ID == expected[i], fmt.Sprintf("expected %s at position %d, got %s", expected[i], i, ev.ID))
	}

	fmt.Println("All vector clock assertions passed.")
}
