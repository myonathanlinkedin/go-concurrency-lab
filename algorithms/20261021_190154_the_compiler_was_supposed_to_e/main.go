package main

import "fmt"

func main() {
	// Test heap push/pop order
	h := NewIntHeap()
	values := []int{5, 3, 8, 1, 9, 2}
	for _, v := range values {
		h.Push(v)
	}
	if h.Len() != len(values) {
		panic(fmt.Sprintf("expected heap length %d, got %d", len(values), h.Len()))
	}
	// Pop should return sorted ascending
	sorted := make([]int, 0, len(values))
	for h.Len() > 0 {
		v, ok := h.Pop()
		if !ok {
			panic("pop returned false on non-empty heap")
		}
		sorted = append(sorted, v)
	}
	expected := []int{1, 2, 3, 5, 8, 9}
	for i, v := range expected {
		if sorted[i] != v {
			panic(fmt.Sprintf("expected %d at index %d, got %d", v, i, sorted[i]))
		}
	}

	// Test Peek on empty heap
	_, ok := h.Peek()
	if ok {
		panic("peek returned true on empty heap")
	}

	// Test HeapSort
	unsorted := []int{10, -1, 7, 3, 0}
	sortedArr := HeapSort(unsorted)
	expectedSorted := []int{-1, 0, 3, 7, 10}
	if len(sortedArr) != len(expectedSorted) {
		panic("HeapSort returned wrong length")
	}
	for i, v := range expectedSorted {
		if sortedArr[i] != v {
			panic(fmt.Sprintf("HeapSort expected %d at index %d, got %d", v, i, sortedArr[i]))
		}
	}

	// Edge case: pop from empty heap
	_, ok = h.Pop()
	if ok {
		panic("pop returned true on empty heap")
	}

	// Edge case: push duplicate values
	h.Push(4)
	h.Push(4)
	h.Push(4)
	if h.Len() != 3 {
		panic("expected heap length 3 after pushing duplicates")
	}
	for i := 0; i < 3; i++ {
		v, ok := h.Pop()
		if !ok || v != 4 {
			panic(fmt.Sprintf("expected 4 on pop %d, got %v, ok=%v", i, v, ok))
		}
	}

	// All tests passed
	fmt.Println("All tests passed")
}
