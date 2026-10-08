package main

import (
	"errors"
	"math"
	"sort"
	"sync"
)

// InMemoryIndex provides a thread‑safe brute‑force vector index.
type InMemoryIndex struct {
	mu      sync.RWMutex
	vectors []Vector
}

// NewInMemoryIndex creates an empty index.
func NewInMemoryIndex() *InMemoryIndex {
	return &InMemoryIndex{
		vectors: make([]Vector, 0),
	}
}

// Insert adds a new vector to the index.
// Returns an error if the vector dimensions are inconsistent with existing vectors.
func (idx *InMemoryIndex) Insert(v Vector) error {
	if len(v.Data) == 0 {
		return errors.New("vector data cannot be empty")
	}
	idx.mu.Lock()
	defer idx.mu.Unlock()
	if len(idx.vectors) > 0 && len(v.Data) != len(idx.vectors[0].Data) {
		return errors.New("vector dimension mismatch")
	}
	idx.vectors = append(idx.vectors, v)
	return nil
}

// Search returns up to k nearest vectors to the query using Euclidean distance.
func (idx *InMemoryIndex) Search(query []float64, k int) []SearchResult {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	if len(idx.vectors) == 0 || k <= 0 {
		return []SearchResult{}
	}
	if len(query) == 0 {
		return []SearchResult{}
	}
	if len(idx.vectors) > 0 && len(query) != len(idx.vectors[0].Data) {
		// dimension mismatch – return empty result set
		return []SearchResult{}
	}
	results := make([]SearchResult, 0, len(idx.vectors))
	for _, v := range idx.vectors {
		dist := euclideanDistance(query, v.Data)
		results = append(results, SearchResult{ID: v.ID, Distance: dist})
	}
	// Sort by ascending distance.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})
	if k > len(results) {
		k = len(results)
	}
	return results[:k]
}

// euclideanDistance computes the Euclidean distance between two equal‑length vectors.
func euclideanDistance(a, b []float64) float64 {
	var sum float64
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}
