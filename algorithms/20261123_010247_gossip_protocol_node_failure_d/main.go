package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	failuresCh := make(chan string, 10)

	failureCallback := func(id string) {
		failuresCh <- id
	}

	threshold := 500 * time.Millisecond
	gossipInt := 100 * time.Millisecond
	checkInt := 200 * time.Millisecond

	var nodeA, nodeB, nodeC *Node

	nodeA = NewNode("A", nil, threshold, gossipInt, checkInt, failureCallback)
	nodeB = NewNode("B", nil, threshold, gossipInt, checkInt, failureCallback)
	nodeC = NewNode("C", nil, threshold, gossipInt, checkInt, failureCallback)

	nodeA.peers = []*Node{nodeB, nodeC}
	nodeB.peers = []*Node{nodeA, nodeC}
	nodeC.peers = []*Node{nodeA, nodeB}

	nodeA.Start()
	nodeB.Start()
	nodeC.Start()

	time.Sleep(1 * time.Second)

	nodeB.Stop()

	time.Sleep(1 * time.Second)

	select {
	case failedID := <-failuresCh:
		if failedID != "B" {
			panic(fmt.Sprintf("expected failure of B, got %s", failedID))
		}
	default:
		panic("expected failure notification for B, but none received")
	}

	select {
	case failedID := <-failuresCh:
		panic(fmt.Sprintf("unexpected failure notification for %s", failedID))
	default:
	}

	nodeA.Stop()
	nodeC.Stop()

	fmt.Println("All assertions passed.")
}
