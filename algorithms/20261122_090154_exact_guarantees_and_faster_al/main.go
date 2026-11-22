package main

import (
	"fmt"
)

func main() {
	// Helper to run a test case.
	run := func(seq []int, k int, expected int) {
		solver := NewSolver(seq, k)
		got := solver.Longest()
		if got != expected {
			panic(fmt.Sprintf("FAIL: seq=%v k=%d expected=%d got=%d", seq, k, expected, got))
		}
	}

	// Basic sanity checks.
	run([]int{}, 0, 0)                                 // empty sequence
	run([]int{5}, 0, 1)                                // single element
	run([]int{1, 2, 3, 4, 5}, 0, 5)                    // already monotone up
	run([]int{5, 4, 3, 2, 1}, 0, 5)                    // already monotone down
	run([]int{1, 3, 2, 4, 3, 5}, 1, 5)                 // one change allowed
	run([]int{1, 3, 2, 4, 3, 5}, 2, 6)                 // two changes allow full length
	run([]int{1, 2, 1, 2, 1, 2, 1}, 3, 7)              // alternating, 3 changes enough for full length
	run([]int{1, 2, 3, 2, 1, 2, 3, 4}, 2, 8)           // complex pattern
	run([]int{10, 20, 10, 20, 10, 20, 10, 20}, 3, 8)   // repeated up/down
	run([]int{1, 2, 3, 4, 5}, 10, 5)                   // k larger than needed
	run([]int{5, 1, 4, 2, 3}, 2, 4)                    // best subsequence length 4

	// Edge case: k = 0 but sequence not monotone.
	run([]int{1, 3, 2, 4, 5}, 0, 3) // longest monotone subsequence (either 1,2,4,5 not allowed because change)

	// Verify that k >= n-1 yields full length.
	for n := 1; n <= 10; n++ {
		seq := make([]int, n)
		for i := 0; i < n; i++ {
			seq[i] = i % 3 // introduce repeats and changes
		}
		solver := NewSolver(seq, n-1)
		if solver.Longest() != n {
			panic(fmt.Sprintf("FAIL: full length not achieved for n=%d", n))
		}
	}

	fmt.Println("All assertions passed.")
}
