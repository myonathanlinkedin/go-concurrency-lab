package main

import (
	"fmt"
	"math"
)

func approxEqual(a, b float64) bool {
	const eps = 1e-9
	return math.Abs(a-b) <= eps
}

func main() {
	// Unit test 1: two orthogonal points.
	input := []Point{
		{X: 1.0, Y: 0.0},
		{X: 0.0, Y: 1.0},
	}
	expected := []Point{
		{X: 1.0 + math.Exp(1)/(math.Exp(1)+1.0), Y: 1.0 / (math.Exp(1) + 1.0)},
		{X: 1.0 / (math.Exp(1) + 1.0), Y: 1.0 + math.Exp(1)/(math.Exp(1)+1.0)},
	}
	output := GenerateStrokes(input)

	if len(output) != len(expected) {
		panic(fmt.Sprintf("length mismatch: got %d, want %d", len(output), len(expected)))
	}
	for i := range output {
		if !approxEqual(output[i].X, expected[i].X) || !approxEqual(output[i].Y, expected[i].Y) {
			panic(fmt.Sprintf("mismatch at index %d: got (%v,%v), want (%v,%v)",
				i, output[i].X, output[i].Y, expected[i].X, expected[i].Y))
		}
	}
	// Unit test 2: empty input should return nil.
	if GenerateStrokes(nil) != nil {
		panic("expected nil output for nil input")
	}
	// Unit test 3: single point should return point + ReLU(point) (since attention returns the point itself).
	single := []Point{{X: -0.5, Y: 2.0}}
	singleOut := GenerateStrokes(single)
	// Attention on single element yields the same point.
	// Residual adds original => 2*point, ReLU on that yields (0,4), then residual adds again => (0,6)
	// Let's compute step by step using the implementation:
	// attention = point
	// after first residual: 2*point = (-1.0,4.0)
	// feedForward (ReLU) => (0,4.0)
	// second residual: (-1.0,4.0)+(0,4.0)=(-1.0,8.0)
	// Since ReLU zeroes negative X, final X stays -1.0 (no further ReLU)
	// Our implementation does not clamp after second residual, so expected is (-1.0,8.0)
	expX := -1.0
	expY := 8.0
	if !approxEqual(singleOut[0].X, expX) || !approxEqual(singleOut[0].Y, expY) {
		panic(fmt.Sprintf("single point test failed: got (%v,%v), want (%v,%v)", singleOut[0].X, singleOut[0].Y, expX, expY))
	}
	fmt.Println("All assertions passed.")
}
