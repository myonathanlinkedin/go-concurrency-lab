package main

import "sort"

// NewVectorClock creates a vector clock with all provided node identifiers initialized to zero.
func NewVectorClock(nodes []string) VectorClock {
	vc := make(VectorClock, len(nodes))
	for _, n := range nodes {
		vc[n] = 0
	}
	return vc
}

// OrderEvents returns a topologically sorted slice of events respecting causal ordering.
// If a cycle is detected (which should not happen with proper vector clocks), the original order is returned.
func OrderEvents(events []Event) []Event {
	// Build adjacency list where edge a -> b means a happens before b.
	n := len(events)
	adj := make([][]int, n)
	inDeg := make([]int, n)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			switch events[i].Clock.Compare(events[j].Clock) {
			case Before:
				adj[i] = append(adj[i], j)
				inDeg[j]++
			}
		}
	}

	// Kahn's algorithm
	queue := make([]int, 0, n)
	for i, d := range inDeg {
		if d == 0 {
			queue = append(queue, i)
		}
	}
	sorted := make([]Event, 0, n)
	for len(queue) > 0 {
		// deterministic selection: smallest ID
		sort.Slice(queue, func(a, b int) bool {
			return events[queue[a]].ID < events[queue[b]].ID
		})
		u := queue[0]
		queue = queue[1:]

		sorted = append(sorted, events[u])
		for _, v := range adj[u] {
			inDeg[v]--
			if inDeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	if len(sorted) != n {
		// Cycle detected; fallback
		return events
	}
	return sorted
}
