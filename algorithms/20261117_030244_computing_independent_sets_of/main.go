package main

import (
	"fmt"
)

func main() {
	// Test 1: No disks
	disks1 := []Disk{}
	sets1 := IndependentSets(disks1)
	if len(sets1) != 0 {
		panic("Test1 failed: expected 0 sets")
	}

	// Test 2: Three non-overlapping disks
	disks2 := []Disk{
		{0, 0, 1},
		{5, 0, 1},
		{0, 5, 1},
	}
	sets2 := IndependentSets(disks2)
	if len(sets2) != 1 {
		panic(fmt.Sprintf("Test2 failed: expected 1 set, got %d", len(sets2)))
	}
	if !(sets2[0][0] == 0 && sets2[0][1] == 1 && sets2[0][2] == 2) {
		panic("Test2 failed: wrong indices")
	}

	// Test 3: Three overlapping disks (all overlap)
	disks3 := []Disk{
		{0, 0, 5},
		{3, 0, 5},
		{0, 4, 5},
	}
	sets3 := IndependentSets(disks3)
	if len(sets3) != 0 {
		panic("Test3 failed: expected 0 sets")
	}

	// Test 4: Four disks, all independent triples
	disks4 := []Disk{
		{0, 0, 1},
		{5, 0, 1},
		{0, 5, 1},
		{5, 5, 1},
	}
	sets4 := IndependentSets(disks4)
	if len(sets4) != 4 {
		panic(fmt.Sprintf("Test4 failed: expected 4 sets, got %d", len(sets4)))
	}
	for _, tri := range sets4 {
		if Overlaps(disks4[tri[0]], disks4[tri[1]]) ||
			Overlaps(disks4[tri[0]], disks4[tri[2]]) ||
			Overlaps(disks4[tri[1]], disks4[tri[2]]) {
			panic("Test4 failed: found overlapping triple")
		}
	}

	// Test 5: Larger random set, consistency check
	disks5 := []Disk{
		{0, 0, 1},
		{2, 0, 1},
		{4, 0, 1},
		{1, 3, 1},
		{3, 3, 1},
	}
	sets5 := IndependentSets(disks5)
	expected := 0
	n := len(disks5)
	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			if Overlaps(disks5[i], disks5[j]) {
				continue
			}
			for k := j + 1; k < n; k++ {
				if Overlaps(disks5[i], disks5[k]) || Overlaps(disks5[j], disks5[k]) {
					continue
				}
				expected++
			}
		}
	}
	if len(sets5) != expected {
		panic(fmt.Sprintf("Test5 failed: expected %d sets, got %d", expected, len(sets5)))
	}

	fmt.Println("All tests passed.")
}
