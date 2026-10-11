package main

import "fmt"

func assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: " + msg)
	}
}

func testBasicConstruction() {
	arr := []int64{1, 2, 3, 4, 5}
	st := NewSegmentTree(arr)
	total := st.Query(0, len(arr))
	assert(total == 15, "total sum mismatch")
	assert(st.Query(1, 4) == 9, "subrange sum mismatch")
}

func testSingleElementUpdates() {
	arr := []int64{10}
	st := NewSegmentTree(arr)
	assert(st.Query(0, 1) == 10, "initial single element")
	st.Update(0, 1, 5)
	assert(st.Query(0, 1) == 15, "after update")
	st.Update(0, 1, -3)
	assert(st.Query(0, 1) == 12, "after negative update")
}

func testRangeUpdatesAndQueries() {
	arr := []int64{0, 0, 0, 0, 0, 0}
	st := NewSegmentTree(arr)
	st.Update(1, 5, 2) // add 2 to indices 1..4
	assert(st.Query(0, 6) == 8, "total after first update")
	assert(st.Query(1, 5) == 8, "subrange after first update")
	st.Update(2, 4, 3) // add 3 to indices 2..3
	assert(st.Query(0, 6) == 14, "total after second update")
	assert(st.Query(2, 4) == 10, "inner subrange after second update")
}

func testEdgeCases() {
	// Empty array
	empty := []int64{}
	stEmpty := NewSegmentTree(empty)
	assert(stEmpty.Query(0, 0) == 0, "empty tree query")
	// Full range updates
	arr := []int64{1, 2, 3, 4}
	st := NewSegmentTree(arr)
	st.Update(0, 4, 1)
	assert(st.Query(0, 4) == 13, "full range update")
	// Single element range
	st.Update(2, 3, 5)
	assert(st.Query(2, 3) == 9, "single element after update")
}

func main() {
	fmt.Println("Running KACTL Segment Tree tests...")
	testBasicConstruction()
	testSingleElementUpdates()
	testRangeUpdatesAndQueries()
	testEdgeCases()
	fmt.Println("All tests passed.")
}
