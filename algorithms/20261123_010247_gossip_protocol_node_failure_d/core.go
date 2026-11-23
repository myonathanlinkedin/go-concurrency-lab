package main

import (
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	ID               string
	peers            []*Node
	lastHeartbeat    map[string]time.Time
	mu               sync.Mutex
	failureThreshold time.Duration
	gossipInterval   time.Duration
	checkInterval    time.Duration
	failureCallback  func(nodeID string)
	failedPeers      map[string]bool
	stopCh           chan struct{}
}

func NewNode(id string, peers []*Node, threshold, gossipInterval, checkInterval time.Duration, failureCallback func(nodeID string)) *Node {
	n := &Node{
		ID:               id,
		peers:            peers,
		lastHeartbeat:    make(map[string]time.Time),
		failureThreshold: threshold,
		gossipInterval:   gossipInterval,
		checkInterval:    checkInterval,
		failureCallback:  failureCallback,
		failedPeers:      make(map[string]bool),
		stopCh:           make(chan struct{}),
	}
	now := time.Now()
	for _, p := range peers {
		n.lastHeartbeat[p.ID] = now
	}
	return n
}

func (n *Node) Start() {
	go n.gossipLoop()
	go n.failureLoop()
}

func (n *Node) Stop() {
	close(n.stopCh)
}

func (n *Node) ReceiveHeartbeat(fromID string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.lastHeartbeat[fromID] = time.Now()
}

func (n *Node) gossipLoop() {
	ticker := time.NewTicker(n.gossipInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			n.sendHeartbeat()
		case <-n.stopCh:
			return
		}
	}
}

func (n *Node) sendHeartbeat() {
	n.mu.Lock()
	peers := make([]*Node, 0, len(n.peers))
	for _, p := range n.peers {
		if p.ID != n.ID {
			peers = append(peers, p)
		}
	}
	n.mu.Unlock()
	if len(peers) == 0 {
		return
	}
	target := peers[rand.Intn(len(peers))]
	target.ReceiveHeartbeat(n.ID)
}

func (n *Node) failureLoop() {
	ticker := time.NewTicker(n.checkInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			n.checkFailures()
		case <-n.stopCh:
			return
		}
	}
}

func (n *Node) checkFailures() {
	n.mu.Lock()
	defer n.mu.Unlock()
	now := time.Now()
	for _, p := range n.peers {
		if p.ID == n.ID {
			continue
		}
		last := n.lastHeartbeat[p.ID]
		if now.Sub(last) > n.failureThreshold {
			if !n.failedPeers[p.ID] {
				n.failedPeers[p.ID] = true
				if n.failureCallback != nil {
					n.failureCallback(p.ID)
				}
			}
		}
	}
}
