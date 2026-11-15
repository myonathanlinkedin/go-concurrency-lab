package main

import (
	"fmt"
	"sync"
)

func main() {
	// Basic functionality tests
	store := NewStore()

	// Initial set
	v1 := store.Set("alpha", "first")
	if v1 != 1 {
		panic(fmt.Sprintf("expected version 1, got %d", v1))
	}
	val, ver, ok := store.Get("alpha")
	if !ok {
		panic("expected key 'alpha' to exist")
	}
	if val != "first" {
		panic(fmt.Sprintf("expected value 'first', got %s", val))
	}
	if ver != v1 {
		panic(fmt.Sprintf("expected version %d, got %d", v1, ver))
	}
	if store.IsStale("alpha", ver) {
		panic("expected not stale when using current version")
	}
	if store.IsStale("alpha", ver-1) {
		panic("expected stale when using older version")
	}

	// Update same key
	v2 := store.Set("alpha", "second")
	if v2 != 2 {
		panic(fmt.Sprintf("expected version 2 after update, got %d", v2))
	}
	_, ver2, _ := store.Get("alpha")
	if ver2 != v2 {
		panic(fmt.Sprintf("expected version %d after update, got %d", v2, ver2))
	}
	if !store.IsStale("alpha", v1) {
		panic("expected stale when using previous version after update")
	}
	if store.IsStale("alpha", v2) {
		panic("expected not stale when using latest version")
	}

	// Non‑existent key handling
	_, _, ok = store.Get("beta")
	if ok {
		panic("expected key 'beta' to be absent")
	}
	if store.IsStale("beta", 0) {
		panic("expected IsStale to be false for missing key")
	}

	// Concurrency stress test
	const goroutineCount = 10
	const opsPerGoroutine = 1000
	var wg sync.WaitGroup
	wg.Add(goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("k-%d-%d", id, j)
				store.Set(key, fmt.Sprintf("v-%d-%d", id, j))
				if _, _, ok := store.Get(key); !ok {
					panic(fmt.Sprintf("key %s should exist after set", key))
				}
			}
		}(i)
	}
	wg.Wait()

	// Verify that version counter is monotonic and matches total writes
	expectedVersion := uint64(2 + goroutineCount*opsPerGoroutine) // 2 initial writes + concurrent writes
	if store.globalVersion != expectedVersion {
		panic(fmt.Sprintf("expected global version %d, got %d", expectedVersion, store.globalVersion))
	}

	// All assertions passed
	fmt.Println("All tests passed.")
}
