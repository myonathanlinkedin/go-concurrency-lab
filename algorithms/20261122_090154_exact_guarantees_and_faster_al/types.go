package main

// Direction indicates the monotonic direction of a step.
type Direction int

const (
	DirUp Direction = iota
	DirDown
)

// RollerCoasterSolver encapsulates the input sequence and the allowed
// number of direction changes (k). It provides methods to compute the
// longest k‑rollercoaster subsequence.
type RollerCoasterSolver struct {
	seq []int
	k   int
}

// NewSolver constructs a RollerCoasterSolver. It panics if k is negative.
func NewSolver(seq []int, k int) *RollerCoasterSolver {
	if k < 0 {
		panic("k must be non‑negative")
	}
	return &RollerCoasterSolver{
		seq: seq,
		k:   k,
	}
}
