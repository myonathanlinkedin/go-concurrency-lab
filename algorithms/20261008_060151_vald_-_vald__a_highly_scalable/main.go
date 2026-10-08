package main

import (
	"fmt"
)

func main() {
	idx := NewInMemoryIndex()

	// Insert a small set of vectors.
	vectors := []Vector{
		{ID: "v1", Data: []float64{1.0, 2.0, 3.0}},
		{ID: "v2", Data: []float64{2.0, 3.0, 4.0}},
		{ID: "v3", Data: []float64{-1.0, -2.0, -3.0}},
		{ID: "v4", Data: []float64{0.0, 0.0, 0.0}},
	}
	for _, v := range vectors {
		if err := idx.Insert(v); err != nil {
			panic(fmt.Sprintf("failed to insert %s: %v", v.ID, err))
		}
	}

	// Edge case: search in empty index.
	emptyIdx := NewInMemoryIndex()
	emptyRes := emptyIdx.Search([]float64{0, 0, 0}, 3)
	if len(emptyRes) != 0 {
		panic("expected empty result set for empty index")
	}

	// Edge case: k larger than number of vectors.
	largeK := idx.Search([]float64{0, 0, 0}, 10)
	if len(largeK) != len(vectors) {
		panic("expected all vectors to be returned when k exceeds count")
	}

	// Basic nearest neighbor test.
	query := []float64{1.0, 2.0, 3.0}
	results := idx.Search(query, 2)
	if len(results) != 2 {
		panic("expected exactly 2 results")
	}
	if results[0].ID != "v1" {
		panic("expected v1 to be the closest vector")
	}
	if results[0].Distance != 0.0 {
		panic("expected distance of zero for identical vector")
	}
	if results[1].ID != "v4" {
		panic("expected v4 to be the second closest vector")
	}

	// Verify distance calculations.
	expectedDistV2 := euclideanDistance(query, vectors[1].Data)
	if mathAbs(results[1].Distance-expectedDistV2) > 1e-9 {
		panic("distance mismatch for second result")
	}

	// Dimension mismatch should yield empty results.
	mismatchRes := idx.Search([]float64{1.0, 2.0}, 2)
	if len(mismatchRes) != 0 {
		panic("expected empty result set for dimension mismatch")
	}

	// Insert vector with mismatched dimension – should error.
	err := idx.Insert(Vector{ID: "bad", Data: []float64{1.0, 2.0}})
	if err == nil {
		panic("expected dimension mismatch error on insert")
	}

	fmt.Println("All assertions passed.")
}

// Helper for floating‑point comparison.
func mathAbs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
