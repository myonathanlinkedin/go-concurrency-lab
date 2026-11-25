package main

import (
	"math"
)

// Point represents a 2‑dimensional coordinate.
type Point struct {
	X float64
	Y float64
}

// dot returns the dot product of two 2‑D vectors.
func dot(a, b Point) float64 {
	return a.X*b.X + a.Y*b.Y
}

// softmax computes the softmax of a slice of scores.
// It returns a new slice with the same length.
func softmax(scores []float64) []float64 {
	max := scores[0]
	for _, v := range scores[1:] {
		if v > max {
			max = v
		}
	}
	expSum := 0.0
	expVals := make([]float64, len(scores))
	for i, v := range scores {
		// subtract max for numerical stability
		e := math.Exp(v - max)
		expVals[i] = e
		expSum += e
	}
	out := make([]float64, len(scores))
	for i, e := range expVals {
		out[i] = e / expSum
	}
	return out
}

// attention performs single‑head self‑attention on a slice of points.
// It uses identity projections for queries, keys and values.
func attention(inputs []Point) []Point {
	n := len(inputs)
	if n == 0 {
		return nil
	}
	// Compute attention output for each position.
	output := make([]Point, n)
	for i := 0; i < n; i++ {
		// Compute scores = dot(query_i, key_j) for all j.
		scores := make([]float64, n)
		for j := 0; j < n; j++ {
			scores[j] = dot(inputs[i], inputs[j])
		}
		weights := softmax(scores)
		// Weighted sum of values (which are the inputs themselves).
		var sumX, sumY float64
		for j := 0; j < n; j++ {
			sumX += weights[j] * inputs[j].X
			sumY += weights[j] * inputs[j].Y
		}
		output[i] = Point{X: sumX, Y: sumY}
	}
	return output
}

// relu applies the ReLU activation element‑wise.
func relu(p Point) Point {
	if p.X < 0 {
		p.X = 0
	}
	if p.Y < 0 {
		p.Y = 0
	}
	return p
}

// feedForward applies a linear‑ReLU‑linear block with identity weights.
func feedForward(inputs []Point) []Point {
	// With identity weights, the linear layers are no‑ops.
	out := make([]Point, len(inputs))
	for i, p := range inputs {
		out[i] = relu(p)
	}
	return out
}

// transformerBlock applies attention, adds a residual connection,
// then feed‑forward with another residual.
func transformerBlock(inputs []Point) []Point {
	att := attention(inputs)
	// Residual connection after attention.
	res1 := make([]Point, len(inputs))
	for i := range inputs {
		res1[i] = Point{
			X: inputs[i].X + att[i].X,
			Y: inputs[i].Y + att[i].Y,
		}
	}
	ff := feedForward(res1)
	// Residual connection after feed‑forward.
	out := make([]Point, len(inputs))
	for i := range res1 {
		out[i] = Point{
			X: res1[i].X + ff[i].X,
			Y: res1[i].Y + ff[i].Y,
		}
	}
	return out
}

// GenerateStrokes runs a single transformer block on the input stroke sequence.
func GenerateStrokes(seq []Point) []Point {
	return transformerBlock(seq)
}
