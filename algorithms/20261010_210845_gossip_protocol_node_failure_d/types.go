package main

import "fmt"

type NodeID int

type GossipMessage struct {
	From    NodeID
	Counter uint64
}

// Node represents a participant in the gossip-based failure detector.
type Node struct {
	ID          NodeID
	Counter     uint64            // Local heartbeat counter.
	Peers       []NodeID          // IDs of nodes to which this node can send gossip.
	Received    map[NodeID]uint64 // Latest known counters from peers.
	LastUpdate  map[NodeID]int    // Step index when the last gossip was received from each peer.
	Suspect     map[NodeID]bool   // Current suspicion status for each peer.
	Active      bool              // Whether the node is still participating (used to simulate failure).
}

// NewNode creates a new Node with the given identifier and peer list.
func NewNode(id NodeID, peers []NodeID) *Node {
	return &Node{
		ID:         id,
		Counter:    0,
		Peers:      peers,
		Received:   make(map[NodeID]uint64),
		LastUpdate: make(map[NodeID]int),
		Suspect:    make(map[NodeID]bool),
		Active:     true,
	}
}

// Tick increments the node's local heartbeat counter and produces a gossip message.
// The step argument is the current simulation step (used only for debugging).
func (n *Node) Tick(step int) GossipMessage {
	if !n.Active {
		// Inactive nodes do not generate messages.
		return GossipMessage{}
	}
	n.Counter++
	// Debug output (optional, can be removed for production).
	_ = fmt.Sprintf("Step %d: Node %d ticked, counter=%d", step, n.ID, n.Counter)
	return GossipMessage{From: n.ID, Counter: n.Counter}
}

// Receive processes an incoming gossip message, updating the node's view of the sender.
func (n *Node) Receive(msg GossipMessage, step int) {
	if !n.Active {
		return
	}
	// Record the latest counter seen from the sender.
	if prev, ok := n.Received[msg.From]; !ok || msg.Counter > prev {
		n.Received[msg.From] = msg.Counter
		n.LastUpdate[msg.From] = step
		// Debug output.
		_ = fmt.Sprintf("Step %d: Node %d received counter %d from Node %d", step, n.ID, msg.Counter, msg.From)
	}
}

// Detect updates the suspicion map based on the elapsed steps since the last update.
// If a peer has not been heard from for more than maxMissedSteps, it is suspected.
func (n *Node) Detect(step int, maxMissedSteps int) {
	if !n.Active {
		return
	}
	for _, peer := range n.Peers {
		last, ok := n.LastUpdate[peer]
		if !ok {
			// Never heard from this peer yet.
			n.Suspect[peer] = true
			continue
		}
		if step-last > maxMissedSteps {
			n.Suspect[peer] = true
		} else {
			n.Suspect[peer] = false
		}
	}
}

// IsSuspecting returns true if the node currently suspects the given peer.
func (n *Node) IsSuspecting(peer NodeID) bool {
	s, ok := n.Suspect[peer]
	if !ok {
		return false
	}
	return s
}
