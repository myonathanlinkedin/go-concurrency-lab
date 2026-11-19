package main

import (
	"fmt"
)

func contains(slice []uint64, val uint64) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func main() {
	// Basic parameters
	width := 1000
	depth := 5
	cms := NewCMS(width, depth)

	// Test data
	items := []uint64{1, 2, 3}
	counts := []uint64{100, 50, 200}

	// Update CMS
	for i, item := range items {
		cms.Update(item, counts[i])
	}

	// Verify estimates are >= actual counts
	for i, item := range items {
		est := cms.Estimate(item)
		if est < counts[i] {
			panic(fmt.Sprintf("Estimate for item %d (%d) is less than actual count %d", item, est, counts[i]))
		}
	}

	// Estimate for unseen item should be zero
	if cms.Estimate(999) != 0 {
		panic("Estimate for unseen item should be zero")
	}

	// Heavy hitters with threshold 80 should return items 1 and 3
	hh := cms.HeavyHitters(items, 80)
	if len(hh) != 2 || !contains(hh, 1) || !contains(hh, 3) {
		panic(fmt.Sprintf("Heavy hitters incorrect: %v", hh))
	}

	// Heavy hitters with threshold 0 should return all items
	hhAll := cms.HeavyHitters(items, 0)
	if len(hhAll) != len(items) {
		panic(fmt.Sprintf("Heavy hitters with threshold 0 returned %d items, expected %d", len(hhAll), len(items)))
	}

	// Update with zero count should not change estimates
	prevEst := cms.Estimate(1)
	cms.Update(1, 0)
	if cms.Estimate(1) != prevEst {
		panic("Update with zero count altered estimate")
	}

	// Verify table dimensions
	if len(cms.table) != depth {
		panic("Depth mismatch in table")
	}
	for i := range cms.table {
		if len(cms.table[i]) != width {
			panic("Width mismatch in table")
		}
	}

	// Edge case: width and depth of 1
	smallCMS := NewCMS(1, 1)
	smallCMS.Update(42, 10)
	if smallCMS.Estimate(42) != 10 {
		panic("Single-cell CMS estimate incorrect")
	}
	if smallCMS.Estimate(43) != 0 {
		panic("Single-cell CMS estimate for unseen item should be zero")
	}

	fmt.Println("All tests passed")
}
