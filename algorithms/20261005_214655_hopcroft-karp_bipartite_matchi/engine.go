package main

// bfs builds distance layers from all free left vertices.
// Returns true if there exists at least one augmenting path.
func bfs(g *Graph, pairU, pairV, dist []int) bool {
	queue := make([]int, 0, g.nLeft)
	const INF = int(^uint(0) >> 1) // max int

	for u := 0; u < g.nLeft; u++ {
		if pairU[u] == -1 {
			dist[u] = 0
			queue = append(queue, u)
		} else {
			dist[u] = INF
		}
	}
	found := false

	for head := 0; head < len(queue); head++ {
		u := queue[head]
		for _, v := range g.adj[u] {
			if pairV[v] != -1 && dist[pairV[v]] == INF {
				dist[pairV[v]] = dist[u] + 1
				queue = append(queue, pairV[v])
			}
			if pairV[v] == -1 {
				found = true // we can potentially reach a free right vertex
			}
		}
	}
	return found
}

// dfs searches augmenting paths from left vertex u using the layering built by bfs.
func dfs(g *Graph, u int, pairU, pairV, dist []int) bool {
	const INF = int(^uint(0) >> 1)

	for _, v := range g.adj[u] {
		pu := pairV[v]
		if pu == -1 || (dist[pu] == dist[u]+1 && dfs(g, pu, pairU, pairV, dist)) {
			pairU[u] = v
			pairV[v] = u
			return true
		}
	}
	dist[u] = INF // prune dead ends
	return false
}

// HopcroftKarp computes a maximum cardinality matching on a bipartite graph.
// It returns a MatchingResult containing the size and the match arrays.
func HopcroftKarp(g *Graph) MatchingResult {
	pairU := make([]int, g.nLeft)  // match for left vertices, -1 if free
	pairV := make([]int, g.nRight) // match for right vertices, -1 if free
	dist := make([]int, g.nLeft)

	for i := range pairU {
		pairU[i] = -1
	}
	for i := range pairV {
		pairV[i] = -1
	}

	matching := 0
	for bfs(g, pairU, pairV, dist) {
		for u := 0; u < g.nLeft; u++ {
			if pairU[u] == -1 {
				if dfs(g, u, pairU, pairV, dist) {
					matching++
				}
			}
		}
	}
	return MatchingResult{
		Size:       matching,
		LeftMatch:  pairU,
		RightMatch: pairV,
	}
}
