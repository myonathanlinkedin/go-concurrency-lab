package main

import "fmt"

// HashFunc defines a function that maps a byte slice to a uint64 hash.
type HashFunc func(data []byte) uint64

// HashRing implements a consistent hashing ring with virtual nodes.
type HashRing struct {
	// numVirtual is the number of virtual nodes per real node.
	numVirtual int
	// hashFn is the hash function used to map keys and nodes to the ring.
	hashFn HashFunc
	// ring holds the sorted hash values of all virtual nodes.
	ring []uint64
	// vnodeMap maps a virtual node hash to its real node identifier.
	vnodeMap map[uint64]string
}

// NewHashRing creates a new HashRing with the specified number of virtual nodes per real node.
// If hashFn is nil, a default FNV-1a 64-bit hash function is used.
func NewHashRing(numVirtual int, hashFn HashFunc) *HashRing {
	if numVirtual <= 0 {
		panic(fmt.Sprintf("numVirtual must be > 0, got %d", numVirtual))
	}
	if hashFn == nil {
		hashFn = defaultHash
	}
	return &HashRing{
		numVirtual: numVirtual,
		hashFn:     hashFn,
		ring:       []uint64{},
		vnodeMap:   make(map[uint64]string),
	}
}
