package main

import "errors"

type Mode int

const (
	ModeCP Mode = iota // Consistency over Availability
	ModeAP             // Availability over Consistency
	ModeCA             // Consistency & Availability (no partitions allowed)
)

type Node struct {
	id    int
	store map[string]string
	alive bool
}

type Cluster struct {
	nodes       []*Node
	mode        Mode
	partitioned bool
}

// Errors returned by Cluster operations
var (
	ErrUnavailable = errors.New("operation unavailable due to partition")
	ErrNotFound    = errors.New("key not found")
	ErrInvalidMode = errors.New("invalid mode for operation")
)
