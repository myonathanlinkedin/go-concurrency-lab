package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	const virtualNodes = 10
	ch := NewConsistentHash(virtualNodes)

	// Test 1: Empty ring
	if ch.GetNode("foo") != "" {
		panic("Test 1 failed: expected empty string for empty ring")
	}

	// Test 2: Add nodes and distribution
	ch.AddNode("node1")
	ch.AddNode("node2")
	ch.AddNode("node3")

	keyCount := 1000
	dist := map[string]int{}
	for i := 0; i < keyCount; i++ {
		key := fmt.Sprintf("key-%d", i)
		node := ch.GetNode(key)
		if node == "" {
			panic("Test 2 failed: got empty node")
		}
		dist[node]++
	}
	if len(dist) != 3 {
		panic("Test 2 failed: expected 3 nodes in distribution")
	}
	for node, cnt := range dist {
		if cnt == 0 {
			panic(fmt.Sprintf("Test 2 failed: node %s has zero keys", node))
		}
	}

	// Test 3: Consistency of mapping
	sampleKey := "consistent-key"
	nodeA := ch.GetNode(sampleKey)
	nodeB := ch.GetNode(sampleKey)
	if nodeA != nodeB {
		panic("Test 3 failed: inconsistent mapping for same key")
	}

	// Test 4: Remove a node
	ch.RemoveNode("node2")
	for i := 0; i < keyCount; i++ {
		key := fmt.Sprintf("key-%d", i)
		node := ch.GetNode(key)
		if node == "node2" {
			panic("Test 4 failed: key mapped to removed node")
		}
	}

	// Test 5: Add node back
	ch.AddNode("node2")
	for i := 0; i < keyCount; i++ {
		key := fmt.Sprintf("key-%d", i)
		node := ch.GetNode(key)
		if node == "" {
			panic("Test 5 failed: got empty node after re-adding")
		}
	}

	// Test 6: Ring sorted
	for i := 1; i < len(ch.ring); i++ {
		if ch.ring[i-1] > ch.ring[i] {
			panic("Test 6 failed: ring not sorted")
		}
	}

	// Test 7: Random keys over time
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("rand-%d", rand.Int63())
		node := ch.GetNode(key)
		if node == "" {
			panic("Test 7 failed: empty node for random key")
		}
	}

	fmt.Println("All tests passed successfully.")
}
