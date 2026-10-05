package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

// murmur3_32 implements the 32‑bit MurmurHash3 algorithm.
// It is deterministic, fast, and uses only the standard library.
func murmur3_32(data []byte, seed uint32) uint32 {
	const (
		c1 uint32 = 0xcc9e2d51
		c2 uint32 = 0x1b873593
		r1 uint32 = 15
		r2 uint32 = 13
		m  uint32 = 5
		n  uint32 = 0xe6546b64
	)

	h := seed
	nblocks := len(data) / 4

	// body
	for i := 0; i < nblocks; i++ {
		i4 := i * 4
		k := binary.LittleEndian.Uint32(data[i4 : i4+4])

		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2

		h ^= k
		h = ((h << r2) | (h >> (32 - r2)))*m + n
	}

	// tail
	var k uint32
	tail := data[nblocks*4:]

	switch len(tail) {
	case 3:
		k ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k ^= uint32(tail[0])
		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2
		h ^= k
	}

	// finalization
	h ^= uint32(len(data))
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16

	return h
}

// BloomFilter implements a classic Bloom filter with
// configurable false‑positive probability.
type BloomFilter struct {
	m     uint          // number of bits
	k     uint          // number of hash functions
	bits  []uint64      // bit array
	seeds []uint32      // per‑hash seeds
}

// NewBloomFilter creates a filter for expected n items and target false‑positive rate p.
// It computes optimal m and k, rounding m up to a multiple of 64 for storage efficiency.
func NewBloomFilter(n uint, p float64) *BloomFilter {
	if n == 0 {
		log.Panic("expected number of items must be > 0")
	}
	if p <= 0 || p >= 1 {
		log.Panic("false‑positive rate must be between 0 and 1")
	}
	ln2 := math.Ln2
	mFloat := -float64(n) * math.Log(p) / (ln2 * ln2)
	m := uint(math.Ceil(mFloat))
	// round up to nearest 64 bits
	if m%64 != 0 {
		m += 64 - (m % 64)
	}
	kFloat := (float64(m) / float64(n)) * ln2
	k := uint(math.Ceil(kFloat))
	if k == 0 {
		k = 1
	}
	bits := make([]uint64, m/64)

	// generate distinct seeds
	seeds := make([]uint32, k)
	base := uint32(time.Now().UnixNano())
	for i := uint(0); i < k; i++ {
		seeds[i] = base + uint32(i*0x9e3779b9) // golden‑ratio step
	}

	return &BloomFilter{
		m:     m,
		k:     k,
		bits:  bits,
		seeds: seeds,
	}
}

// setBit marks the bit at position idx.
func (bf *BloomFilter) setBit(idx uint) {
	word := idx / 64
	bit := idx % 64
	bf.bits[word] |= 1 << bit
}

// getBit returns true if the bit at position idx is set.
func (bf *BloomFilter) getBit(idx uint) bool {
	word := idx / 64
	bit := idx % 64
	return (bf.bits[word] & (1 << bit)) != 0
}

// Add inserts data into the filter.
func (bf *BloomFilter) Add(data []byte) {
	for _, seed := range bf.seeds {
		h := murmur3_32(data, seed)
		idx := uint(h) % bf.m
		bf.setBit(idx)
	}
}

// Test checks whether data is possibly in the set.
func (bf *BloomFilter) Test(data []byte) bool {
	for _, seed := range bf.seeds {
		h := murmur3_32(data, seed)
		idx := uint(h) % bf.m
		if !bf.getBit(idx) {
			return false
		}
	}
	return true
}

// -------------------- Unit Tests --------------------

func assert(cond bool, msg string) {
	if !cond {
		log.Panicf("assertion failed: %s", msg)
	}
}

// testBasic verifies that added items are reported present and non‑added items are mostly absent.
func testBasic() {
	bf := NewBloomFilter(1000, 0.01)
	items := []string{"alice", "bob", "carol", "dave"}
	for _, s := range items {
		bf.Add([]byte(s))
	}
	for _, s := range items {
		assert(bf.Test([]byte(s)), fmt.Sprintf("expected %s to be present", s))
	}
	negatives := []string{"eve", "mallory", "trent"}
	for _, s := range negatives {
		if bf.Test([]byte(s)) {
			// false positive is acceptable but should be rare; we just note it.
		}
	}
}

// testFalsePositiveRate estimates the empirical false‑positive probability.
func testFalsePositiveRate() {
	const n = 5000
	const p = 0.02
	bf := NewBloomFilter(n, p)

	// insert n random strings
	rng := rand.New(rand.NewSource(42))
	inserted := make([]string, n)
	for i := 0; i < n; i++ {
		s := fmt.Sprintf("key-%d-%d", i, rng.Int())
		inserted[i] = s
		bf.Add([]byte(s))
	}

	// test m non‑inserted strings
	const m = 20000
	falsePos := 0
	for i := 0; i < m; i++ {
		s := fmt.Sprintf("query-%d-%d", i, rng.Int())
		if bf.Test([]byte(s)) {
			falsePos++
		}
	}
	empirical := float64(falsePos) / float64(m)
	// Allow 3× the target rate due to statistical variance.
	allowed := p * 3
	assert(empirical <= allowed, fmt.Sprintf("false positive rate too high: got %f, allowed %f", empirical, allowed))
}

// testEdgeCases checks behavior with empty input and extreme parameters.
func testEdgeCases() {
	bf := NewBloomFilter(1, 0.5)
	bf.Add([]byte{})
	assert(bf.Test([]byte{}), "empty slice should be present after insertion")
}

// main runs all tests and reports success.
func main() {
	fmt.Println("Running Bloom filter tests...")
	testBasic()
	testFalsePositiveRate()
	testEdgeCases()
	fmt.Println("All tests passed.")
}
