package main

import (
	"fmt"
	"math"
)

func assertEqualSlice(got, want []int64) {
	if len(got) != len(want) {
		panic(fmt.Sprintf("length mismatch: got %d, want %d", len(got), len(want)))
	}
	for i := range got {
		if got[i] != want[i] {
			panic(fmt.Sprintf("at index %d: got %d, want %d", i, got[i], want[i]))
		}
	}
}

// testSimple verifies Dijkstra on a small directed graph.
func testSimple() {
	// Graph:
	// 0 -> 1 (4), 0 -> 2 (1)
	// 2 -> 1 (2), 2 -> 3 (5)
	// 1 -> 3 (1)
	// 3 -> 4 (3)
	g := Graph{
		0: {{to: 1, weight: 4}, {to: 2, weight: 1}},
		1: {{to: 3, weight: 1}},
		2: {{to: 1, weight: 2}, {to: 3, weight: 5}},
		3: {{to: 4, weight: 3}},
		4: {},
	}
	dist := Dijkstra(g, 0)
	want := []int64{0, 3, 1, 4, 7}
	assertEqualSlice(dist, want)
}

// testDisconnected verifies handling of unreachable vertices.
func testDisconnected() {
	// Two components: 0-1 and 2 isolated.
	g := Graph{
		0: {{to: 1, weight: 2}},
		1: {{to: 0, weight: 2}},
		2: {},
	}
	dist := Dijkstra(g, 0)
	want := []int64{0, 2, math.MaxInt64}
	assertEqualSlice(dist, want)
}

// testZeroWeight verifies zero-weight edges.
func testZeroWeight() {
	g := Graph{
		0: {{to: 1, weight: 0}, {to: 2, weight: 5}},
		1: {{to: 2, weight: 0}},
		2: {},
	}
	dist := Dijkstra(g, 0)
	want := []int64{0, 0, 0}
	assertEqualSlice(dist, want)
}

// benchmark runs Dijkstra on a moderately sized random graph multiple times.
func benchmark() {
	const n = 200
	const iterations = 20
	// Build a dense graph where each vertex connects to the next vertex with weight 1.
	g := make(Graph, n)
	for i := 0; i < n-1; i++ {
		g[i] = []Edge{{to: i + 1, weight: 1}}
	}
	g[n-1] = []Edge{}
	for i := 0; i < iterations; i++ {
		_ = Dijkstra(g, 0)
	}
}

// main runs all tests and the benchmark.
func main() {
	testSimple()
	testDisconnected()
	testZeroWeight()
	benchmark()
	fmt.Println("All tests passed")
}
