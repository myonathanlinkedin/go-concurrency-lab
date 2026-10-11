package main

import (
	"fmt"
)

func main() {
	// Test 1: Basic LRU behavior.
	cache := NewLRUCache(3)
	sequence := []int{1, 2, 3, 1, 4, 5, 2, 1, 3}
	expectedHits := []bool{false, false, false, true, false, false, false, true, false}
	for i, page := range sequence {
		hit := cache.Get(page)
		if hit != expectedHits[i] {
			panic(fmt.Sprintf("Test1: unexpected hit result at index %d: got %v, want %v", i, hit, expectedHits[i]))
		}
		if !hit {
			cache.Put(page)
		}
	}
	h, m := cache.Stats()
	if h != 2 || m != 7 {
		panic(fmt.Sprintf("Test1: stats mismatch: hits=%d want 2, misses=%d want 7", h, m))
	}

	// Test 2: Zero capacity cache – should always miss.
	zeroCache := NewLRUCache(0)
	for i := 0; i < 5; i++ {
		if zeroCache.Get(i) {
			panic("Test2: zero capacity cache reported a hit")
		}
		zeroCache.Put(i) // should be a no‑op
	}
	h2, m2 := zeroCache.Stats()
	if h2 != 0 || m2 != 5 {
		panic(fmt.Sprintf("Test2: stats mismatch: hits=%d want 0, misses=%d want 5", h2, m2))
	}

	// Test 3: Re‑inserting existing key does not increase size.
	cap5 := NewLRUCache(5)
	for i := 0; i < 5; i++ {
		cap5.Put(i)
	}
	if cap5.size != 5 {
		panic(fmt.Sprintf("Test3: size after initial fill = %d, want 5", cap5.size))
	}
	cap5.Put(2) // promote existing key
	if cap5.size != 5 {
		panic(fmt.Sprintf("Test3: size after promoting existing key = %d, want 5", cap5.size))
	}
	// Insert new key causing eviction.
	cap5.Put(6)
	if cap5.size != 5 {
		panic(fmt.Sprintf("Test3: size after eviction = %d, want 5", cap5.size))
	}
	if _, ok := cap5.lookup[0]; ok {
		panic("Test3: least recently used key (0) was not evicted")
	}

	fmt.Println("All assertions passed.")
}
