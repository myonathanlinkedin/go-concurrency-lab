package main

// Graph represents a bipartite graph with vertices partitioned into
// a left set {0..nLeft-1} and a right set {0..nRight-1}.
// Edges are stored as adjacency lists from left vertices to right vertices.
type Graph struct {
	nLeft  int
	nRight int
	adj    [][]int // adj[u] = slice of right vertices v adjacent to left vertex u
}

// NewGraph creates a new bipartite graph with the given number of left and right vertices.
func NewGraph(nLeft, nRight int) *Graph {
	adj := make([][]int, nLeft)
	for i := range adj {
		adj[i] = []int{}
	}
	return &Graph{
		nLeft:  nLeft,
		nRight: nRight,
		adj:    adj,
	}
}

// AddEdge adds an undirected edge between left vertex u and right vertex v.
// u must be in [0, nLeft) and v in [0, nRight). Panics on out‑of‑range indices.
func (g *Graph) AddEdge(u, v int) {
	if u < 0 || u >= g.nLeft {
		panic("AddEdge: left vertex out of range")
	}
	if v < 0 || v >= g.nRight {
		panic("AddEdge: right vertex out of range")
	}
	g.adj[u] = append(g.adj[u], v)
}

// MatchingResult holds the outcome of a maximum bipartite matching.
type MatchingResult struct {
	Size      int   // number of matched pairs
	LeftMatch []int // LeftMatch[u] = matched right vertex or -1
	RightMatch []int // RightMatch[v] = matched left vertex or -1
}
