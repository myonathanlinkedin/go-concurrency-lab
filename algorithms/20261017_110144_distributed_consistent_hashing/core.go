package main

import (
	"hash/fnv"
	"sort"
	"strconv"
)

type ConsistentHashRouter struct {
	virtualNodes int
	hashRing     []uint64
	nodeMap      map[uint64]string
}

func NewConsistentHashRouter(virtualNodes int) *ConsistentHashRouter {
	return &ConsistentHashRouter{
		virtualNodes: virtualNodes,
		hashRing:     make([]uint64, 0),
		nodeMap:      make(map[uint64]string),
	}
}

func (r *ConsistentHashRouter) hashKey(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	return h.Sum64()
}

func (r *ConsistentHashRouter) AddNode(nodeID string) {
	for i := 0; i < r.virtualNodes; i++ {
		vnodeKey := nodeID + "#" + strconv.Itoa(i)
		hash := r.hashKey(vnodeKey)
		r.hashRing = append(r.hashRing, hash)
		r.nodeMap[hash] = nodeID
	}
	sort.Slice(r.hashRing, func(i, j int) bool { return r.hashRing[i] < r.hashRing[j] })
}

func (r *ConsistentHashRouter) RemoveNode(nodeID string) {
	newRing := make([]uint64, 0, len(r.hashRing))
	for _, h := range r.hashRing {
		if r.nodeMap[h] != nodeID {
			newRing = append(newRing, h)
		} else {
			delete(r.nodeMap, h)
		}
	}
	r.hashRing = newRing
}

func (r *ConsistentHashRouter) GetNode(key string) string {
	if len(r.hashRing) == 0 {
		return ""
	}
	hash := r.hashKey(key)
	idx := sort.Search(len(r.hashRing), func(i int) bool { return r.hashRing[i] >= hash })
	if idx == len(r.hashRing) {
		idx = 0
	}
	return r.nodeMap[r.hashRing[idx]]
}
