package main

import (
	"encoding/binary"
	"hash/fnv"
)

// CMS implements a Count-Min Sketch for frequency estimation.
type CMS struct {
	width  int
	depth  int
	table  [][]uint64
	seeds  []uint64
}

// NewCMS creates a new Count-Min Sketch with the specified width and depth.
func NewCMS(width, depth int) *CMS {
	if width <= 0 || depth <= 0 {
		panic("width and depth must be positive")
	}
	table := make([][]uint64, depth)
	for i := 0; i < depth; i++ {
		table[i] = make([]uint64, width)
	}
	seeds := make([]uint64, depth)
	for i := 0; i < depth; i++ {
		seeds[i] = uint64(i + 1)
	}
	return &CMS{
		width: width,
		depth: depth,
		table: table,
		seeds: seeds,
	}
}

// hash computes a 64-bit hash of item using the provided seed.
func (c *CMS) hash(item, seed uint64) uint64 {
	h := fnv.New64a()
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], seed)
	h.Write(buf[:])
	binary.BigEndian.PutUint64(buf[:], item)
	h.Write(buf[:])
	return h.Sum64()
}

// Update increments the estimated count of item by count.
func (c *CMS) Update(item, count uint64) {
	if count == 0 {
		return
	}
	for i := 0; i < c.depth; i++ {
		idx := int(c.hash(item, c.seeds[i]) % uint64(c.width))
		c.table[i][idx] += count
	}
}

// Estimate returns the estimated count of item.
func (c *CMS) Estimate(item uint64) uint64 {
	min := uint64(^uint64(0)) // max uint64
	for i := 0; i < c.depth; i++ {
		idx := int(c.hash(item, c.seeds[i]) % uint64(c.width))
		if c.table[i][idx] < min {
			min = c.table[i][idx]
		}
	}
	return min
}

// HeavyHitters returns items from the provided slice whose estimated count
// is at least threshold.
func (c *CMS) HeavyHitters(items []uint64, threshold uint64) []uint64 {
	var result []uint64
	for _, item := range items {
		if c.Estimate(item) >= threshold {
			result = append(result, item)
		}
	}
	return result
}
