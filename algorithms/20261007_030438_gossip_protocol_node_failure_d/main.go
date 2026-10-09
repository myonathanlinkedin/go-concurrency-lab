package main

import (
	"fmt"
)

func main() {
	engine := NewEngine(3)

	// Add five nodes
	engine.AddNode("A")
	engine.AddNode("B")
	engine.AddNode("C")
	engine.AddNode("D")
	engine.AddNode("E")

	// Run 2 rounds; no failures yet
	engine.RunRounds(2)

	// Assert no node suspects any other
	for _, node := range engine.nodes {
		if len(node.suspected) != 0 {
			panic(fmt.Sprintf("Node %s unexpectedly suspects nodes: %v", node.ID, node.suspected))
		}
	}

	// Simulate failure of node C
	engine.RemoveNode("C")

	// Run 5 rounds; C should be suspected by all remaining nodes
	engine.RunRounds(5)

	for _, node := range engine.nodes {
		if !node.suspected["C"] {
			panic(fmt.Sprintf("Node %s did not suspect failed node C", node.ID))
		}
	}

	// Verify that remaining nodes still know each other
	for _, node := range engine.nodes {
		for _, peer := range engine.nodes {
			if node.ID == peer.ID {
				continue
			}
			if _, ok := node.known[peer.ID]; !ok {
				panic(fmt.Sprintf("Node %s missing knowledge of node %s", node.ID, peer.ID))
			}
		}
	}

	fmt.Println("All assertions passed. Gossip failure detector simulation succeeded.")
}
