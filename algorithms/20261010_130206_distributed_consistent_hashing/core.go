package main

import (
	"hash/fnv"
	"sort"
	"strconv"
)

type ConsistentHash struct {
	ring         []uint64
	nodes        map[uint64]string
	virtualNodes int
	hashFunc     func(string) uint64
}

func NewConsistentHash(virtualNodes int) *ConsistentHash {
	return &ConsistentHash{
		ring:         []uint64{},
		nodes:        make(map[uint64]string),
		virtualNodes: virtualNodes,
		hashFunc: func(s string) uint64 {
			h := fnv.New64a()
			_, _ = h.Write([]byte(s))
			return h.Sum64()
		},
	}
}

func (c *ConsistentHash) AddNode(node string) {
	for i := 0; i < c.virtualNodes; i++ {
		vnodeKey := node + "#" + strconv.Itoa(i)
		hash := c.hashFunc(vnodeKey)
		c.ring = append(c.ring, hash)
		c.nodes[hash] = node
	}
	sort.Slice(c.ring, func(i, j int) bool { return c.ring[i] < c.ring[j] })
}

func (c *ConsistentHash) RemoveNode(node string) {
	toRemove := make(map[uint64]struct{})
	for hash, n := range c.nodes {
		if n == node {
			toRemove[hash] = struct{}{}
		}
	}
	if len(toRemove) == 0 {
		return
	}
	for hash := range toRemove {
		delete(c.nodes, hash)
	}
	newRing := make([]uint64, 0, len(c.ring)-len(toRemove))
	for _, h := range c.ring {
		if _, ok := toRemove[h]; !ok {
			newRing = append(newRing, h)
		}
	}
	c.ring = newRing
}

func (c *ConsistentHash) GetNode(key string) string {
	if len(c.ring) == 0 {
		return ""
	}
	hash := c.hashFunc(key)
	idx := sort.Search(len(c.ring), func(i int) bool { return c.ring[i] >= hash })
	if idx == len(c.ring) {
		idx = 0
	}
	return c.nodes[c.ring[idx]]
}
