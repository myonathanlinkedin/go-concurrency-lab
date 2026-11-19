package main

import (
    "fmt"
    "math"
)

func assertFloatEqual(a, b, eps float64) {
    if math.Abs(a-b) > eps {
        panic(fmt.Sprintf("assertion failed: %v != %v (eps %v)", a, b, eps))
    }
}

func assertIntSliceEqual(a, b []int) {
    if len(a) != len(b) {
        panic(fmt.Sprintf("assertion failed: slice lengths differ %d != %d", len(a), len(b)))
    }
    for i := range a {
        if a[i] != b[i] {
            panic(fmt.Sprintf("assertion failed: slice elements differ at %d: %d != %d", i, a[i], b[i]))
        }
    }
}

func bruteForceTSP(cities []City) (float64, []int) {
    n := len(cities)
    if n == 0 {
        return 0, nil
    }
    if n == 1 {
        return 0, []int{0, 0}
    }
    best := math.Inf(1)
    bestPath := []int{}
    perm := make([]int, n-1)
    for i := 0; i < n-1; i++ {
        perm[i] = i + 1
    }
    for {
        path := make([]int, n+1)
        path[0] = 0
        copy(path[1:], perm)
        path[n] = 0
        dist := 0.0
        for i := 0; i < n; i++ {
            a := path[i]
            b := path[i+1]
            dx := cities[a].X - cities[b].X
            dy := cities[a].Y - cities[b].Y
            dist += math.Sqrt(dx*dx + dy*dy)
        }
        if dist < best {
            best = dist
            bestPath = append([]int(nil), path...)
        }
        // next permutation
        i := n - 2
        for i >= 0 && perm[i] > perm[i+1] {
            i--
        }
        if i < 0 {
            break
        }
        j := n - 2
        for perm[j] < perm[i] {
            j--
        }
        perm[i], perm[j] = perm[j], perm[i]
        for l, r := i+1, n-2; l < r; l, r = l+1, r-1 {
            perm[l], perm[r] = perm[r], perm[l]
        }
    }
    return best, bestPath
}

func main() {
    // Test 1: square
    square := []City{
        {0, 0},
        {0, 1},
        {1, 1},
        {1, 0},
    }
    tsp1 := NewTSP(square)
    best1, path1 := tsp1.HeldKarp()
    assertFloatEqual(best1, 4.0, 1e-9)
    if len(path1) != len(square)+1 {
        panic("path length mismatch")
    }
    if path1[0] != 0 || path1[len(path1)-1] != 0 {
        panic("path does not start/end at 0")
    }

    // Test 2: triangle
    triangle := []City{
        {0, 0},
        {1, 0},
        {0, 1},
    }
    tsp2 := NewTSP(triangle)
    best2, _ := tsp2.HeldKarp()
    expected2 := 1.0 + math.Sqrt(2) + 1.0
    assertFloatEqual(best2, expected2, 1e-9)

    // Test 3: single city
    single := []City{{0, 0}}
    tsp3 := NewTSP(single)
    best3, path3 := tsp3.HeldKarp()
    assertFloatEqual(best3, 0.0, 1e-9)
    if len(path3) != 2 || path3[0] != 0 || path3[1] != 0 {
        panic("single city path incorrect")
    }

    // Test 4: random 5 cities vs brute force
    random := []City{
        {0, 0},
        {2, 3},
        {5, 1},
        {1, 4},
        {3, 2},
    }
    tsp4 := NewTSP(random)
    best4, path4 := tsp4.HeldKarp()
    bfBest, _ := bruteForceTSP(random)
    assertFloatEqual(best4, bfBest, 1e-9)

    // verify path4 is a valid Hamiltonian cycle
    visited := make([]bool, len(random))
    for _, v := range path4 {
        if v < 0 || v >= len(random) {
            panic("path contains invalid city index")
        }
        if visited[v] && v != 0 {
            panic("city visited twice")
        }
        visited[v] = true
    }
    for i := 1; i < len(random); i++ {
        if !visited[i] {
            panic("city not visited")
        }
    }

    fmt.Println("All tests passed")
}
