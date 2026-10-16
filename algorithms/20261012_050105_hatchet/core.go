package main

import (
	"container/heap"
	"fmt"
)

// IntMinHeap implements a min-heap for ints.
type IntMinHeap []int

func (h IntMinHeap) Len() int           { return len(h) }
func (h IntMinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntMinHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntMinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// IntMaxHeap implements a max-heap for ints.
type IntMaxHeap []int

func (h IntMaxHeap) Len() int           { return len(h) }
func (h IntMaxHeap) Less(i, j int) bool { return h[i] > h[j] } // reverse order
func (h IntMaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntMaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntMaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Hatchet maintains a multiset of integers and can return the median in O(1).
type Hatchet struct {
	lower   *IntMaxHeap      // max-heap for the lower half
	upper   *IntMinHeap      // min-heap for the upper half
	delayed map[int]int      // lazy deletion counts
	size    int              // total number of elements (excluding delayed deletions)
}

// NewHatchet creates an empty Hatchet.
func NewHatchet() *Hatchet {
	l := &IntMaxHeap{}
	u := &IntMinHeap{}
	heap.Init(l)
	heap.Init(u)
	return &Hatchet{
		lower:   l,
		upper:   u,
		delayed: make(map[int]int),
		size:    0,
	}
}

// Insert adds a value to the Hatchet.
func (h *Hatchet) Insert(val int) {
	if h.lower.Len() == 0 || val <= (*h.lower)[0] {
		heap.Push(h.lower, val)
	} else {
		heap.Push(h.upper, val)
	}
	h.size++
	h.rebalance()
}

// Remove deletes one occurrence of val from the Hatchet, if present.
func (h *Hatchet) Remove(val int) {
	if h.size == 0 {
		return
	}
	// Record lazy deletion.
	h.delayed[val]++
	h.size--

	// Clean tops if they are scheduled for deletion.
	if h.lower.Len() > 0 && val <= (*h.lower)[0] {
		h.prune(h.lower)
	} else {
		h.prune(h.upper)
	}
	h.rebalance()
}

// Median returns the median of the current multiset.
// If the size is even, the median is the average of the two middle values.
func (h *Hatchet) Median() float64 {
	if h.size == 0 {
		panic("median of empty Hatchet")
	}
	h.prune(h.lower)
	h.prune(h.upper)

	if h.lower.Len() > h.upper.Len() {
		return float64((*h.lower)[0])
	}
	if h.upper.Len() > h.lower.Len() {
		return float64((*h.upper)[0])
	}
	// even number of elements
	return (float64((*h.lower)[0]) + float64((*h.upper)[0])) / 2.0
}

// rebalance ensures the size property: |lower| == |upper| or |lower| == |upper|+1
func (h *Hatchet) rebalance() {
	// Ensure lower has at most one more element than upper.
	if h.lower.Len() > h.upper.Len()+1 {
		h.prune(h.lower)
		val := heap.Pop(h.lower).(int)
		heap.Push(h.upper, val)
	} else if h.upper.Len() > h.lower.Len() {
		h.prune(h.upper)
		val := heap.Pop(h.upper).(int)
		heap.Push(h.lower, val)
	}
	// After moving, prune possible delayed elements at tops.
	h.prune(h.lower)
	h.prune(h.upper)
}

// prune removes elements from the top of the heap that are marked for lazy deletion.
func (h *Hatchet) prune(hp heap.Interface) {
	for hp.Len() > 0 {
		var top int
		switch heapTyped := hp.(type) {
		case *IntMaxHeap:
			top = (*heapTyped)[0]
		case *IntMinHeap:
			top = (*heapTyped)[0]
		default:
			panic("unknown heap type")
		}
		if cnt, ok := h.delayed[top]; ok && cnt > 0 {
			// Remove the top element.
			heap.Pop(hp)
			if cnt == 1 {
				delete(h.delayed, top)
			} else {
				h.delayed[top] = cnt - 1
			}
			// Continue pruning.
		} else {
			break
		}
	}
}

// DebugString returns a string representation for debugging.
func (h *Hatchet) DebugString() string {
	h.prune(h.lower)
	h.prune(h.upper)
	return fmt.Sprintf("lower:%v upper:%v size:%d delayed:%v", *h.lower, *h.upper, h.size, h.delayed)
}
