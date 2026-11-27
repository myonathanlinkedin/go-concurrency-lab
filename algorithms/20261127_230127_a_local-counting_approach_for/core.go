package main

import (
	"fmt"
)

// Graph represents an undirected simple graph using adjacency sets.
type Graph struct {
	adj map[int]map[int]struct{}
}

// NewGraph creates a graph with no edges.
func NewGraph() *Graph {
	return &Graph{adj: make(map[int]map[int]struct{})}
}

// AddEdge inserts an undirected edge (u, v). Self‑loops are ignored.
func (g *Graph) AddEdge(u, v int) {
	if u == v {
		return
	}
	if g.adj[u] == nil {
		g.adj[u] = make(map[int]struct{})
	}
	if g.adj[v] == nil {
		g.adj[v] = make(map[int]struct{})
	}
	g.adj[u][v] = struct{}{}
	g.adj[v][u] = struct{}{}
}

// HasInducedPath reports whether the graph contains an induced simple path
// of exactly k edges (i.e., k+1 distinct vertices).
// k must be >= 1; for k == 0 the function returns true (a single vertex).
func (g *Graph) HasInducedPath(k int) bool {
	if k < 0 {
		return false
	}
	if k == 0 {
		// Trivially true: any vertex forms a path of length 0.
		return len(g.adj) > 0
	}
	// Pre‑compute local non‑adjacent neighbor pair counts for pruning.
	localPairs := make(map[int]int)
	for v := range g.adj {
		localPairs[v] = g.localNonAdjPairs(v)
	}
	// Try every vertex as a start of the path.
	for start := range g.adj {
		visited := map[int]struct{}{start: {}}
		if g.dfsInduced(start, start, k, visited, localPairs) {
			return true
		}
	}
	return false
}

// dfsInduced performs a depth‑first search respecting the induced‑path condition.
// remaining is the number of edges still to add to the current partial path.
func (g *Graph) dfsInduced(start, current int, remaining int, visited map[int]struct{}, localPairs map[int]int) bool {
	if remaining == 0 {
		return true
	}
	// Prune: if we need at least two more edges, the current vertex must have
	// at least one pair of non‑adjacent neighbors to continue the induced path.
	if remaining >= 2 && localPairs[current] == 0 {
		return false
	}
	for nb := range g.adj[current] {
		if _, ok := visited[nb]; ok {
			continue
		}
		// Induced condition: nb must not be adjacent to any previously visited
		// vertex except the immediate predecessor (current).
		if !g.isInducedExtension(nb, visited, current) {
			continue
		}
		visited[nb] = struct{}{}
		if g.dfsInduced(start, nb, remaining-1, visited, localPairs) {
			return true
		}
		delete(visited, nb)
	}
	return false
}

// isInducedExtension checks that candidate vertex nb does not connect to any
// vertex in visited other than the allowed neighbor prev.
func (g *Graph) isInducedExtension(nb int, visited map[int]struct{}, prev int) bool {
	for v := range visited {
		if v == prev {
			continue
		}
		if _, adj := g.adj[nb][v]; adj {
			return false
		}
	}
	return true
}

// localNonAdjPairs returns the number of unordered neighbor pairs (u,w) of v
// such that u and w are NOT adjacent. This metric is used for pruning.
func (g *Graph) localNonAdjPairs(v int) int {
	neighbors := make([]int, 0, len(g.adj[v]))
	for nb := range g.adj[v] {
		neighbors = append(neighbors, nb)
	}
	count := 0
	for i := 0; i < len(neighbors); i++ {
		for j := i + 1; j < len(neighbors); j++ {
			u, w := neighbors[i], neighbors[j]
			if _, adj := g.adj[u][w]; !adj {
				count++
			}
		}
	}
	return count
}

// DebugString returns a printable representation of the graph (for diagnostics).
func (g *Graph) DebugString() string {
	s := ""
	for v, nbrs := range g.adj {
		s += fmt.Sprintf("%d:", v)
		for nb := range nbrs {
			s += fmt.Sprintf(" %d", nb)
		}
		s += "\n"
	}
	return s
}
