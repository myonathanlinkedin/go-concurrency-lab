package main

import (
	"sync"
)

type PNCounter struct {
	mu sync.Mutex
	p  map[string]int
	n  map[string]int
}

func NewPNCounter() *PNCounter {
	return &PNCounter{
		p: make(map[string]int),
		n: make(map[string]int),
	}
}

func (c *PNCounter) Increment(node string, delta int) {
	if delta <= 0 {
		return
	}
	c.mu.Lock()
	c.p[node] += delta
	c.mu.Unlock()
}

func (c *PNCounter) Decrement(node string, delta int) {
	if delta <= 0 {
		return
	}
	c.mu.Lock()
	c.n[node] += delta
	c.mu.Unlock()
}

func (c *PNCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	sumP, sumN := 0, 0
	for _, v := range c.p {
		sumP += v
	}
	for _, v := range c.n {
		sumN += v
	}
	return sumP - sumN
}

func (c *PNCounter) Merge(other *PNCounter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	other.mu.Lock()
	defer other.mu.Unlock()
	for node, v := range other.p {
		if cur, ok := c.p[node]; !ok || v > cur {
			c.p[node] = v
		}
	}
	for node, v := range other.n {
		if cur, ok := c.n[node]; !ok || v > cur {
			c.n[node] = v
		}
	}
}
