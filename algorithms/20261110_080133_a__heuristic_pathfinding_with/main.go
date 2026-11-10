package main

import (
	"fmt"
)

func main() {
	// Test 1: Simple path with obstacles
	gridCost := func(x, y int) int {
		// 5x5 grid, obstacles at (1,1),(1,2),(1,3)
		if x < 0 || x > 4 || y < 0 || y > 4 {
			return -1 // outside bounds treated as obstacle
		}
		if x == 1 && (y == 1 || y == 2 || y == 3) {
			return -1
		}
		return 1 // normal cost
	}
	start := Coord{0, 0}
	goal := Coord{4, 4}
	path, cost, found := AStar(start, goal, gridCost)
	if !found {
		panic("Test 1 failed: path not found")
	}
	if cost != 8 {
		panic(fmt.Sprintf("Test 1 failed: expected cost 8, got %d", cost))
	}
	if len(path) != 9 {
		panic(fmt.Sprintf("Test 1 failed: expected path length 9, got %d", len(path)))
	}

	// Test 2: Unreachable goal
	unreachCost := func(x, y int) int {
		if x < 0 || x > 2 || y < 0 || y > 2 {
			return -1
		}
		if x == 1 && y == 1 {
			return -1
		}
		return 1
	}
	start2 := Coord{0, 0}
	goal2 := Coord{2, 2}
	_, _, found2 := AStar(start2, goal2, unreachCost)
	if found2 {
		panic("Test 2 failed: path should not be found")
	}

	// Test 3: Dynamic cost change
	dynamicCost := func(x, y int) int {
		if x == 2 && y == 2 {
			return 5 // high cost
		}
		return 1
	}
	start3 := Coord{0, 0}
	goal3 := Coord{4, 4}
	path3, cost3, found3 := AStar(start3, goal3, dynamicCost)
	if !found3 {
		panic("Test 3 failed: path not found")
	}
	if cost3 != 8 {
		panic(fmt.Sprintf("Test 3 failed: expected cost 8, got %d", cost3))
	}
	if len(path3) != 9 {
		panic(fmt.Sprintf("Test 3 failed: expected path length 9, got %d", len(path3)))
	}

	// Test 4: Start equals goal
	start4 := Coord{2, 2}
	goal4 := Coord{2, 2}
	path4, cost4, found4 := AStar(start4, goal4, dynamicCost)
	if !found4 {
		panic("Test 4 failed: path not found")
	}
	if cost4 != 0 {
		panic(fmt.Sprintf("Test 4 failed: expected cost 0, got %d", cost4))
	}
	if len(path4) != 1 {
		panic(fmt.Sprintf("Test 4 failed: expected path length 1, got %d", len(path4)))
	}

	fmt.Println("All tests passed.")
}
