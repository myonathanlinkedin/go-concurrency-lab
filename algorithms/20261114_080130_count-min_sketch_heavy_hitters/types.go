package main

import (
	"hash/fnv"
	"math"
)

// CountMinSketch is a probabilistic data structure for frequency estimation.
type CountMinSketch struct {
	depth   uint
	width   uint
	table   [][]uint64
	hashes  []uint64 // per-row seeds for hash diversification
	total   uint64   // total count of all updates
}

// NewCountMinSketch creates a CountMinSketch with error guarantees.
// epsilon: additive error factor (e.g., 0.01 for 1% error)
// delta: probability of error exceeding epsilon (e.g., 0.01 for 99% confidence)
func NewCountMinSketch(epsilon, delta float64) *CountMinSketch {
	if epsilon <= 0 || delta <= 0 || epsilon >= 1 || delta >= 1 {
		panic("epsilon and delta must be in (0,1)")
	}
	w := uint(math.Ceil(math.E / epsilon))
	d := uint(math.Ceil(math.Log(1 / delta)))
	// Initialize 2D counter table
	table := make([][]uint64, d)
	for i := range table {
		table[i] = make([]uint64, w)
	}
	// Predefined seeds for each row (using large odd constants)
	const seedBase uint64 = 0x9e3779b97f4a7c15
	hashes := make([]uint64, d)
	for i := uint(0); i < d; i++ {
		hashes[i] = seedBase ^ uint64(i*0x9e3779b97f4a7c15)
	}
	return &CountMinSketch{
		depth:  d,
		width:  w,
		table:  table,
		hashes: hashes,
	}
}

// hash computes a 64‑bit FNV‑1a hash of data combined with a row‑specific seed.
func (cms *CountMinSketch) hash(data []byte, seed uint64) uint64 {
	h := fnv.New64a()
	// Write seed first (as 8 bytes, little endian)
	var b [8]byte
	for i := 0; i < 8; i++ {
		b[i] = byte(seed >> (8 * uint(i)))
	}
	_, _ = h.Write(b[:])
	_, _ = h.Write(data)
	return h.Sum64()
}

// Add increments the count of item by delta (default 1 if delta==0).
func (cms *CountMinSketch) Add(item []byte, delta uint64) {
	if delta == 0 {
		delta = 1
	}
	for i := uint(0); i < cms.depth; i++ {
		h := cms.hash(item, cms.hashes[i])
		idx := uint(h % uint64(cms.width))
		cms.table[i][idx] += delta
	}
	cms.total += delta
}

// Estimate returns the estimated frequency of item (upper bound).
func (cms *CountMinSketch) Estimate(item []byte) uint64 {
	min := uint64(math.MaxUint64)
	for i := uint(0); i < cms.depth; i++ {
		h := cms.hash(item, cms.hashes[i])
		idx := uint(h % uint64(cms.width))
		if cms.table[i][idx] < min {
			min = cms.table[i][idx]
		}
	}
	return min
}

// HeavyHittersTracker maintains a map of observed items to enable extraction of heavy hitters.
type HeavyHittersTracker struct {
	sketch    *CountMinSketch
	frequencies map[string]uint64
}

// NewHeavyHittersTracker creates a tracker with given sketch parameters.
func NewHeavyHittersTracker(epsilon, delta float64) *HeavyHittersTracker {
	return &HeavyHittersTracker{
		sketch:      NewCountMinSketch(epsilon, delta),
		frequencies: make(map[string]uint64),
	}
}

// Add records an occurrence of item (as string) with optional delta.
func (hht *HeavyHittersTracker) Add(item string, delta uint64) {
	hht.sketch.Add([]byte(item), delta)
	hht.frequencies[item] += delta
}

// Estimate returns the sketch's frequency estimate for item.
func (hht *HeavyHittersTracker) Estimate(item string) uint64 {
	return hht.sketch.Estimate([]byte(item))
}

// HeavyHitters returns all items whose estimated frequency exceeds threshold.
func (hht *HeavyHittersTracker) HeavyHitters(threshold uint64) []string {
	var result []string
	for item := range hht.frequencies {
		if hht.Estimate(item) >= threshold {
			result = append(result, item)
		}
	}
	return result
}
