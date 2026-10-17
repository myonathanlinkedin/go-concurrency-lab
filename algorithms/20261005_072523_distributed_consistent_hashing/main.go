package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"sort"
	"strconv"
	"sync"
)

type ConsistentHash struct {
	replicas int
	hashRing []uint32
	hashMap  map[uint32]string
	mu       sync.RWMutex
}

func NewConsistentHash(replicas int) *ConsistentHash {
	return &ConsistentHash{
		replicas: replicas,
		hashMap:  make(map[uint32]string),
	}
}

func (c *ConsistentHash) Add(node string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := 0; i < c.replicas; i++ {
		hash := c.hashKey(node + "#" + strconv.Itoa(i))
		c.hashRing = append(c.hashRing, hash)
		c.hashMap[hash] = node
	}
	sort.Slice(c.hashRing, func(i, j int) bool { return c.hashRing[i] < c.hashRing[j] })
}

func (c *ConsistentHash) Remove(node string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := 0; i < c.replicas; i++ {
		hash := c.hashKey(node + "#" + strconv.Itoa(i))
		idx := sort.Search(len(c.hashRing), func(i int) bool { return c.hashRing[i] >= hash })
		if idx < len(c.hashRing) && c.hashRing[idx] == hash {
			c.hashRing = append(c.hashRing[:idx], c.hashRing[idx+1:]...)
			delete(c.hashMap, hash)
		}
	}
}

func (c *ConsistentHash) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.hashRing) == 0 {
		return ""
	}
	hash := c.hashKey(key)
	idx := sort.Search(len(c.hashRing), func(i int) bool { return c.hashRing[i] >= hash })
	if idx == len(c.hashRing) {
		idx = 0
	}
	return c.hashMap[c.hashRing[idx]]
}

func (c *ConsistentHash) hashKey(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func assert(cond bool, msg string) {
	if !cond {
		fmt.Println("Assertion failed:", msg)
		os.Exit(1)
	}
}

func runTests() bool {
	r := NewConsistentHash(3)
	r.Add("nodeA")
	r.Add("nodeB")
	r.Add("nodeC")

	keys := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	mapping := make(map[string]string)
	for _, k := range keys {
		mapping[k] = r.Get(k)
	}
	// same key should map to same node
	for _, k := range keys {
		assert(r.Get(k) == mapping[k], "consistent mapping failed")
	}

	// remove a node and check redistribution
	r.Remove("nodeB")
	for _, k := range keys {
		newNode := r.Get(k)
		if newNode == "" {
			assert(false, "node removed but key has no mapping")
		}
		if newNode != mapping[k] {
			// acceptable redistribution
			continue
		}
		// if still same, ensure at least one key changed
	}
	// ensure at least one key changed
	changed := false
	for _, k := range keys {
		if r.Get(k) != mapping[k] {
			changed = true
			break
		}
	}
	assert(changed, "no key changed after node removal")

	// test empty ring
	empty := NewConsistentHash(3)
	assert(empty.Get("anything") == "", "empty ring should return empty string")

	return true
}

func main() {
	if !runTests() {
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}
