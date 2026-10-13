package main

import (
	"fmt"
	"math/rand"
)

const failureThreshold uint64 = 3 // rounds without update before suspecting failure

// Engine orchestrates the simulation of nodes and gossip rounds.
type Engine struct {
	Nodes map[NodeID]*Node
	rng   *rand.Rand
}

// NewEngine constructs an Engine with a deterministic random source.
func NewEngine(seed int64) *Engine {
	return &Engine{
		Nodes: make(map[NodeID]*Node),
		rng:   rand.New(rand.NewSource(seed)),
	}
}

// AddNode registers a new node in the engine.
func (e *Engine) AddNode(id NodeID) {
	if _, exists := e.Nodes[id]; exists {
		panic(fmt.Sprintf("node %s already exists", id))
	}
	e.Nodes[id] = NewNode(id)
}

// DeactivateNode stops a node from participating in future rounds (simulates crash).
func (e *Engine) DeactivateNode(id NodeID) {
	if n, ok := e.Nodes[id]; ok {
		n.Active = false
	}
}

// SimulateRound performs one gossip round:
// 1. Active nodes increment their own counter.
// 2. Each active node gossips its view to a deterministic subset of peers (here: all peers).
// 3. After gossip, each node updates failure suspicion based on the failureThreshold.
func (e *Engine) SimulateRound() {
	// Phase 1: increment counters for active nodes.
	for _, n := range e.Nodes {
		if n.Active {
			n.Counter++
			// Node always knows its own latest counter.
			n.LastSeen[n.ID] = n.Counter
		}
	}

	// Phase 2: gossip exchange.
	// Collect messages first to avoid order dependencies.
	messages := make([]GossipMessage, 0)
	for _, sender := range e.Nodes {
		if !sender.Active {
			continue
		}
		// Prepare a shallow copy of LastSeen for transmission.
		seenCopy := make(map[NodeID]uint64, len(sender.LastSeen))
		for pid, cnt := range sender.LastSeen {
			seenCopy[pid] = cnt
		}
		msg := GossipMessage{
			From:    sender.ID,
			Counter: sender.Counter,
			Seen:    seenCopy,
		}
		messages = append(messages, msg)
	}

	// Deliver each message to all other nodes (including the sender for completeness).
	for _, msg := range messages {
		for _, receiver := range e.Nodes {
			// Even a failed (inactive) node does not process incoming gossip.
			if !receiver.Active {
				continue
			}
			// Update receiver's view of the sender.
			if prev, ok := receiver.LastSeen[msg.From]; !ok || msg.Counter > prev {
				receiver.LastSeen[msg.From] = msg.Counter
			}
			// Merge the sender's knowledge of other peers.
			for pid, cnt := range msg.Seen {
				if pid == receiver.ID {
					// Receiver already knows its own counter via the direct update above.
					continue
				}
				if prev, ok := receiver.LastSeen[pid]; !ok || cnt > prev {
					receiver.LastSeen[pid] = cnt
				}
			}
		}
	}

	// Phase 3: failure detection.
	for _, n := range e.Nodes {
		if !n.Active {
			// Inactive nodes are considered failed by everyone else; they do not evaluate.
			continue
		}
		for pid, last := range n.LastSeen {
			if pid == n.ID {
				continue
			}
			// If the difference exceeds the threshold, mark as failed.
			if n.Counter-last > failureThreshold {
				n.Failed[pid] = true
			}
		}
	}
}

// GetFailedPeers returns a slice of peer IDs that the given node currently suspects as failed.
func (e *Engine) GetFailedPeers(nodeID NodeID) []NodeID {
	n, ok := e.Nodes[nodeID]
	if !ok {
		panic(fmt.Sprintf("node %s not found", nodeID))
	}
	failed := make([]NodeID, 0)
	for pid := range n.Failed {
		failed = append(failed, pid)
	}
	return failed
}
