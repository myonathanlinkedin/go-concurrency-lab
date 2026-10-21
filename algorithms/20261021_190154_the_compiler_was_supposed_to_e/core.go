package main

type IntHeap struct {
	data []int
}

func NewIntHeap() *IntHeap {
	return &IntHeap{data: make([]int, 0)}
}

func (h *IntHeap) Len() int { return len(h.data) }

func (h *IntHeap) Peek() (int, bool) {
	if len(h.data) == 0 {
		return 0, false
	}
	return h.data[0], true
}

func (h *IntHeap) Push(val int) {
	h.data = append(h.data, val)
	h.up(len(h.data) - 1)
}

func (h *IntHeap) Pop() (int, bool) {
	if len(h.data) == 0 {
		return 0, false
	}
	root := h.data[0]
	last := h.data[len(h.data)-1]
	h.data = h.data[:len(h.data)-1]
	if len(h.data) > 0 {
		h.data[0] = last
		h.down(0)
	}
	return root, true
}

func (h *IntHeap) up(i int) {
	for i > 0 {
		parent := (i - 1) >> 1
		if h.data[parent] <= h.data[i] {
			break
		}
		h.data[parent], h.data[i] = h.data[i], h.data[parent]
		i = parent
	}
}

func (h *IntHeap) down(i int) {
	n := len(h.data)
	for {
		left := (i << 1) + 1
		right := left + 1
		smallest := i
		if left < n && h.data[left] < h.data[smallest] {
			smallest = left
		}
		if right < n && h.data[right] < h.data[smallest] {
			smallest = right
		}
		if smallest == i {
			break
		}
		h.data[i], h.data[smallest] = h.data[smallest], h.data[i]
		i = smallest
	}
}

func HeapSort(arr []int) []int {
	h := NewIntHeap()
	for _, v := range arr {
		h.Push(v)
	}
	sorted := make([]int, 0, len(arr))
	for h.Len() > 0 {
		v, _ := h.Pop()
		sorted = append(sorted, v)
	}
	return sorted
}
