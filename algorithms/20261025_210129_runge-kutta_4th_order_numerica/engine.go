package main

// Step performs a single RK4 integration step.
// It returns a newly allocated slice containing the state at t+h.
func (r *RK4) Step(y []float64, t, h float64) []float64 {
	n := len(y)

	// k1 = f(t, y)
	k1 := r.f(t, y)
	if len(k1) != n {
		panic("ODE function returned slice of incorrect length for k1")
	}

	// yTemp = y + h/2 * k1
	yTemp := make([]float64, n)
	for i := 0; i < n; i++ {
		yTemp[i] = y[i] + h*0.5*k1[i]
	}
	// k2 = f(t + h/2, yTemp)
	k2 := r.f(t+0.5*h, yTemp)
	if len(k2) != n {
		panic("ODE function returned slice of incorrect length for k2")
	}

	// yTemp = y + h/2 * k2
	for i := 0; i < n; i++ {
		yTemp[i] = y[i] + h*0.5*k2[i]
	}
	// k3 = f(t + h/2, yTemp)
	k3 := r.f(t+0.5*h, yTemp)
	if len(k3) != n {
		panic("ODE function returned slice of incorrect length for k3")
	}

	// yTemp = y + h * k3
	for i := 0; i < n; i++ {
		yTemp[i] = y[i] + h*k3[i]
	}
	// k4 = f(t + h, yTemp)
	k4 := r.f(t+h, yTemp)
	if len(k4) != n {
		panic("ODE function returned slice of incorrect length for k4")
	}

	// Combine increments
	yNext := make([]float64, n)
	for i := 0; i < n; i++ {
		yNext[i] = y[i] + (h/6.0)*(k1[i]+2.0*k2[i]+2.0*k3[i]+k4[i])
	}
	return yNext
}

// Integrate advances the system from t0 to t1 using the given number of steps.
// It returns a new slice containing the final state.
func (r *RK4) Integrate(y0 []float64, t0, t1 float64, steps int) []float64 {
	if steps <= 0 {
		panic("steps must be a positive integer")
	}
	h := (t1 - t0) / float64(steps)
	y := make([]float64, len(y0))
	copy(y, y0)
	t := t0
	for i := 0; i < steps; i++ {
		y = r.Step(y, t, h)
		t += h
	}
	return y
}
