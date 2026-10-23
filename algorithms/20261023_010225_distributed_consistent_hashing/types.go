package main

import "sync"

type ConsistentHash struct {
    ring         []uint64
    hashMap      map[uint64]string
    virtualNodes int
    mutex        sync.RWMutex
}

func NewConsistentHash(virtualNodes int) *ConsistentHash {
    return &ConsistentHash{
        ring:         make([]uint64, 0),
        hashMap:      make(map[uint64]string),
        virtualNodes: virtualNodes,
    }
}
