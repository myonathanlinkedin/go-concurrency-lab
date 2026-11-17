package main

type Disk struct {
	X, Y, R float64
}

func Overlaps(a, b Disk) bool {
	dx := a.X - b.X
	dy := a.Y - b.Y
	distSq := dx*dx + dy*dy
	sumR := a.R + b.R
	return distSq <= sumR*sumR
}

func IndependentSets(disks []Disk) [][]int {
	n := len(disks)
	var result [][]int
	if n < 3 {
		return result
	}
	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]bool, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if Overlaps(disks[i], disks[j]) {
				adj[i][j] = true
				adj[j][i] = true
			}
		}
	}
	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			if adj[i][j] {
				continue
			}
			for k := j + 1; k < n; k++ {
				if adj[i][k] || adj[j][k] {
					continue
				}
				result = append(result, []int{i, j, k})
			}
		}
	}
	return result
}

func CountIndependentSets(disks []Disk) int {
	return len(IndependentSets(disks))
}
