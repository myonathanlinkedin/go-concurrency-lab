package main

import (
	"container/heap"
	"math"
)

// Coord represents a point in the grid.
type Coord struct {
	X, Y int
}

// CostFunc returns the traversal cost for a given coordinate.
// A negative value indicates an impassable obstacle.
type CostFunc func(x, y int) int

// Node represents a state in the A* search.
type Node struct {
	coord   Coord
	g, h, f int
	parent  *Node
	index   int // index in the heap
}

// NodeHeap implements a min-heap based on node f-values.
type NodeHeap []*Node

func (h NodeHeap) Len() int           { return len(h) }
func (h NodeHeap) Less(i, j int) bool { return h[i].f < h[j].f }
func (h NodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *NodeHeap) Push(x interface{}) {
	n := x.(*Node)
	n.index = len(*h)
	*h = append(*h, n)
}
func (h *NodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	*h = old[:n-1]
	node.index = -1
	return node
}

// Heuristic estimates cost from a to b using Manhattan distance.
func Heuristic(a, b Coord) int {
	return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}

// AStar performs the A* search on a grid defined by costFunc.
// Returns the path from start to goal, the total cost, and a boolean indicating success.
func AStar(start, goal Coord, costFunc CostFunc) ([]Coord, int, bool) {
	if start == goal {
		return []Coord{start}, 0, true
	}

	open := &NodeHeap{}
	heap.Init(open)
	startNode := &Node{coord: start, g: 0, h: Heuristic(start, goal)}
	startNode.f = startNode.g + startNode.h
	heap.Push(open, startNode)

	closed := make(map[Coord]bool)

	for open.Len() > 0 {
		current := heap.Pop(open).(*Node)
		if current.coord == goal {
			// Reconstruct path
			path := []Coord{}
			for n := current; n != nil; n = n.parent {
				path = append([]Coord{n.coord}, path...)
			}
			return path, current.g, true
		}
		closed[current.coord] = true

		neighbors := []Coord{
			{current.coord.X + 1, current.coord.Y},
			{current.coord.X - 1, current.coord.Y},
			{current.coord.X, current.coord.Y + 1},
			{current.coord.X, current.coord.Y - 1},
		}
		for _, nb := range neighbors {
			if closed[nb] {
				continue
			}
			cost := costFunc(nb.X, nb.Y)
			if cost < 0 {
				continue // obstacle
			}
			tentativeG := current.g + cost
			// Check if node is already in open with lower g
			inOpen := false
			for _, n := range *open {
				if n.coord == nb && tentativeG >= n.g {
					inOpen = true
					break
				}
			}
			if inOpen {
				continue
			}
			h := Heuristic(nb, goal)
			f := tentativeG + h
			heap.Push(open, &Node{coord: nb, g: tentativeG, h: h, f: f, parent: current})
		}
	}
	return nil, 0, false
}
