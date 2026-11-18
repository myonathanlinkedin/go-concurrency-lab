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
	// Create a skip list with reasonable parameters.
	sl := NewSkipList(16, 0.5)

	// Edge case: search in empty list.
	_, found := sl.Search(10)
	assert(!found, "search in empty list should fail")

	// Insert a series of timestamps with string values.
	keys := []int64{5, 1, 3, 8, 7, 2, 10}
	for _, k := range keys {
		inserted := sl.Insert(k, fmt.Sprintf("val-%d", k))
		assert(inserted, fmt.Sprintf("key %d should be inserted", k))
	}

	// Duplicate insert should update, not increase size.
	updated := sl.Insert(5, "updated-5")
	assert(!updated, "duplicate insert should return false")
	val, ok := sl.Search(5)
	assert(ok && val == "updated-5", "value for key 5 should be updated")

	// Verify all keys are searchable.
	for _, k := range keys {
		v, ok := sl.Search(k)
		assert(ok, fmt.Sprintf("key %d should be found", k))
		expected := fmt.Sprintf("val-%d", k)
		if k == 5 {
			expected = "updated-5"
		}
		assert(v == expected, fmt.Sprintf("key %d value mismatch", k))
	}

	// Delete a middle key.
	removed := sl.Delete(3)
	assert(removed, "key 3 should be removed")
	_, ok = sl.Search(3)
	assert(!ok, "key 3 should no longer be found")

	// Delete non‑existent key.
	removed = sl.Delete(42)
	assert(!removed, "deleting non‑existent key should return false")

	// Range query tests.
	rangeVals := sl.RangeQuery(2, 8)
	expectedVals := []string{"val-2", "updated-5", "val-7", "val-8"}
	assert(len(rangeVals) == len(expectedVals), "range query length mismatch")
	for i, ev := range expectedVals {
		assert(rangeVals[i] == ev, fmt.Sprintf("range query element %d mismatch", i))
	}

	// Edge range: start > end.
	emptyRange := sl.RangeQuery(9, 4)
	assert(len(emptyRange) == 0, "empty range should return empty slice")

	// Delete all keys.
	for _, k := range []int64{1, 2, 5, 7, 8, 10} {
		removed = sl.Delete(k)
		assert(removed, fmt.Sprintf("key %d should be removed", k))
	}
	// List should be empty now.
	_, found = sl.Search(1)
	assert(!found, "list should be empty after deletions")
	assert(len(sl.RangeQuery(-100, 100)) == 0, "range query on empty list should be empty")

	// Verify internal consistency (no panics) by printing debug string.
	fmt.Println("All assertions passed.")
	_ = sl.DebugString() // ensure method compiles and is reachable
}
