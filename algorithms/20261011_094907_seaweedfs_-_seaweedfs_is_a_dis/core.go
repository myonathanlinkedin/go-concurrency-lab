package main

import (
	"fmt"
	"hash/fnv"
	"sort"
	"sync"
)

// ---------- Needle and Volume ----------

// Needle represents a small data chunk (a "needle") stored in a SeaweedFS volume.
type Needle struct {
	ID     uint64 // unique identifier
	Offset int64  // offset within the volume file
	Size   int32  // size of the data payload
}

// Volume stores needles in a sorted slice for O(log n) lookup.
type Volume struct {
	mu      sync.RWMutex
	needles []Needle // sorted by Needle.ID
}

// Add inserts or replaces a needle while keeping the slice sorted.
func (v *Volume) Add(n Needle) {
	v.mu.Lock()
	defer v.mu.Unlock()
	i := sort.Search(len(v.needles), func(i int) bool { return v.needles[i].ID >= n.ID })
	if i < len(v.needles) && v.needles[i].ID == n.ID {
		v.needles[i] = n
		return
	}
	// Insert at position i.
	v.needles = append(v.needles, Needle{}) // grow by one
	copy(v.needles[i+1:], v.needles[i:])
	v.needles[i] = n
}

// Get retrieves a needle by ID. The second return value reports existence.
func (v *Volume) Get(id uint64) (Needle, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	i := sort.Search(len(v.needles), func(i int) bool { return v.needles[i].ID >= id })
	if i < len(v.needles) && v.needles[i].ID == id {
		return v.needles[i], true
	}
	return Needle{}, false
}

// Delete removes a needle by ID. Returns true if the needle existed.
func (v *Volume) Delete(id uint64) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	i := sort.Search(len(v.needles), func(i int) bool { return v.needles[i].ID >= id })
	if i >= len(v.needles) || v.needles[i].ID != id {
		return false
	}
	// Remove element at i.
	copy(v.needles[i:], v.needles[i+1:])
	v.needles = v.needles[:len(v.needles)-1]
	return true
}

// List returns a copy of all needles in ascending ID order.
func (v *Volume) List() []Needle {
	v.mu.RLock()
	defer v.mu.RUnlock()
	out := make([]Needle, len(v.needles))
	copy(out, v.needles)
	return out
}

// ---------- Consistent Hash Ring ----------

// FNV-1a 64‑bit constants (typed to avoid overflow warnings).
const (
	fnvOffset64 uint64 = 0xcbf29ce484222325
	fnvPrime64  uint64 = 0x100000001b3
)

// hashString computes a 64‑bit FNV‑1a hash of a string.
func hashString(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}

// Ring implements a thread‑safe consistent hashing ring with virtual nodes.
type Ring struct {
	mu       sync.RWMutex
	hashes   []uint64          // sorted hash values of virtual nodes
	nodeMap  map[uint64]string // hash -> real node identifier
	replicas int               // number of virtual nodes per real node
}

// NewRing creates a new Ring with the desired number of replicas per node.
func NewRing(replicas int) *Ring {
	if replicas <= 0 {
		replicas = 1
	}
	return &Ring{
		hashes:   []uint64{},
		nodeMap:  make(map[uint64]string),
		replicas: replicas,
	}
}

// AddNode inserts a real node and its virtual replicas into the ring.
func (r *Ring) AddNode(nodeID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := 0; i < r.replicas; i++ {
		virtualKey := fmt.Sprintf("%s#%d", nodeID, i)
		h := hashString(virtualKey)
		r.hashes = append(r.hashes, h)
		r.nodeMap[h] = nodeID
	}
	sort.Slice(r.hashes, func(i, j int) bool { return r.hashes[i] < r.hashes[j] })
}

// RemoveNode deletes a real node and all its virtual replicas.
func (r *Ring) RemoveNode(nodeID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	newHashes := r.hashes[:0]
	for _, h := range r.hashes {
		if r.nodeMap[h] != nodeID {
			newHashes = append(newHashes, h)
		} else {
			delete(r.nodeMap, h)
		}
	}
	r.hashes = newHashes
}

// GetNode returns the identifier of the node responsible for the given key.
func (r *Ring) GetNode(key string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.hashes) == 0 {
		return ""
	}
	h := hashString(key)
	// Locate the first hash >= h; wrap around if necessary.
	idx := sort.Search(len(r.hashes), func(i int) bool { return r.hashes[i] >= h })
	if idx == len(r.hashes) {
		idx = 0
	}
	return r.nodeMap[r.hashes[idx]]
}

// ---------- Utility ----------

// assert panics with a formatted message if condition is false.
func assert(cond bool, msg string) {
	if !cond {
		panic(fmt.Sprintf("assertion failed: %s", msg))
	}
}

// === END OF core.go ===
