package main

import "fmt"

func main() {
	// Test 1: basic functionality
	arr1 := []int64{1, 2, 3, 4, 5}
	st1 := NewSegmentTree(arr1)

	// Update indices 1..3 by +10
	st1.UpdateRange(1, 3, 10)
	expectedSum1 := int64(45) // 1 + (2+10)+(3+10)+(4+10) +5 = 45
	if got := st1.QueryRange(0, 4); got != expectedSum1 {
		panic(fmt.Sprintf("Test1 failed: expected sum %d, got %d", expectedSum1, got))
	}

	// Query subrange 2..4
	expectedSub2 := int64(13 + 14 + 5) // (3+10)+(4+10)+5 = 32
	if got := st1.QueryRange(2, 4); got != expectedSub2 {
		panic(fmt.Sprintf("Test1 subrange failed: expected %d, got %d", expectedSub2, got))
	}

	// Test 2: multiple updates and queries
	arr2 := []int64{0, 0, 0, 0, 0, 0, 0, 0}
	st2 := NewSegmentTree(arr2)

	// Apply a series of updates
	st2.UpdateRange(0, 7, 5)   // all +5
	st2.UpdateRange(2, 5, -2)  // indices 2..5 -2
	st2.UpdateRange(4, 7, 3)   // indices 4..7 +3

	// Brute-force compute expected array
	expected := make([]int64, 8)
	for i := 0; i < 8; i++ {
		expected[i] = 5
	}
	for i := 2; i <= 5; i++ {
		expected[i] -= 2
	}
	for i := 4; i <= 7; i++ {
		expected[i] += 3
	}
	// Verify each element
	for i := 0; i < 8; i++ {
		if got := st2.QueryRange(i, i); got != expected[i] {
			panic(fmt.Sprintf("Test2 element %d failed: expected %d, got %d", i, expected[i], got))
		}
	}
	// Verify total sum
	var total int64
	for _, v := range expected {
		total += v
	}
	if got := st2.QueryRange(0, 7); got != total {
		panic(fmt.Sprintf("Test2 total sum failed: expected %d, got %d", total, got))
	}

	// Test 3: edge cases - single element array
	arr3 := []int64{42}
	st3 := NewSegmentTree(arr3)
	if got := st3.QueryRange(0, 0); got != 42 {
		panic(fmt.Sprintf("Test3 initial query failed: expected 42, got %d", got))
	}
	st3.UpdateRange(0, 0, -7)
	if got := st3.QueryRange(0, 0); got != 35 {
		panic(fmt.Sprintf("Test3 after update failed: expected 35, got %d", got))
	}

	// Test 4: full range updates and queries on larger array
	n := 100
	arr4 := make([]int64, n)
	for i := 0; i < n; i++ {
		arr4[i] = int64(i)
	}
	st4 := NewSegmentTree(arr4)

	// Update whole range by +1
	st4.UpdateRange(0, n-1, 1)

	// Verify sum matches arithmetic series + n
	expectedSum4 := int64(n*(n-1)/2) + int64(n) // original sum + n*1
	if got := st4.QueryRange(0, n-1); got != expectedSum4 {
		panic(fmt.Sprintf("Test4 full range sum failed: expected %d, got %d", expectedSum4, got))
	}

	// Random mid-range update
	st4.UpdateRange(20, 79, 5)
	// Brute recompute for verification
	for i := 20; i <= 79; i++ {
		arr4[i] += 1 + 5 // previous +1 then +5
	}
	var bruteSum int64
	for _, v := range arr4 {
		bruteSum += v
	}
	if got := st4.QueryRange(0, n-1); got != bruteSum {
		panic(fmt.Sprintf("Test4 after mid update failed: expected %d, got %d", bruteSum, got))
	}

	fmt.Println("All tests passed")
}
