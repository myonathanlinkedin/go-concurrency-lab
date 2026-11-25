package main

import (
	"math/bits"
)

// BuddyAllocator implements a simple buddy memory allocation system.
// The allocator manages a contiguous byte slice (the memory pool) and
// provides Allocate and Free operations based on power‑of‑two block sizes.
type BuddyAllocator struct {
	pool      []byte   // underlying memory pool
	minSize   uint64   // minimum block size (must be power of two)
	maxOrder  uint   // highest order (0‑based)
	freeLists [][]uint64 // freeLists[order] holds offsets of free blocks of that order
}

// NewBuddyAllocator creates a new allocator managing a pool of totalSize bytes.
// totalSize must be a power of two and at least minBlockSize.
// minBlockSize must also be a power of two.
func NewBuddyAllocator(totalSize uint64, minBlockSize uint64) *BuddyAllocator {
	if !isPowerOfTwo(totalSize) || !isPowerOfTwo(minBlockSize) {
		panic("sizes must be powers of two")
	}
	if totalSize < minBlockSize {
		panic("totalSize must be >= minBlockSize")
	}
	orders := log2(totalSize/minBlockSize) + 1 // include order 0
	freeLists := make([][]uint64, orders)
	// initially the whole pool is a single free block at the highest order
	freeLists[orders-1] = []uint64{0}
	return &BuddyAllocator{
		pool:      make([]byte, totalSize),
		minSize:   minBlockSize,
		maxOrder:  orders - 1,
		freeLists: freeLists,
	}
}

// Allocate returns the offset within the pool of a block that can hold at least size bytes.
// The returned offset is aligned to the block size (a power of two).
// If allocation fails, returns ^uint64(0) (max uint64) as sentinel.
func (a *BuddyAllocator) Allocate(size uint64) uint64 {
	if size == 0 {
		panic("cannot allocate zero bytes")
	}
	// Determine the smallest order that can satisfy the request.
	blockSize := roundUpToPowerOfTwo(size)
	if blockSize < a.minSize {
		blockSize = a.minSize
	}
	order := a.sizeToOrder(blockSize)

	// Find a free block at order or higher.
	foundOrder := uint(0xFFFFFFFF)
	for o := order; o <= a.maxOrder; o++ {
		if len(a.freeLists[o]) > 0 {
			foundOrder = o
			break
		}
	}
	if foundOrder == 0xFFFFFFFF {
		// No suitable block.
		return ^uint64(0)
	}

	// Remove block from foundOrder list.
	offset := a.freeLists[foundOrder][len(a.freeLists[foundOrder])-1]
	a.freeLists[foundOrder] = a.freeLists[foundOrder][:len(a.freeLists[foundOrder])-1]

	// Split down to desired order.
	for o := foundOrder; o > order; o-- {
		o-- // split from o+1 to o
		halfSize := a.minSize << o
		// The two buddies: left at current offset, right at offset+halfSize
		rightBuddy := offset + halfSize
		// Push right buddy onto free list of order o.
		a.freeLists[o] = append(a.freeLists[o], rightBuddy)
		// Continue with left buddy (offset unchanged) for next iteration.
	}
	return offset
}

// Free releases a previously allocated block.
// The caller must provide the original allocation size (the request size, not the rounded block size).
func (a *BuddyAllocator) Free(offset uint64, requestSize uint64) {
	if offset >= uint64(len(a.pool)) {
		panic("invalid offset")
	}
	if requestSize == 0 {
		panic("cannot free zero bytes")
	}
	// Determine the order of the block that was actually allocated.
	blockSize := roundUpToPowerOfTwo(requestSize)
	if blockSize < a.minSize {
		blockSize = a.minSize
	}
	order := a.sizeToOrder(blockSize)

	// Coalesce with buddy if possible.
	for {
		buddyOffset := a.buddyOffset(offset, order)
		// Search for buddy in free list of current order.
		idx := indexOf(a.freeLists[order], buddyOffset)
		if idx == -1 {
			// Buddy not free; stop coalescing.
			break
		}
		// Remove buddy from free list.
		lastIdx := len(a.freeLists[order]) - 1
		a.freeLists[order][idx] = a.freeLists[order][lastIdx]
		a.freeLists[order] = a.freeLists[order][:lastIdx]

		// Determine new combined block offset (the lower of the two).
		if buddyOffset < offset {
			offset = buddyOffset
		}
		// Move to next higher order.
		if order == a.maxOrder {
			break
		}
		order++
	}
	// Insert the (possibly coalesced) block into its free list.
	a.freeLists[order] = append(a.freeLists[order], offset)
}

// Helper: compute order from block size (size is guaranteed to be power of two and >= minSize).
func (a *BuddyAllocator) sizeToOrder(size uint64) uint {
	return uint(bits.Len64(size/a.minSize) - 1)
}

// Helper: compute buddy offset for a block at given order.
func (a *BuddyAllocator) buddyOffset(offset uint64, order uint) uint64 {
	blockSize := a.minSize << order
	return offset ^ blockSize
}

// Utility: round up to the next power of two (returns 1 for size==0).
func roundUpToPowerOfTwo(v uint64) uint64 {
	if v == 0 {
		return 1
	}
	if isPowerOfTwo(v) {
		return v
	}
	return 1 << (bits.Len64(v))
}

// Utility: check if a number is a power of two.
func isPowerOfTwo(v uint64) bool {
	return v != 0 && (v&(v-1)) == 0
}

// Utility: integer log2 of a power‑of‑two value.
func log2(v uint64) uint {
	if v == 0 {
		panic("log2 of zero")
	}
	return uint(bits.Len64(v) - 1)
}

// Utility: find index of value in slice, -1 if not found.
func indexOf(slice []uint64, val uint64) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}
