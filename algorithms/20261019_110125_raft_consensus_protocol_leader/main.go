package main

import (
	"fmt"
	"time"
)

func assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: " + msg)
	}
}

// simulateCluster runs a deterministic simulation of a Raft cluster.
func simulateCluster(nodeCount int, maxTicks int) []*Node {
	// Create nodes and wire peers.
	nodes := make([]*Node, nodeCount)
	for i := 0; i < nodeCount; i++ {
		nodes[i] = NewNode(i, nil)
	}
	// Assign peers (excluding self).
	for i := 0; i < nodeCount; i++ {
		peers := make([]*Node, 0, nodeCount-1)
		for j := 0; j < nodeCount; j++ {
			if i != j {
				peers = append(peers, nodes[j])
			}
		}
		nodes[i].Peers = peers
	}
	// Run ticks.
	for tick := 0; tick < maxTicks; tick++ {
		for _, n := range nodes {
			n.Tick()
		}
	}
	return nodes
}

// countLeaders returns the number of nodes that consider themselves leader.
func countLeaders(nodes []*Node) (count int, term int, leaderIDs []int) {
	term = -1
	for _, n := range nodes {
		if n.State == Leader {
			count++
			leaderIDs = append(leaderIDs, n.ID)
			if term == -1 {
				term = n.CurrentTerm
			} else {
				// All leaders must share the same term.
				assert(n.CurrentTerm == term, fmt.Sprintf("leader term mismatch: %d vs %d", term, n.CurrentTerm))
			}
		}
	}
	return
}

func main() {
	// Test 1: Basic leader election in a 3-node cluster.
	nodes := simulateCluster(3, 1000)
	leaderCount, term, leaderIDs := countLeaders(nodes)
	assert(leaderCount == 1, fmt.Sprintf("expected 1 leader, got %d", leaderCount))
	assert(term >= 1, fmt.Sprintf("expected term >=1, got %d", term))
	fmt.Printf("Test1 passed: Leader %v elected in term %d\n", leaderIDs, term)

	// Test 2: Leader steps down when a higher term appears.
	// Force node 0 to become leader with term 1.
	leader := nodes[0]
	leader.becomeLeader()
	leader.CurrentTerm = 1

	// Simulate node 1 receiving a higher term heartbeat.
	hb := AppendEntriesArgs{Term: 5, LeaderID: 1}
	_ = nodes[1].HandleAppendEntries(hb)

	// After processing, node 0 should step down to follower with term 5.
	assert(leader.State == Follower, "leader should have stepped down")
	assert(leader.CurrentTerm == 5, fmt.Sprintf("leader term should be 5, got %d", leader.CurrentTerm))
	fmt.Println("Test2 passed: Leader stepped down on higher term")

	// Test 3: No split brain – after many ticks, still exactly one leader.
	nodes = simulateCluster(5, 2000)
	leaderCount, term, leaderIDs = countLeaders(nodes)
	assert(leaderCount == 1, fmt.Sprintf("expected 1 leader in 5-node cluster, got %d", leaderCount))
	assert(term >= 1, fmt.Sprintf("expected term >=1, got %d", term))
	fmt.Printf("Test3 passed: Single leader %v in term %d for 5-node cluster\n", leaderIDs, term)

	// Simple benchmark: measure ticks per second for 7-node cluster.
	start := time.Now()
	_ = simulateCluster(7, 5000)
	elapsed := time.Since(start)
	fmt.Printf("Benchmark: 7-node cluster processed 5000 ticks in %v (%.2f ticks/sec)\n",
		elapsed, float64(5000)/elapsed.Seconds())
}
