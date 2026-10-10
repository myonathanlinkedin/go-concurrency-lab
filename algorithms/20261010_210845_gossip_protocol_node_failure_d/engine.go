package main

import (
	"math/rand"
)

// GossipStep performs one simulation step:
// 1. Each active node ticks and creates a gossip message.
// 2. Each node selects a random peer (if any) and sends its message.
// 3. Receivers process the messages.
// 4. Each node runs its failure detection logic.
func GossipStep(nodes map[NodeID]*Node, step int, rng *rand.Rand, maxMissedSteps int) {
	// Collect messages to be delivered this step.
	type delivery struct {
		to   NodeID
		msg  GossipMessage
	}
	var deliveries []delivery

	// Phase 1: each node ticks and prepares a message.
	for _, node := range nodes {
		if !node.Active {
			continue
		}
		msg := node.Tick(step)
		if len(node.Peers) == 0 {
			continue
		}
		// Choose a random peer to gossip to.
		peerIdx := rng.Intn(len(node.Peers))
		peerID := node.Peers[peerIdx]
		deliveries = append(deliveries, delivery{to: peerID, msg: msg})
	}

	// Phase 2: deliver messages.
	for _, d := range deliveries {
		receiver, ok := nodes[d.to]
		if ok {
			receiver.Receive(d.msg, step)
		}
	}

	// Phase 3: each node runs detection.
	for _, node := range nodes {
		node.Detect(step, maxMissedSteps)
	}
}
