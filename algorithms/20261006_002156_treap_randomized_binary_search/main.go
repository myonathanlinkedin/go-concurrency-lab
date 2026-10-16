package main

import (
	"fmt"
	"sort"
)

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func main() {
	// Create a new Treap
	t := NewTreap()

	// Test insertion
	keysToInsert := []int{50, 30, 70, 20, 40, 60, 80}
	for _, k := range keysToInsert {
		t.Insert(k)
	}

	// Verify all inserted keys are searchable
	for _, k := range keysToInsert {
		assert(t.Search(k), fmt.Sprintf("Key %d should be found after insertion", k))
	}

	// Verify in-order traversal yields sorted keys
	sorted := t.InOrder()
	expected := make([]int, len(keysToInsert))
	copy(expected, keysToInsert)
	sort.Ints(expected)
	assert(len(sorted) == len(expected), "In-order traversal length mismatch")
	for i, v := range sorted {
		assert(v == expected[i], fmt.Sprintf("In-order mismatch at index %d: got %d, want %d", i, v, expected[i]))
	}

	// Test duplicate insertion (should be ignored)
	t.Insert(30)
	assert(len(t.InOrder()) == len(keysToInsert), "Duplicate insertion altered treap size")

	// Test deletion
	keysToDelete := []int{20, 70}
	for _, k := range keysToDelete {
		t.Delete(k)
	}

	// Verify deleted keys are no longer searchable
	for _, k := range keysToDelete {
		assert(!t.Search(k), fmt.Sprintf("Key %d should not be found after deletion", k))
	}

	// Verify remaining keys
	remaining := []int{50, 30, 40, 60, 80}
	sorted = t.InOrder()
	assert(len(sorted) == len(remaining), "Remaining keys count mismatch")
	sort.Ints(remaining)
	for i, v := range sorted {
		assert(v == remaining[i], fmt.Sprintf("Remaining key mismatch at index %d: got %d, want %d", i, v, remaining[i]))
	}

	// Edge case: delete non-existent key
	t.Delete(999) // should not panic or alter treap
	assert(len(t.InOrder()) == len(remaining), "Deleting non-existent key altered treap size")

	fmt.Println("All tests passed successfully.")
}
