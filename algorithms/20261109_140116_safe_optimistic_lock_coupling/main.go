package main

import (
	"fmt"
	"sync"
)

func main() {
	// Basic sequential sanity checks.
	list := NewList()
	if list.Insert(10, "ten") != true {
		panic("insert failed")
	}
	if list.Insert(5, "five") != true {
		panic("insert failed")
	}
	if list.Insert(10, "ten-again") != false {
		panic("duplicate insert should return false")
	}
	if v, ok := list.Get(10); !ok || v != "ten-again" {
		panic("get after duplicate insert failed")
	}
	if list.Delete(5) != true {
		panic("delete existing key failed")
	}
	if _, ok := list.Get(5); ok {
		panic("deleted key still present")
	}
	if list.Delete(5) != false {
		panic("delete non‑existent key should return false")
	}

	// Concurrent stress test.
	const n = 1000
	var wg sync.WaitGroup
	list = NewList()

	// Concurrent inserts.
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(k int) {
			defer wg.Done()
			if !list.Insert(k, k*2) {
				// duplicate insert is acceptable in race; ignore.
			}
		}(i)
	}
	wg.Wait()

	// Verify all keys present.
	for i := 0; i < n; i++ {
		if v, ok := list.Get(i); !ok || v != i*2 {
			panic(fmt.Sprintf("missing or incorrect value for key %d", i))
		}
	}

	// Concurrent deletes of even keys.
	wg.Add(n / 2)
	for i := 0; i < n; i += 2 {
		go func(k int) {
			defer wg.Done()
			if !list.Delete(k) {
				// deletion may race; ignore false.
			}
		}(i)
	}
	wg.Wait()

	// Verify even keys absent, odd keys present.
	for i := 0; i < n; i++ {
		_, ok := list.Get(i)
		if i%2 == 0 {
			if ok {
				panic(fmt.Sprintf("even key %d should have been deleted", i))
			}
		} else {
			if !ok {
				panic(fmt.Sprintf("odd key %d should be present", i))
			}
		}
	}

	// Final sanity output.
	fmt.Println("All tests passed. Final list snapshot:", list.DebugString())
}
