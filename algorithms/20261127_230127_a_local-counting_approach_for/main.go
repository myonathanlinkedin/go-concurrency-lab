package main

import (
	"fmt"
)

func main() {
	// Test 1: Simple path 1-2-3-4 (induced path of length 3 exists)
	g1 := NewGraph()
	g1.AddEdge(1, 2)
	g1.AddEdge(2, 3)
	g1.AddEdge(3, 4)
	if !g1.HasInducedPath(3) {
		panic("Test 1 failed: expected induced path of length 3")
	}
	if !g1.HasInducedPath(2) {
		panic("Test 1 failed: expected induced path of length 2")
	}
	if !g1.HasInducedPath(1) {
		panic("Test 1 failed: expected induced path of length 1")
	}
	// Test 2: Triangle (3-cycle) – no induced path of length 2
	g2 := NewGraph()
	g2.AddEdge(1, 2)
	g2.AddEdge(2, 3)
	g2.AddEdge(3, 1)
	if g2.HasInducedPath(2) {
		panic("Test 2 failed: triangle should not have induced path of length 2")
	}
	if !g2.HasInducedPath(1) {
		panic("Test 2 failed: triangle should have induced path of length 1 (any edge)")
	}
	// Test 3: Square with a diagonal (vertices 1-2-3-4-1, diagonal 1-3)
	g3 := NewGraph()
	g3.AddEdge(1, 2)
	g3.AddEdge(2, 3)
	g3.AddEdge(3, 4)
	g3.AddEdge(4, 1)
	g3.AddEdge(1, 3) // diagonal
	// Induced path of length 2 exists: 2-3-4 (2 and 4 are not adjacent)
	if !g3.HasInducedPath(2) {
		panic("Test 3 failed: expected induced path of length 2")
	}
	// Induced path of length 3 does NOT exist because any 4‑vertex sequence contains a chord.
	if g3.HasInducedPath(3) {
		panic("Test 3 failed: should not have induced path of length 3")
	}
	// Test 4: Disconnected graph – no edges
	g4 := NewGraph()
	if g4.HasInducedPath(1) {
		panic("Test 4 failed: empty graph should not have induced path of length 1")
	}
	// Test 5: Star graph (center 0 connected to 1,2,3,4). Induced path of length 2 exists: leaf‑center‑leaf where leaves are not adjacent.
	g5 := NewGraph()
	for i := 1; i <= 4; i++ {
		g5.AddEdge(0, i)
	}
	if !g5.HasInducedPath(2) {
		panic("Test 5 failed: star should have induced path of length 2")
	}
	if g5.HasInducedPath(3) {
		panic("Test 5 failed: star should not have induced path of length 3")
	}
	// Summary
	fmt.Println("All assertions passed.")
}
