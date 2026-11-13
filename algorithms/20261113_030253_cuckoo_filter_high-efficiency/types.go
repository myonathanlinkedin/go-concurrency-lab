package main

import (
	"hash/fnv"
)

const (
	// Fingerprint size in bits (8 bits => 1 byte)
	fingerprintSize uint8 = 8

	// Number of entries per bucket
	bucketSize int = 4

	// Maximum number of evictions during an insert before giving up
	maxKicks int = 500

	// Constant used to hash fingerprint to alternate index
	constAltHash uint64 = 0x5bd1e995
)

// CuckooFilter is a high‑efficiency probabilistic set with support for deletions.
type CuckooFilter struct {
	buckets    [][]uint8 // each bucket holds up to bucketSize fingerprints
	bucketMask uint64    // mask for fast modulo (bucketCount is power of two)
}

// hash64 returns a 64‑bit FNV‑1a hash of the given data.
func hash64(data []byte) uint64 {
	h := fnv.New64a()
	_, _ = h.Write(data)
	return h.Sum64()
}

// fingerprint derives a non‑zero fingerprint from a 64‑bit hash.
func fingerprint(hash uint64) uint8 {
	// Use the lower 8 bits; ensure fingerprint is never zero.
	f := uint8(hash & 0xFF)
	if f == 0 {
		f = 1
	}
	return f
}

// altIndex computes the alternate bucket index given an index and fingerprint.
func altIndex(index uint64, fp uint8, mask uint64) uint64 {
	// Simple mixing: hash the fingerprint with a constant and xor.
	h := uint64(fp) * constAltHash
	return (index ^ h) & mask
}

// nextPowerOfTwo returns the smallest power of two >= n (n > 0).
func nextPowerOfTwo(n uint64) uint64 {
	if n == 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return n
}
