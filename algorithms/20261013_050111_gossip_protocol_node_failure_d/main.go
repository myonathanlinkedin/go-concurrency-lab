package main

import (
	"fmt"
)

func main() {
	// Initialize engine with a fixed seed for deterministic behavior.
	engine := NewEngine(42)

	// Create three nodes: A, B, C.
	engine.AddNode("A")
	engine.AddNode("B")
	engine.AddNode("C")

	// Helper to assert conditions; panics on failure.
	assert := func(cond bool, msg string) {
		if !cond {
			panic("assertion failed: " + msg)
		}
	}

	// Simulate two healthy rounds.
	for i := 0; i < 2; i++ {
		engine.SimulateRound()
	}

	// Verify that no node suspects any peer yet.
	for _, id := range []NodeID{"A", "B", "C"} {
		failed := engine.GetFailedPeers(id)
		assert(len(failed) == 0, fmt.Sprintf("node %s should have no failures after healthy rounds", id))
	}

	// Deactivate node C to simulate a crash.
	engine.DeactivateNode("C")

	// Simulate enough rounds for failure detection (threshold + 1).
	for i := 0; i < int(failureThreshold)+1; i++ {
		engine.SimulateRound()
	}

	// After the extra rounds, nodes A and B should suspect C.
	for _, id := range []NodeID{"A", "B"} {
		failed := engine.GetFailedPeers(id)
		assert(len(failed) == 1 && failed[0] == "C", fmt.Sprintf("node %s should suspect C as failed", id))
	}

	// Node C is inactive; it should not suspect anyone.
	failedC := engine.GetFailedPeers("C")
	assert(len(failedC) == 0, "inactive node C should not have failure suspicions")

	// Demonstrate that after a node is marked failed, it remains failed even if it later becomes active again.
	// Reactivate C (simulating a node that recovers after being suspected).
	engine.Nodes["C"].Active = true
	// Run a single round to allow C to gossip its revived state.
	engine.SimulateRound()

	// Nodes A and B should still consider C failed because they have not received a fresh heartbeat within the threshold.
	for _, id := range []NodeID{"A", "B"} {
		failed := engine.GetFailedPeers(id)
		assert(len(failed) == 1 && failed[0] == "C", fmt.Sprintf("node %s should still suspect C after recovery", id))
	}

	fmt.Println("All assertions passed. Gossip failure detector behaves as expected.")
}
