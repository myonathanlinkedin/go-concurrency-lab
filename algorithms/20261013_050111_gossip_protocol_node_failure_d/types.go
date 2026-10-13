package main

type NodeID string

// Node represents a participant in the gossip failure detector.
type Node struct {
	// Unique identifier of the node.
	ID NodeID
	// Counter is the node's own heartbeat, incremented each active round.
	Counter uint64
	// LastSeen maps peer IDs to the most recent heartbeat counter known for that peer.
	LastSeen map[NodeID]uint64
	// Failed maps peer IDs to a boolean indicating whether the peer is considered failed.
	Failed map[NodeID]bool
	// Active indicates whether the node participates in gossip (simulates a live node).
	Active bool
}

// NewNode creates a fresh Node with initialized maps.
func NewNode(id NodeID) *Node {
	return &Node{
		ID:       id,
		Counter:  0,
		LastSeen: make(map[NodeID]uint64),
		Failed:   make(map[NodeID]bool),
		Active:   true,
	}
}

// GossipMessage carries a node's view of heartbeats to a peer.
type GossipMessage struct {
	From    NodeID
	Counter uint64
	// Snapshot of the sender's LastSeen map (deep copy not required for this deterministic simulation).
	Seen map[NodeID]uint64
}
