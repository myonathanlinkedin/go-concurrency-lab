package main

import (
	"math/rand"
	"time"
)

func NewEngine(threshold int) *Engine {
	return &Engine{
		nodes:              make(map[string]*Node),
		suspicionThreshold: threshold,
		rand:               rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (e *Engine) AddNode(id string) {
	e.nodes[id] = &Node{
		ID:        id,
		known:     make(map[string]int),
		suspected: make(map[string]bool),
	}
}

func (e *Engine) RemoveNode(id string) {
	delete(e.nodes, id)
}

func (e *Engine) RunRounds(n int) {
	for i := 0; i < n; i++ {
		e.runRound()
	}
}

func (e *Engine) runRound() {
	e.round++
	// Snapshot of node IDs to avoid mutation during iteration
	ids := make([]string, 0, len(e.nodes))
	for id := range e.nodes {
		ids = append(ids, id)
	}
	// Each node updates its own known timestamp
	for _, id := range ids {
		node := e.nodes[id]
		node.known[id] = e.round
	}
	// Gossip phase
	for _, id := range ids {
		node := e.nodes[id]
		if len(ids) <= 1 {
			continue
		}
		peerID := e.selectPeer(id, ids)
		peer := e.nodes[peerID]
		e.gossip(node, peer)
	}
	// Update suspicion lists
	for _, id := range ids {
		node := e.nodes[id]
		for knownID, lastRound := range node.known {
			if knownID == node.ID {
				continue
			}
			if e.round-lastRound > e.suspicionThreshold {
				node.suspected[knownID] = true
			} else {
				delete(node.suspected, knownID)
			}
		}
	}
}

func (e *Engine) selectPeer(excludeID string, ids []string) string {
	var peers []string
	for _, id := range ids {
		if id != excludeID {
			peers = append(peers, id)
		}
	}
	return peers[e.rand.Intn(len(peers))]
}

func (e *Engine) gossip(node, peer *Node) {
	// Merge peer's known into node's known
	for id, round := range peer.known {
		if existing, ok := node.known[id]; !ok || round > existing {
			node.known[id] = round
		}
	}
	// Merge node's known into peer's known
	for id, round := range node.known {
		if existing, ok := peer.known[id]; !ok || round > existing {
			peer.known[id] = round
		}
	}
}
