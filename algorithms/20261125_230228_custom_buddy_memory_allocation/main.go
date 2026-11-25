package main

import (
	"fmt"
)

func main() {
	// Create a buddy allocator with 1024 bytes total and minimum block size 32 bytes.
	const totalSize uint64 = 1024
	const minBlock uint64 = 32
	alloc := NewBuddyAllocator(totalSize, minBlock)

	// Helper to assert conditions.
	assert := func(cond bool, msg string) {
		if !cond {
			panic("assertion failed: " + msg)
		}
	}

	// 1. Simple allocation and free.
	off1 := alloc.Allocate(50) // should round up to 64
	assert(off1 != ^uint64(0), "allocation 1 failed")
	alloc.Free(off1, 50)

	// After freeing, the allocator should have a single free block of max order.
	assert(len(alloc.freeLists[alloc.maxOrder]) == 1, "expected one block at max order after free")
	assert(alloc.freeLists[alloc.maxOrder][0] == 0, "expected block offset 0 after free")

	// 2. Allocate multiple blocks to cause fragmentation.
	offA := alloc.Allocate(100) // 128
	offB := alloc.Allocate(200) // 256
	offC := alloc.Allocate(60)  // 64
	assert(offA != ^uint64(0) && offB != ^uint64(0) && offC != ^uint64(0), "multiple allocations failed")

	// Free middle block and ensure it appears in appropriate free list.
	alloc.Free(offB, 200) // frees a 256‑byte block
	order256 := alloc.sizeToOrder(256)
	assert(indexOf(alloc.freeLists[order256], offB) != -1, "freed 256‑byte block not in correct free list")

	// Free remaining blocks and verify full coalescence.
	alloc.Free(offA, 100) // frees 128‑byte block
	alloc.Free(offC, 60)  // frees 64‑byte block

	// After all frees, we should be back to a single block covering the whole pool.
	assert(len(alloc.freeLists[alloc.maxOrder]) == 1, "expected single max‑order block after all frees")
	assert(alloc.freeLists[alloc.maxOrder][0] == 0, "expected offset 0 after full coalescence")

	// 3. Exhaustive allocation until failure.
	var offsets []uint64
	for {
		off := alloc.Allocate(32) // smallest block size
		if off == ^uint64(0) {
			break
		}
		offsets = append(offsets, off)
	}
	// All 32‑byte blocks should have been allocated: totalSize / 32 = 32 blocks.
	assert(len(offsets) == int(totalSize/minBlock), "unexpected number of smallest blocks allocated")

	// Free them all.
	for i, off := range offsets {
		alloc.Free(off, 32)
		_ = i // silence unused variable warning
	}
	// Verify again single max‑order block.
	assert(len(alloc.freeLists[alloc.maxOrder]) == 1 && alloc.freeLists[alloc.maxOrder][0] == 0,
		"allocator not fully restored after exhaustive free")

	fmt.Println("All buddy allocator tests passed.")
}
