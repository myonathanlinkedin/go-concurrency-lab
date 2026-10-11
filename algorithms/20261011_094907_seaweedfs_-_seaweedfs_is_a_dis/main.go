package main

import (
	"fmt"
	"time"
)

func main() {
	// ----- Volume tests -----
	vol := &Volume{}
	needles := []Needle{
		{ID: 42, Offset: 1000, Size: 256},
		{ID: 7, Offset: 2000, Size: 128},
		{ID: 99, Offset: 3000, Size: 512},
	}
	// Insert needles.
	for _, n := range needles {
		vol.Add(n)
	}
	// Verify ordering.
	list := vol.List()
	assert(len(list) == 3, "volume should contain three needles")
	assert(list[0].ID == 7 && list[1].ID == 42 && list[2].ID == 99, "needles must be sorted by ID")

	// Retrieve existing needle.
	n, ok := vol.Get(42)
	assert(ok, "needle 42 must exist")
	assert(n.Offset == 1000 && n.Size == 256, "needle 42 fields mismatch")

	// Retrieve non‑existent needle.
	_, ok = vol.Get(123)
	assert(!ok, "needle 123 must not exist")

	// Delete a needle.
	removed := vol.Delete(7)
	assert(removed, "needle 7 should be removed")
	_, ok = vol.Get(7)
	assert(!ok, "needle 7 must no longer exist")
	assert(len(vol.List()) == 2, "volume should now contain two needles")

	// Replace an existing needle.
	vol.Add(Needle{ID: 99, Offset: 4000, Size: 1024})
	n, ok = vol.Get(99)
	assert(ok && n.Offset == 4000 && n.Size == 1024, "needle 99 should be updated")

	// ----- Ring tests -----
	ring := NewRing(3) // three virtual nodes per real node
	ring.AddNode("nodeA")
	ring.AddNode("nodeB")
	ring.AddNode("nodeC")

	// Ensure deterministic mapping for a set of keys.
	keys := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	mapping := make(map[string]string)
	for _, k := range keys {
		node := ring.GetNode(k)
		assert(node != "", fmt.Sprintf("key %s must map to a node", k))
		mapping[k] = node
	}
	// Remove a node and verify that keys previously mapped to it migrate.
	ring.RemoveNode("nodeB")
	for _, k := range keys {
		newNode := ring.GetNode(k)
		assert(newNode != "", fmt.Sprintf("key %s must still map after removal", k))
		// If the original node was nodeB, the new node must differ.
		if mapping[k] == "nodeB" {
			assert(newNode != "nodeB", fmt.Sprintf("key %s should have migrated", k))
		}
	}

	// ----- Simple benchmark (non‑formal) -----
	const ops = 100_000
	start := time.Now()
	for i := 0; i < ops; i++ {
		vol.Add(Needle{ID: uint64(i), Offset: int64(i * 10), Size: 64})
	}
	elapsed := time.Since(start)
	fmt.Printf("Inserted %d needles in %s (≈%0.2f µs/op)\n", ops, elapsed, float64(elapsed.Microseconds())/float64(ops))

	fmt.Println("All assertions passed.")
}

// === END OF main.go ===
