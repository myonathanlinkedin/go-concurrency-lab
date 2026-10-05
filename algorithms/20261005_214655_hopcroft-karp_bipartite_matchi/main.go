package main

import (
	"fmt"
)

func assert(condition bool, msg string) {
	if !condition {
		panic("assertion failed: " + msg)
	}
}

func testEmptyGraph() {
	g := NewGraph(0, 0)
	res := HopcroftKarp(g)
	assert(res.Size == 0, "empty graph should have size 0")
}

func testSingleEdge() {
	g := NewGraph(1, 1)
	g.AddEdge(0, 0)
	res := HopcroftKarp(g)
	assert(res.Size == 1, "single edge graph should have size 1")
	assert(res.LeftMatch[0] == 0 && res.RightMatch[0] == 0, "matching should pair vertices 0-0")
}

func testMultipleComponents() {
	// Component 1: perfect matching of size 2
	// Component 2: one extra left vertex with no edges
	g := NewGraph(5, 4)
	// Component 1
	g.AddEdge(0, 0)
	g.AddEdge(0, 1)
	g.AddEdge(1, 0)
	g.AddEdge(1, 1)
	// Component 2
	g.AddEdge(2, 2)
	g.AddEdge(3, 3)
	// Vertex 4 is isolated on left side
	res := HopcroftKarp(g)
	assert(res.Size == 4, fmt.Sprintf("expected matching size 4, got %d", res.Size))
	// Verify that each matched left vertex indeed points to a distinct right vertex
	used := make(map[int]bool)
	for u, v := range res.LeftMatch {
		if v != -1 {
			assert(!used[v], fmt.Sprintf("right vertex %d matched twice", v))
			used[v] = true
			assert(res.RightMatch[v] == u, "inconsistent reverse mapping")
		}
	}
}

func testDenseGraph() {
	n := 100
	g := NewGraph(n, n)
	// Complete bipartite graph K_{n,n}
	for u := 0; u < n; u++ {
		for v := 0; v < n; v++ {
			g.AddEdge(u, v)
		}
	}
	res := HopcroftKarp(g)
	assert(res.Size == n, fmt.Sprintf("complete graph should have matching size %d", n))
}

func main() {
	fmt.Println("Running Hopcroft‑Karp unit tests...")
	testEmptyGraph()
	testSingleEdge()
	testMultipleComponents()
	testDenseGraph()
	fmt.Println("All tests passed.")
}
