package main

import (
	"fmt"
)

func main() {
	// Parameters: 1% error, 99% confidence
	epsilon := 0.01
	delta := 0.01
	tracker := NewHeavyHittersTracker(epsilon, delta)

	// Simulated stream
	stream := []struct {
		item  string
		count uint64
	}{
		{"apple", 120},
		{"banana", 80},
		{"cherry", 30},
		{"date", 5},
		{"elderberry", 2},
		{"fig", 1},
		{"grape", 150},
		{"apple", 30},   // additional occurrences
		{"banana", 20},
		{"cherry", 10},
	}

	// Insert into tracker
	for _, entry := range stream {
		tracker.Add(entry.item, entry.count)
	}

	// Verify estimates are within theoretical bounds
	totalUpdates := uint64(0)
	for _, e := range stream {
		totalUpdates += e.count
	}
	// The additive error bound is epsilon * totalUpdates
	allowedError := uint64(epsilon * float64(totalUpdates))

	for _, e := range stream {
		est := tracker.Estimate(e.item)
		if est < e.count {
			panic(fmt.Sprintf("underestimate for %s: got %d, want at least %d", e.item, est, e.count))
		}
		if est > e.count+allowedError {
			panic(fmt.Sprintf("overestimate for %s: got %d, allowed %d", e.item, est, e.count+allowedError))
		}
	}

	// Heavy hitter detection: threshold set to 100
	threshold := uint64(100)
	hh := tracker.HeavyHitters(threshold)

	// Expected heavy hitters: apple (150), grape (150)
	expected := map[string]bool{
		"apple": true,
		"grape": true,
	}
	if len(hh) != len(expected) {
		panic(fmt.Sprintf("heavy hitters count mismatch: got %d, want %d", len(hh), len(expected)))
	}
	for _, item := range hh {
		if !expected[item] {
			panic(fmt.Sprintf("unexpected heavy hitter: %s", item))
		}
	}
	// Print success message
	fmt.Println("All assertions passed. Heavy hitters:", hh)
}
