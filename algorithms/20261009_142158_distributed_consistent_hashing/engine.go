package main

import (
	"fmt"

	"hash/fnv"
	"sort"
)

// defaultHash is a deterministic FNV-1a 64-bit hash implementation.
func defaultHash(data []byte) uint64 {
	h := fnv.New64a()
	_, _ = h.Write(data)
	return h.Sum64()
}

// hashKey computes the hash of a string key using the ring's hash function.
func (hr *HashRing) hashKey(key string) uint64 {
	return hr.hashFn([]byte(key))
}

// AddNode inserts a real node identified by nodeID into the ring, creating virtual nodes.
func (hr *HashRing) AddNode(nodeID string) {
	for i := 0; i < hr.numVirtual; i++ {
		// Create a unique identifier for each virtual node.
		vnodeID := fmt.Sprintf("%s#%d", nodeID, i)
		vhash := hr.hashFn([]byte(vnodeID))
		hr.ring = append(hr.ring, vhash)
		hr.vnodeMap[vhash] = nodeID
	}
	// Keep the ring sorted for binary search.
	sort.Slice(hr.ring, func(i, j int) bool { return hr.ring[i] < hr.ring[j] })
}

// RemoveNode deletes a real node and all its virtual nodes from the ring.
func (hr *HashRing) RemoveNode(nodeID string) {
	newRing := hr.ring[:0] // reuse underlying array
	for _, vhash := range hr.ring {
		if hr.vnodeMap[vhash] != nodeID {
			newRing = append(newRing, vhash)
		} else {
			delete(hr.vnodeMap, vhash)
		}
	}
	hr.ring = newRing
	// Ring remains sorted because we only removed elements.
}

// GetNode returns the identifier of the real node responsible for the given key.
// If the ring is empty, an empty string is returned.
func (hr *HashRing) GetNode(key string) string {
	if len(hr.ring) == 0 {
		return ""
	}
	keyHash := hr.hashKey(key)
	// Locate the first virtual node hash >= keyHash using binary search.
	idx := sort.Search(len(hr.ring), func(i int) bool { return hr.ring[i] >= keyHash })
	if idx == len(hr.ring) {
		// Wrap around to the first node.
		idx = 0
	}
	vhash := hr.ring[idx]
	return hr.vnodeMap[vhash]
}
