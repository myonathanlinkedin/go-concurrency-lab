package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	router := NewConsistentHashRouter(10)

	// Test 1: No nodes
	if node := router.GetNode("anykey"); node != "" {
		panic("expected empty node when no nodes present")
	}

	// Add nodes
	router.AddNode("nodeA")
	router.AddNode("nodeB")
	router.AddNode("nodeC")

	// Test 2: GetNode returns non-empty
	key := "testkey"
	node := router.GetNode(key)
	if node == "" {
		panic("GetNode returned empty for existing nodes")
	}

	// Test 3: Distribution roughly balanced
	counts := map[string]int{"nodeA": 0, "nodeB": 0, "nodeC": 0}
	for i := 0; i < 10000; i++ {
		k := fmt.Sprintf("key-%d", i)
		n := router.GetNode(k)
		counts[n]++
	}
	for n, c := range counts {
		if c < 2500 || c > 3500 {
			panic(fmt.Sprintf("node %s has %d keys, out of expected range", n, c))
		}
	}

	// Test 4: Remove node and ensure keys re-map
	router.RemoveNode("nodeB")
	for i := 0; i < 1000; i++ {
		k := fmt.Sprintf("key-%d", i)
		n := router.GetNode(k)
		if n == "nodeB" {
			panic("key still mapped to removed node")
		}
	}

	// Test 5: Adding duplicate node increases virtual nodes
	prevRingSize := len(router.hashRing)
	router.AddNode("nodeA")
	if len(router.hashRing) != prevRingSize+router.virtualNodes {
		panic("adding duplicate node did not increase ring size correctly")
	}

	fmt.Println("All tests passed.")
}
