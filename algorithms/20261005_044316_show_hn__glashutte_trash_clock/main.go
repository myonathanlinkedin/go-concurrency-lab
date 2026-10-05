package main

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// ErrInvalidLength is returned when a pendulum is created with a non‑positive length.
var ErrInvalidLength = errors.New("pendulum length must be greater than zero")

// ErrInvalidGravity is returned when a pendulum is created with a non‑positive gravity constant.
var ErrInvalidGravity = errors.New("gravity must be greater than zero")

// Pendulum models a simple pendulum under the small‑angle approximation.
// Length is the distance from pivot to centre of mass in metres.
// Gravity is the local acceleration due to gravity in metres per second squared.
type Pendulum struct {
	Length  float64 // metres, must be > 0
	Gravity float64 // m/s², must be > 0
}

// NewPendulum constructs a Pendulum, validating its parameters.
// Returns an error if length or gravity are non‑positive.
func NewPendulum(length, gravity float64) (*Pendulum, error) {
	if length <= 0 {
		return nil, ErrInvalidLength
	}
	if gravity <= 0 {
		return nil, ErrInvalidGravity
	}
	return &Pendulum{
		Length:  length,
		Gravity: gravity,
	}, nil
}

// Period returns the period T of the pendulum in seconds using the formula
//   T = 2π * sqrt(L / g)
// where L is the length and g is the gravitational acceleration.
// The result is mathematically exact for the small‑angle approximation.
func (p *Pendulum) Period() float64 {
	return 2 * math.Pi * math.Sqrt(p.Length/p.Gravity)
}

// TickTimes returns a slice of time.Time values representing the moments
// when the pendulum reaches each extreme (i.e. every half period) over the
// supplied duration. The first tick is the start time (now).
// The slice always contains at least one element (the start time).
func (p *Pendulum) TickTimes(duration time.Duration) []time.Time {
	if duration < 0 {
		return []time.Time{}
	}
	halfPeriod := p.Period() / 2.0
	if halfPeriod <= 0 {
		return []time.Time{}
	}
	start := time.Now()
	// Number of intervals = floor(duration / halfPeriod)
	intervals := int(math.Floor(duration.Seconds() / halfPeriod))
	ticks := make([]time.Time, intervals+1)
	for i := 0; i <= intervals; i++ {
		ticks[i] = start.Add(time.Duration(float64(i) * halfPeriod * float64(time.Second)))
	}
	return ticks
}

// assert panics with a formatted message if condition is false.
func assert(condition bool, format string, args ...any) {
	if !condition {
		panic(fmt.Sprintf(format, args...))
	}
}

// runTests executes a suite of internal tests and reports failures via panic.
// It is deliberately lightweight to keep the module self‑contained.
func runTests() {
	// Test NewPendulum validation
	_, err := NewPendulum(0, 9.81)
	assert(err == ErrInvalidLength, "expected ErrInvalidLength, got %v", err)

	_, err = NewPendulum(1, -1)
	assert(err == ErrInvalidGravity, "expected ErrInvalidGravity, got %v", err)

	// Test Period calculation against known value.
	p, err := NewPendulum(1.0, 9.80665) // length 1 m, standard gravity
	assert(err == nil, "unexpected error: %v", err)
	expectedPeriod := 2 * math.Pi * math.Sqrt(1.0/9.80665)
	assert(math.Abs(p.Period()-expectedPeriod) < 1e-12, "period mismatch: got %v, want %v", p.Period(), expectedPeriod)

	// Test TickTimes for a 30‑minute simulation.
	simDuration := 30 * time.Minute
	ticks := p.TickTimes(simDuration)
	halfPeriod := p.Period() / 2.0
	expectedTicks := int(math.Floor(simDuration.Seconds()/halfPeriod)) + 1
	assert(len(ticks) == expectedTicks, "tick count mismatch: got %d, want %d", len(ticks), expectedTicks)

	// Ensure monotonic increase.
	for i := 1; i < len(ticks); i++ {
		assert(ticks[i].After(ticks[i-1]), "ticks not strictly increasing at index %d", i)
	}
}

// main constructs a representative pendulum (the "trash clock") and
// demonstrates its period and tick schedule for a 30‑minute run.
// It also invokes the internal test suite; any failure will cause a panic.
func main() {
	// Run internal assertions first.
	runTests()

	// Example: a 30‑minute pendulum clock built from recycled material.
	// Assume the effective length is 0.5 m.
	const length = 0.5 // metres
	const gravity = 9.80665 // m/s² (standard Earth gravity)

	clock, err := NewPendulum(length, gravity)
	if err != nil {
		panic(fmt.Sprintf("failed to create pendulum: %v", err))
	}

	period := clock.Period()
	fmt.Printf("Pendulum length: %.3f m\n", length)
	fmt.Printf("Calculated period: %.6f s (full swing)\n", period)
	fmt.Printf("Half‑period (tick interval): %.6f s\n", period/2)

	// Simulate 30 minutes of operation.
	simulation := 30 * time.Minute
	ticks := clock.TickTimes(simulation)
	fmt.Printf("Simulating %v yields %d ticks (including start).\n", simulation, len(ticks))

	// Display first few tick timestamps as a sanity check.
	fmt.Println("First five tick times:")
	for i := 0; i < 5 && i < len(ticks); i++ {
		fmt.Printf("  %d: %s\n", i, ticks[i].Format(time.RFC3339Nano))
	}
}