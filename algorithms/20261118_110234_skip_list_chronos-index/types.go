package main

import (
	"math/rand"
	"sync"
)

// Node represents a single element in the skip list.
type Node struct {
	key     int64
	value   interface{}
	forward []*Node
}

// SkipList is a probabilistic ordered data structure supporting O(log n) search, insert, and delete.
type SkipList struct {
	header    *Node
	maxLevel  int
	level     int
	prob      float64
	randSrc   *rand.Rand
	mutex     sync.RWMutex
}

// NewSkipList creates a new empty skip list with the given maximum level and probability.
// maxLevel must be >= 1; prob must be in (0,1).
func NewSkipList(maxLevel int, prob float64) *SkipList {
	if maxLevel < 1 {
		maxLevel = 1
	}
	if prob <= 0 || prob >= 1 {
		prob = 0.5
	}
	// Header node with nil value and forward slice sized to maxLevel.
	header := &Node{
		forward: make([]*Node, maxLevel),
	}
	// Deterministic source for reproducibility.
	src := rand.New(rand.NewSource(42))
	return &SkipList{
		header:   header,
		maxLevel: maxLevel,
		level:    1,
		prob:     prob,
		randSrc:  src,
	}
}
