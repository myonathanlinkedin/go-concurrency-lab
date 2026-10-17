package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// Deterministic random source.
	rng := rand.New(rand.NewSource(42))

	// Configuration.
	const nodeCount = 5
	const maxSteps = 12
	const failureStep = 4          // Step at which node 2 stops being active.
	const maxMissedSteps = 2       // Threshold for suspicion.

	// Initialize nodes with full mesh topology (each node knows all others).
	nodes := make(map[NodeID]*Node)
	var ids []NodeID
	for i := 0; i < nodeCount; i++ {
		id := NodeID(i)
		ids = append(ids, id)
	}
	for _, id := range ids {
		// Peers are all other IDs.
		var peers []NodeID
		for _, pid := range ids {
			if pid != id {
				peers = append(peers, pid)
			}
		}
		nodes[id] = NewNode(id, peers)
	}

	// Simulation loop.
	for step := 0; step < maxSteps; step++ {
		// Simulate failure of node 2 at the designated step.
		if step == failureStep {
			if node, ok := nodes[NodeID(2)]; ok {
				node.Active = false
			}
		}
		GossipStep(nodes, step, rng, maxMissedSteps)
	}

	// Assertions.

	// After the simulation, all alive nodes (0,1,3,4) should suspect node 2.
	for _, id := range []NodeID{0, 1, 3, 4} {
		if !nodes[id].IsSuspecting(NodeID(2)) {
			panic(fmt.Sprintf("Assertion failed: Node %d does not suspect failed Node 2", id))
		}
	}

	// Alive nodes should NOT suspect each other.
	for _, a := range []NodeID{0, 1, 3, 4} {
		for _, b := range []NodeID{0, 1, 3, 4} {
			if a == b {
				continue
			}
			if nodes[a].IsSuspecting(b) {
				panic(fmt.Sprintf("Assertion failed: Node %d incorrectly suspects alive Node %d", a, b))
			}
		}
	}

	// Edge case: single node system should never suspect itself.
	singleNode := NewNode(NodeID(99), []NodeID{})
	singleNode.Tick(0)
	singleNode.Detect(0, maxMissedSteps)
	if singleNode.IsSuspecting(NodeID(99)) {
		panic("Assertion failed: Single node suspects itself")
	}

	// Edge case: empty system (no nodes) – nothing to do, just ensure no panic.
	emptyNodes := make(map[NodeID]*Node)
	GossipStep(emptyNodes, 0, rng, maxMissedSteps)

	fmt.Println("All assertions passed.")
}
