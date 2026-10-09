package main

import (
	"fmt"
)

func main() {
	// Create a consistent hash ring with 100 virtual nodes per real node.
	ring := NewHashRing(100, nil)

	// Edge case: GetNode on empty ring should return empty string.
	if got := ring.GetNode("anykey"); got != "" {
		panic(fmt.Sprintf("expected empty string for empty ring, got %s", got))
	}

	// Add three nodes.
	ring.AddNode("NodeA")
	ring.AddNode("NodeB")
	ring.AddNode("NodeC")

	// Verify that the same key always maps to the same node.
	keys := []string{
		"user:1001", "user:1002", "order:5001", "session:abcd1234",
	}
	mapping := make(map[string]string)
	for _, k := range keys {
		node := ring.GetNode(k)
		if node == "" {
			panic(fmt.Sprintf("expected non-empty node for key %s", k))
		}
		mapping[k] = node
	}

	// Re-query and assert stability.
	for _, k := range keys {
		if node := ring.GetNode(k); node != mapping[k] {
			panic(fmt.Sprintf("unstable mapping for key %s: got %s, want %s", k, node, mapping[k]))
		}
	}

	// Add a new node and ensure only a subset of keys change owners.
	ring.AddNode("NodeD")
	changed := 0
	for _, k := range keys {
		newNode := ring.GetNode(k)
		if newNode != mapping[k] {
			changed++
		}
	}
	// With consistent hashing, at most (1/numNodes) of keys should move.
	if changed > len(keys)/2 {
		panic(fmt.Sprintf("too many keys remapped after adding NodeD: %d of %d", changed, len(keys)))
	}

	// Remove a node and verify that keys previously owned by it are reassigned.
	ownerBeforeRemoval := ring.GetNode("user:1001")
	ring.RemoveNode(ownerBeforeRemoval)
	newOwner := ring.GetNode("user:1001")
	if newOwner == "" || newOwner == ownerBeforeRemoval {
		panic(fmt.Sprintf("expected key to be reassigned after removal, got %s", newOwner))
	}

	// Final sanity print.
	fmt.Println("All assertions passed. Consistent hashing router behaves as expected.")
}
