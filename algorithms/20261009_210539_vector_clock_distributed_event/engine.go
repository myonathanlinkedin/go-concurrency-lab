package main

import "fmt"

// String provides a human‑readable representation of the vector clock.
func (vc *VectorClock) String() string {
	// Produce a deterministic ordering for readability.
	type kv struct {
		pid int
		ts  int
	}
	var entries []kv
	for pid, ts := range vc.clock {
		entries = append(entries, kv{pid, ts})
	}
	// Simple insertion sort (vector clocks are tiny in typical tests).
	for i := 1; i < len(entries); i++ {
		j := i
		for j > 0 && entries[j-1].pid > entries[j].pid {
			entries[j-1], entries[j] = entries[j], entries[j-1]
			j--
		}
	}
	result := "{"
	for i, e := range entries {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%d:%d", e.pid, e.ts)
	}
	result += "}"
	return result
}
