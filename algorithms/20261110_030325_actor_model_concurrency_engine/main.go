package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// Test 1: Simple summation actor.
	var sum int64
	sumHandler := func(msg interface{}) {
		if v, ok := msg.(int); ok {
			atomic.AddInt64(&sum, int64(v))
		}
	}
	sumActor := NewActor(10, sumHandler)
	for i := 1; i <= 5; i++ {
		sumActor.Send(i)
	}
	time.Sleep(10 * time.Millisecond)
	sumActor.Stop()
	if atomic.LoadInt64(&sum) != 15 {
		panic(fmt.Sprintf("Test 1 failed: expected sum 15, got %d", atomic.LoadInt64(&sum)))
	}

	// Test 2: Order preservation.
	var order []int
	var orderMu sync.Mutex
	orderHandler := func(msg interface{}) {
		if v, ok := msg.(int); ok {
			orderMu.Lock()
			order = append(order, v)
			orderMu.Unlock()
		}
	}
	orderActor := NewActor(5, orderHandler)
	for i := 0; i < 3; i++ {
		orderActor.Send(i)
	}
	time.Sleep(10 * time.Millisecond)
	orderActor.Stop()
	expectedOrder := []int{0, 1, 2}
	if len(order) != len(expectedOrder) {
		panic(fmt.Sprintf("Test 2 failed: length mismatch %d vs %d", len(order), len(expectedOrder)))
	}
	for i, v := range expectedOrder {
		if order[i] != v {
			panic(fmt.Sprintf("Test 2 failed: at index %d expected %d got %d", i, v, order[i]))
		}
	}

	// Test 3: Concurrency safety with multiple senders.
	var processed int64
	concurrencyHandler := func(msg interface{}) {
		atomic.AddInt64(&processed, 1)
	}
	concurrencyActor := NewActor(100, concurrencyHandler)

	var wg sync.WaitGroup
	numSenders := 5
	msgsPerSender := 1000
	wg.Add(numSenders)
	for s := 0; s < numSenders; s++ {
		go func() {
			defer wg.Done()
			for i := 0; i < msgsPerSender; i++ {
				concurrencyActor.Send(i)
			}
		}()
	}
	wg.Wait()
	// Allow some time for processing.
	time.Sleep(50 * time.Millisecond)
	concurrencyActor.Stop()
	expectedProcessed := int64(numSenders * msgsPerSender)
	if atomic.LoadInt64(&processed) != expectedProcessed {
		panic(fmt.Sprintf("Test 3 failed: expected %d processed messages, got %d", expectedProcessed, atomic.LoadInt64(&processed)))
	}

	// All tests passed.
	fmt.Println("All actor model tests passed.")
}
