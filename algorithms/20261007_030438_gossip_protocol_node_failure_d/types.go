package main

import (
	"math/rand"
)

type Node struct {
	ID        string
	known     map[string]int // last round heard
	suspected map[string]bool
}

type Engine struct {
	nodes              map[string]*Node
	round              int
	suspicionThreshold int
	rand               *rand.Rand
}
