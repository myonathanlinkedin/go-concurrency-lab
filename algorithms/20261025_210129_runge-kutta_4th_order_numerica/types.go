package main

// ODEFunc defines a first‑order ordinary differential equation system.
// Given time t and state vector y, it returns dy/dt as a new slice.
type ODEFunc func(t float64, y []float64) []float64

// RK4 implements the classical fourth‑order Runge‑Kutta integrator.
type RK4 struct {
	f ODEFunc
}
