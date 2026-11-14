package main

import (
	"fmt"
)

func assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: " + msg)
	}
}

func testBasicOperations() {
	cache := NewLRUCache(2)
	cache.Put(1, 1)
	cache.Put(2, 2)

	v, ok := cache.Get(1)
	assert(ok && v == 1, "Get existing key 1")

	cache.Put(3, 3) // evicts key 2
	_, ok = cache.Get(2)
	assert(!ok, "Key 2 should be evicted")

	cache.Put(4, 4) // evicts key 1
	_, ok = cache.Get(1)
	assert(!ok, "Key 1 should be evicted")

	v, ok = cache.Get(3)
	assert(ok && v == 3, "Get key 3 after eviction")

	v, ok = cache.Get(4)
	assert(ok && v == 4, "Get key 4 after eviction")
}

func testUpdateExisting() {
	cache := NewLRUCache(2)
	cache.Put(1, 1)
	cache.Put(1, 10) // update existing key

	v, ok := cache.Get(1)
	assert(ok && v == 10, "Update existing key 1")

	cache.Put(2, 2)
	cache.Put(3, 3) // evicts key 2 (key 1 was recently used)

	_, ok = cache.Get(2)
	assert(!ok, "Key 2 should be evicted after update")

	v, ok = cache.Get(1)
	assert(ok && v == 10, "Key 1 should remain after eviction")
}

func testZeroCapacity() {
	cache := NewLRUCache(0)
	cache.Put(1, 1)
	_, ok := cache.Get(1)
	assert(!ok, "Cache with zero capacity must not store items")
}

func testCapacityOne() {
	cache := NewLRUCache(1)
	cache.Put(1, 1)
	cache.Put(2, 2) // evicts key 1

	_, ok := cache.Get(1)
	assert(!ok, "Key 1 should be evicted in capacity‑1 cache")

	v, ok := cache.Get(2)
	assert(ok && v == 2, "Key 2 should be present in capacity‑1 cache")
}

func main() {
	fmt.Println("Running LRU Cache unit tests...")
	testBasicOperations()
	testUpdateExisting()
	testZeroCapacity()
	testCapacityOne()
	fmt.Println("All tests passed.")
}
