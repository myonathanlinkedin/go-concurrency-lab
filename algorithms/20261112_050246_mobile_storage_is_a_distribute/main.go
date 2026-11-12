package main

import (
	"fmt"
	"sync"
)

func main() {
	// Basic functionality test.
	ms := NewMobileStorage(2)

	// Put two items; both should be cached.
	ms.Put("a", "alpha")
	ms.Put("b", "bravo")
	if ms.CacheSize() != 2 {
		panic(fmt.Sprintf("expected cache size 2, got %d", ms.CacheSize()))
	}
	if ms.Size() != 2 {
		panic(fmt.Sprintf("expected remote size 2, got %d", ms.Size()))
	}

	// Access 'a' to make it most-recently used.
	if v, ok := ms.Get("a"); !ok || v != "alpha" {
		panic("failed to get existing key 'a'")
	}

	// Insert third item, causing eviction of least-recently used ('b').
	ms.Put("c", "charlie")
	if ms.CacheSize() != 2 {
		panic(fmt.Sprintf("expected cache size 2 after eviction, got %d", ms.CacheSize()))
	}
	if _, ok := ms.cache["b"]; ok {
		panic("key 'b' should have been evicted from cache")
	}
	// 'b' should still exist in remote storage.
	if v, ok := ms.Get("b"); !ok || v != "bravo" {
		panic("failed to retrieve evicted key 'b' from remote")
	}
	// After fetching, 'b' should be back in cache, evicting the LRU ('c' now).
	if _, ok := ms.cache["c"]; ok {
		panic("key 'c' should have been evicted after fetching 'b'")
	}
	if _, ok := ms.cache["b"]; !ok {
		panic("key 'b' should be present in cache after fetch")
	}

	// Delete a key and ensure it disappears from both layers.
	ms.Delete("a")
	if _, ok := ms.Get("a"); ok {
		panic("key 'a' should have been deleted")
	}
	if _, ok := ms.remote["a"]; ok {
		panic("key 'a' still present in remote after delete")
	}
	if _, ok := ms.cache["a"]; ok {
		panic("key 'a' still present in cache after delete")
	}

	// Concurrency test: perform many puts and gets from multiple goroutines.
	const workers = 10
	const opsPerWorker = 1000
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerWorker; j++ {
				key := fmt.Sprintf("w%d-%d", id, j)
				val := fmt.Sprintf("val-%d-%d", id, j)
				ms.Put(key, val)
				if got, ok := ms.Get(key); !ok || got != val {
					panic(fmt.Sprintf("concurrent get mismatch for %s: got %s, want %s", key, got, val))
				}
			}
		}(i)
	}
	wg.Wait()

	// Verify total number of remote entries matches expected count.
	expectedRemote := 2 + // initial a,b,c (a deleted, b and c remain)
		workers*opsPerWorker // all puts from workers
	if ms.Size() != expectedRemote {
		panic(fmt.Sprintf("remote size mismatch: expected %d, got %d", expectedRemote, ms.Size()))
	}

	// Final sanity print.
	fmt.Println("All assertions passed.")
}

// The following unused variable assignments silence the Go compiler's
// 'unused variable' checks for the imported packages that are only used
// in assertions or panic messages.
var (
	_ = fmt.Sprintf
	_ = sync.WaitGroup{}
)
