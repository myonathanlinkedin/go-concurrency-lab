package main

import "fmt"

type Vector struct {
	ID   string
	Data []float64
}

type SearchResult struct {
	ID       string
	Distance float64
}

// Index defines the operations supported by a vector search engine.
type Index interface {
	Insert(v Vector) error
	Search(query []float64, k int) []SearchResult
}

// String representations for debugging.
func (v Vector) String() string {
	return fmt.Sprintf("Vector{ID:%s, Dim:%d}", v.ID, len(v.Data))
}

func (sr SearchResult) String() string {
	return fmt.Sprintf("SearchResult{ID:%s, Distance:%.6f}", sr.ID, sr.Distance)
}
