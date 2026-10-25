package main

import (
	"fmt"
	"math"
)

func main() {
	// Test 1: dy/dt = y  (exponential growth)
	fExp := func(t float64, y []float64) []float64 {
		return []float64{y[0]}
	}
	rkExp := &RK4{f: fExp}
	y0 := []float64{1.0}
	t0 := 0.0
	t1 := 1.0
	steps := 100
	y := rkExp.Integrate(y0, t0, t1, steps)
	exact := math.Exp(t1)
	if math.Abs(y[0]-exact) > 1e-5 {
		panic(fmt.Sprintf("RK4 failed for dy/dt = y: got %v, want %v", y[0], exact))
	}

	// Test 2: dy/dt = -y  (exponential decay)
	fDecay := func(t float64, y []float64) []float64 {
		return []float64{-y[0]}
	}
	rkDecay := &RK4{f: fDecay}
	y = rkDecay.Integrate(y0, t0, t1, steps)
	exact = math.Exp(-t1)
	if math.Abs(y[0]-exact) > 1e-5 {
		panic(fmt.Sprintf("RK4 failed for dy/dt = -y: got %v, want %v", y[0], exact))
	}

	// Test 3: Simple harmonic oscillator
	// dy1/dt = y2, dy2/dt = -y1  => solution y1 = cos(t), y2 = -sin(t)
	fOsc := func(t float64, y []float64) []float64 {
		return []float64{y[1], -y[0]}
	}
	rkOsc := &RK4{f: fOsc}
	y0 = []float64{1.0, 0.0}
	t1 = math.Pi / 2 // quarter period
	steps = 200
	y = rkOsc.Integrate(y0, t0, t1, steps)
	// Expected: y1 = cos(pi/2) = 0, y2 = -sin(pi/2) = -1
	if math.Abs(y[0]) > 1e-5 || math.Abs(y[1]+1) > 1e-5 {
		panic(fmt.Sprintf("RK4 failed for harmonic oscillator: got %v, want [0, -1]", y))
	}

	fmt.Println("All RK4 tests passed")
}
