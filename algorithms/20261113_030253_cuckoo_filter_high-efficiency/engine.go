package main

import (
	"math/rand"
)

// NewCuckooFilter creates a filter with at least the requested number of buckets.
// The actual bucket count is rounded up to the next power of two.
func NewCuckooFilter(requestedBuckets uint64) *CuckooFilter {
	if requestedBuckets == 0 {
		requestedBuckets = 1
	}
	bucketCount := nextPowerOfTwo(requestedBuckets)
	buckets := make([][]uint8, bucketCount)
	for i := range buckets {
		buckets[i] = make([]uint8, 0, bucketSize)
	}
	return &CuckooFilter{
		buckets:    buckets,
		bucketMask: bucketCount - 1,
	}
}

// insertIntoBucket tries to place fp into bucket idx. Returns true on success.
func (cf *CuckooFilter) insertIntoBucket(idx uint64, fp uint8) bool {
	b := cf.buckets[idx]
	if len(b) < bucketSize {
		cf.buckets[idx] = append(b, fp)
		return true
	}
	return false
}

// containsInBucket checks whether fp exists in bucket idx.
func (cf *CuckooFilter) containsInBucket(idx uint64, fp uint8) bool {
	for _, v := range cf.buckets[idx] {
		if v == fp {
			return true
		}
	}
	return false
}

// deleteFromBucket removes fp from bucket idx if present. Returns true on success.
func (cf *CuckooFilter) deleteFromBucket(idx uint64, fp uint8) bool {
	b := cf.buckets[idx]
	for i, v := range b {
		if v == fp {
			// Remove element by swapping with last and truncating.
			b[i] = b[len(b)-1]
			cf.buckets[idx] = b[:len(b)-1]
			return true
		}
	}
	return false
}

// Insert adds an item to the filter. Returns true on success, false if the filter is full.
func (cf *CuckooFilter) Insert(item []byte) bool {
	h := hash64(item)
	fp := fingerprint(h)
	i1 := h & cf.bucketMask
	i2 := altIndex(i1, fp, cf.bucketMask)

	// Try direct insertion.
	if cf.insertIntoBucket(i1, fp) || cf.insertIntoBucket(i2, fp) {
		return true
	}

	// Eviction loop.
	idx := i1
	curFP := fp
	for kick := 0; kick < maxKicks; kick++ {
		// Randomly select a victim entry in the current bucket.
		b := cf.buckets[idx]
		evictPos := rand.Intn(len(b))
		evictedFP := b[evictPos]
		// Replace victim with current fingerprint.
		cf.buckets[idx][evictPos] = curFP
		// Prepare to insert evicted fingerprint into its alternate bucket.
		curFP = evictedFP
		idx = altIndex(idx, curFP, cf.bucketMask)

		if cf.insertIntoBucket(idx, curFP) {
			return true
		}
	}
	// Insertion failed after maxKicks.
	return false
}

// Lookup checks whether an item is possibly in the filter.
func (cf *CuckooFilter) Lookup(item []byte) bool {
	h := hash64(item)
	fp := fingerprint(h)
	i1 := h & cf.bucketMask
	i2 := altIndex(i1, fp, cf.bucketMask)

	return cf.containsInBucket(i1, fp) || cf.containsInBucket(i2, fp)
}

// Delete removes an item from the filter if present. Returns true on success.
func (cf *CuckooFilter) Delete(item []byte) bool {
	h := hash64(item)
	fp := fingerprint(h)
	i1 := h & cf.bucketMask
	i2 := altIndex(i1, fp, cf.bucketMask)

	if cf.deleteFromBucket(i1, fp) {
		return true
	}
	if cf.deleteFromBucket(i2, fp) {
		return true
	}
	return false
}
