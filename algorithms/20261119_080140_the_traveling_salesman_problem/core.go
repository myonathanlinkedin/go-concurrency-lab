package main

import "math"

type City struct {
    X, Y float64
}

type TSP struct {
    cities []City
    dist   [][]float64
}

func NewTSP(cities []City) *TSP {
    n := len(cities)
    dist := make([][]float64, n)
    for i := 0; i < n; i++ {
        dist[i] = make([]float64, n)
        for j := 0; j < n; j++ {
            dx := cities[i].X - cities[j].X
            dy := cities[i].Y - cities[j].Y
            dist[i][j] = math.Sqrt(dx*dx + dy*dy)
        }
    }
    return &TSP{cities: cities, dist: dist}
}

func (t *TSP) distance(i, j int) float64 {
    return t.dist[i][j]
}

// HeldKarp computes the optimal tour length and one optimal path.
// Path is a slice of city indices starting and ending at 0.
func (t *TSP) HeldKarp() (float64, []int) {
    n := len(t.cities)
    if n == 0 {
        return 0, nil
    }
    if n == 1 {
        return 0, []int{0, 0}
    }
    size := 1 << n
    dp := make([]float64, size*n)
    pred := make([]int, size*n)
    for i := range dp {
        dp[i] = math.Inf(1)
        pred[i] = -1
    }
    // base: from start city 0 to j
    for j := 1; j < n; j++ {
        mask := (1 << 0) | (1 << j)
        idx := mask*n + j
        dp[idx] = t.distance(0, j)
    }
    // iterate over masks
    for mask := 0; mask < size; mask++ {
        if mask&(1<<0) == 0 {
            continue
        }
        for last := 1; last < n; last++ {
            if mask&(1<<last) == 0 {
                continue
            }
            prevMask := mask ^ (1 << last)
            if prevMask == 0 {
                continue
            }
            for k := 1; k < n; k++ {
                if prevMask&(1<<k) == 0 {
                    continue
                }
                val := dp[prevMask*n+k] + t.distance(k, last)
                idx := mask*n + last
                if val < dp[idx] {
                    dp[idx] = val
                    pred[idx] = k
                }
            }
        }
    }
    // close tour
    fullMask := size - 1
    best := math.Inf(1)
    bestLast := -1
    for last := 1; last < n; last++ {
        val := dp[fullMask*n+last] + t.distance(last, 0)
        if val < best {
            best = val
            bestLast = last
        }
    }
    // reconstruct path
    path := make([]int, n+1)
    path[0] = 0
    mask := fullMask
    last := bestLast
    for i := n - 1; i >= 1; i-- {
        path[i] = last
        mask ^= 1 << last
        prev := pred[mask*n+last]
        last = prev
    }
    path[n] = 0
    return best, path
}
