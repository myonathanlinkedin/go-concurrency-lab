package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// BoundedQueue is a thread‑safe bounded blocking queue.
type BoundedQueue[T any] struct {
	mu       sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
	buf      []T
	head     int
	tail     int
	size     int
	capacity int
}

// NewBoundedQueue creates a new queue with the given capacity.
// Capacity must be > 0.
func NewBoundedQueue[T any](capacity int) *BoundedQueue[T] {
	if capacity <= 0 {
		panic("capacity must be greater than zero")
	}
	q := &BoundedQueue[T]{
		buf:      make([]T, capacity),
		capacity: capacity,
	}
	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)
	return q
}

// Enqueue adds an item to the queue, blocking if the queue is full.
func (q *BoundedQueue[T]) Enqueue(item T) {
	q.mu.Lock()
	for q.size == q.capacity {
		q.notFull.Wait()
	}
	q.buf[q.tail] = item
	q.tail = (q.tail + 1) % q.capacity
	q.size++
	q.notEmpty.Signal()
	q.mu.Unlock()
}

// Dequeue removes and returns an item from the queue, blocking if the queue is empty.
func (q *BoundedQueue[T]) Dequeue() T {
	q.mu.Lock()
	for q.size == 0 {
		q.notEmpty.Wait()
	}
	item := q.buf[q.head]
	q.head = (q.head + 1) % q.capacity
	q.size--
	q.notFull.Signal()
	q.mu.Unlock()
	return item
}

// Size returns the current number of elements in the queue.
// It is provided for testing/debugging and does not block.
func (q *BoundedQueue[T]) Size() int {
	q.mu.Lock()
	s := q.size
	q.mu.Unlock()
	return s
}

// -------------------- Unit Test --------------------

func testBoundedQueue() {
	const (
		capacity   = 5
		producers  = 3
		consumers  = 3
		itemsPerP  = 1000
		totalItems = producers * itemsPerP
	)

	q := NewBoundedQueue[int](capacity)

	var produced int64
	var consumed int64
	var sumProduced int64
	var sumConsumed int64

	var wg sync.WaitGroup

	// Producers
	for p := 0; p < producers; p++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			base := id * itemsPerP
			for i := 0; i < itemsPerP; i++ {
				val := base + i + 1 // avoid zero
				q.Enqueue(val)
				atomic.AddInt64(&produced, 1)
				atomic.AddInt64(&sumProduced, int64(val))
			}
		}(p)
	}

	// Consumers
	for c := 0; c < consumers; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				// Stop condition: when all items have been consumed.
				if atomic.LoadInt64(&consumed) >= int64(totalItems) {
					return
				}
				val := q.Dequeue()
				atomic.AddInt64(&consumed, 1)
				atomic.AddInt64(&sumConsumed, int64(val))
			}
		}()
	}

	// Wait for producers to finish, then give consumers time to drain.
	wg.Wait()

	// Ensure queue is empty.
	if q.Size() != 0 {
		panic(fmt.Sprintf("queue not empty after test, size=%d", q.Size()))
	}
	if produced != int64(totalItems) {
		panic(fmt.Sprintf("produced count mismatch: got %d, want %d", produced, totalItems))
	}
	if consumed != int64(totalItems) {
		panic(fmt.Sprintf("consumed count mismatch: got %d, want %d", consumed, totalItems))
	}
	if sumProduced != sumConsumed {
		panic(fmt.Sprintf("sum mismatch: produced %d, consumed %d", sumProduced, sumConsumed))
	}
	fmt.Println("All assertions passed.")
}

// -------------------- Main --------------------

func main() {
	// Run the test with a timeout to avoid deadlocks.
	done := make(chan struct{})
	go func() {
		testBoundedQueue()
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(10 * time.Second):
		panic("test timed out")
	}
}
