package main

import "fmt"

const INF int = int(^uint(0) >> 1)

// HopcroftKarp holds the state for the algorithm.
type HopcroftKarp struct {
	leftSize  int
	rightSize int
	adj       [][]int // adjacency from left vertices to right vertices
	pairU     []int   // matched right vertex for each left vertex, -1 if free
	pairV     []int   // matched left vertex for each right vertex, -1 if free
	dist      []int   // distances in BFS
}

// NewHopcroftKarp creates a new instance given the bipartite graph.
func NewHopcroftKarp(leftSize, rightSize int, edges [][2]int) *HopcroftKarp {
	adj := make([][]int, leftSize)
	for _, e := range edges {
		u, v := e[0], e[1]
		if u < 0 || u >= leftSize || v < 0 || v >= rightSize {
			panic(fmt.Sprintf("edge (%d,%d) out of bounds", u, v))
		}
		adj[u] = append(adj[u], v)
	}
	return &HopcroftKarp{
		leftSize:  leftSize,
		rightSize: rightSize,
		adj:       adj,
		pairU:     make([]int, leftSize),
		pairV:     make([]int, rightSize),
		dist:      make([]int, leftSize),
	}
}

// bfs builds layers and returns true if there is an augmenting path.
func (hk *HopcroftKarp) bfs() bool {
	queue := make([]int, 0, hk.leftSize)
	for u := 0; u < hk.leftSize; u++ {
		if hk.pairU[u] == -1 {
			hk.dist[u] = 0
			queue = append(queue, u)
		} else {
			hk.dist[u] = INF
		}
	}
	found := false
	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for _, v := range hk.adj[u] {
			pu := hk.pairV[v]
			if pu != -1 && hk.dist[pu] == INF {
				hk.dist[pu] = hk.dist[u] + 1
				queue = append(queue, pu)
			}
			if pu == -1 {
				found = true
			}
		}
	}
	return found
}

// dfs searches augmenting paths from vertex u.
func (hk *HopcroftKarp) dfs(u int) bool {
	for _, v := range hk.adj[u] {
		pu := hk.pairV[v]
		if pu == -1 || (hk.dist[pu] == hk.dist[u]+1 && hk.dfs(pu)) {
			hk.pairU[u] = v
			hk.pairV[v] = u
			return true
		}
	}
	hk.dist[u] = INF
	return false
}

// MaxMatching computes the size of the maximum matching.
func (hk *HopcroftKarp) MaxMatching() int {
	// initialise all matches as free
	for i := range hk.pairU {
		hk.pairU[i] = -1
	}
	for i := range hk.pairV {
		hk.pairV[i] = -1
	}
	matching := 0
	for hk.bfs() {
		for u := 0; u < hk.leftSize; u++ {
			if hk.pairU[u] == -1 && hk.dfs(u) {
				matching++
			}
		}
	}
	return matching
}

// MaxBipartiteMatching is a convenience wrapper.
func MaxBipartiteMatching(leftSize, rightSize int, edges [][2]int) int {
	hk := NewHopcroftKarp(leftSize, rightSize, edges)
	return hk.MaxMatching()
}
