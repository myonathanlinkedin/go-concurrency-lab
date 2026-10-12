package main

import (
	"fmt"
)

func main() {
	// Basic insertion and median checks.
	h := NewHatchet()
	h.Insert(10)
	if m := h.Median(); m != 10 {
		panic(fmt.Sprintf("expected median 10, got %v", m))
	}
	h.Insert(20)
	if m := h.Median(); m != 15 {
		panic(fmt.Sprintf("expected median 15, got %v", m))
	}
	h.Insert(30)
	if m := h.Median(); m != 20 {
		panic(fmt.Sprintf("expected median 20, got %v", m))
	}
	h.Insert(40)
	if m := h.Median(); m != 25 {
		panic(fmt.Sprintf("expected median 25, got %v", m))
	}
	// Removal tests.
	h.Remove(20)
	if m := h.Median(); m != 30 {
		panic(fmt.Sprintf("expected median 30 after removing 20, got %v", m))
	}
	h.Remove(10)
	if m := h.Median(); m != 35 {
		panic(fmt.Sprintf("expected median 35 after removing 10, got %v", m))
	}
	h.Remove(30)
	if m := h.Median(); m != 40 {
		panic(fmt.Sprintf("expected median 40 after removing 30, got %v", m))
	}
	h.Remove(40)
	// Expect panic on empty median.
	func() {
		defer func() {
			if r := recover(); r == nil {
				panic("expected panic on median of empty Hatchet")
			}
		}()
		_ = h.Median()
	}()

	// Stress test with duplicate values.
	h2 := NewHatchet()
	values := []int{5, 1, 5, 3, 5, 2, 5}
	for _, v := range values {
		h2.Insert(v)
	}
	if m := h2.Median(); m != 5 {
		panic(fmt.Sprintf("expected median 5 with duplicates, got %v", m))
	}
	// Remove all 5s one by one.
	for i := 0; i < 4; i++ {
		h2.Remove(5)
	}
	if m := h2.Median(); m != 2.5 {
		panic(fmt.Sprintf("expected median 2.5 after removing 5s, got %v", m))
	}

	// Large random test (deterministic sequence).
	h3 := NewHatchet()
	seq := []int{100, 90, 80, 70, 60, 50, 40, 30, 20, 10}
	for i, v := range seq {
		h3.Insert(v)
		// Compute expected median manually.
		sorted := make([]int, i+1)
		copy(sorted, seq[:i+1])
		// Simple insertion sort for small i.
		for j := i; j > 0 && sorted[j] < sorted[j-1]; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
		var expected float64
		n := len(sorted)
		if n%2 == 1 {
			expected = float64(sorted[n/2])
		} else {
			expected = (float64(sorted[n/2-1]) + float64(sorted[n/2])) / 2.0
		}
		if m := h3.Median(); m != expected {
			panic(fmt.Sprintf("step %d: expected median %v, got %v", i, expected, m))
		}
	}
	fmt.Println("All Hatchet tests passed.")
}
