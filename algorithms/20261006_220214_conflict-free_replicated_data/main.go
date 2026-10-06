package main

import (
	"fmt"
	"sync"
	"time"
)

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func main() {
	// Basic functionality test
	c1 := NewPNCounter()
	c2 := NewPNCounter()

	c1.Increment("nodeA", 5)
	c1.Decrement("nodeA", 2)
	assert(c1.Value() == 3, "c1 value should be 3")

	c2.Increment("nodeB", 4)
	c2.Decrement("nodeB", 1)
	assert(c2.Value() == 3, "c2 value should be 3")

	// Merge counters
	c1.Merge(c2)
	assert(c1.Value() == 6, "merged value should be 6")

	// Concurrent increments
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			node := fmt.Sprintf("node%d", id)
			for j := 0; j < 100; j++ {
				c1.Increment(node, 1)
			}
		}(i)
	}
	wg.Wait()
	expected := 6 + 10*100
	assert(c1.Value() == expected, fmt.Sprintf("after concurrent increments expected %d got %d", expected, c1.Value()))

	// Benchmark: simple loop
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		c1.Increment("bench", 1)
	}
	duration := time.Since(start)
	fmt.Printf("Benchmark: incremented 1,000,000 times in %v\n", duration)

	fmt.Println("All assertions passed.")
}
