package main

import (
	"fmt"
	"math/rand"
	"time"
)

func assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: " + msg)
	}
}

func testSimple() {
	left, right := 3, 3
	edges := [][2]int{
		{0, 0},
		{0, 1},
		{1, 1},
		{2, 2},
	}
	size := MaxBipartiteMatching(left, right, edges)
	assert(size == 3, fmt.Sprintf("expected 3, got %d", size))
}

func testNoEdges() {
	left, right := 4, 5
	edges := [][2]int{}
	size := MaxBipartiteMatching(left, right, edges)
	assert(size == 0, fmt.Sprintf("expected 0, got %d", size))
}

func testPartial() {
	left, right := 4, 4
	edges := [][2]int{
		{0, 0},
		{0, 1},
		{1, 0},
		{2, 2},
	}
	size := MaxBipartiteMatching(left, right, edges)
	assert(size == 3, fmt.Sprintf("expected 3, got %d", size))
}

func testRandomFixedSeed() {
	rng := rand.New(rand.NewSource(42))
	left, right := 50, 50
	edges := make([][2]int, 0, 500)
	for i := 0; i < 500; i++ {
		u := rng.Intn(left)
		v := rng.Intn(right)
		edges = append(edges, [2]int{u, v})
	}
	size := MaxBipartiteMatching(left, right, edges)
	// In a random dense graph with 500 edges, matching size should be at least 20.
	assert(size >= 20, fmt.Sprintf("expected >=20, got %d", size))
}

func benchmarkLarge() {
	left, right := 1000, 1000
	edges := make([][2]int, 0, 5000)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 5000; i++ {
		edges = append(edges, [2]int{rng.Intn(left), rng.Intn(right)})
	}
	start := time.Now()
	_ = MaxBipartiteMatching(left, right, edges)
	elapsed := time.Since(start)
	fmt.Printf("Benchmark large graph: %v\n", elapsed)
	_ = elapsed // silence unused warning if fmt removed
}

func main() {
	testSimple()
	testNoEdges()
	testPartial()
	testRandomFixedSeed()
	benchmarkLarge()
	fmt.Println("All tests passed.")
}
