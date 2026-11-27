package main

import (
	"fmt"
)

func main() {
	// Create a fresh event loop.
	loop := NewEventLoop()

	// Helper to assert conditions; panics with a clear message on failure.
	assert := func(cond bool, msg string) {
		if !cond {
			panic(fmt.Sprintf("assertion failed: %s", msg))
		}
	}

	/*
		Scenario:
		- Task A schedules microtask M1.
		- Task B schedules microtask M2 and also schedules another task C.
		- Task C schedules microtask M3.
		- Microtasks may schedule further microtasks.
		Expected execution order:
		1. A
		2. M1
		3. B
		4. M2
		5. C
		6. M3
	*/

	// Define Task A.
	loop.EnqueueTask("A", func(el *EventLoop) {
		// Schedule microtask M1.
		el.EnqueueMicroTask("M1", func(el2 *EventLoop) {
			// No further scheduling.
		})
	})

	// Define Task B.
	loop.EnqueueTask("B", func(el *EventLoop) {
		// Schedule microtask M2 which itself schedules another microtask M2a.
		el.EnqueueMicroTask("M2", func(el2 *EventLoop) {
			el2.EnqueueMicroTask("M2a", func(el3 *EventLoop) {
				// No further scheduling.
			})
		})
		// Schedule another macrotask C.
		el.EnqueueTask("C", func(el3 *EventLoop) {
			// Schedule microtask M3.
			el3.EnqueueMicroTask("M3", func(el4 *EventLoop) {
				// No further scheduling.
			})
		})
	})

	// Run the event loop.
	loop.Run()

	// Expected log.
	expected := []string{"A", "M1", "B", "M2", "M2a", "C", "M3"}

	// Verify length.
	assert(len(loop.Log) == len(expected), fmt.Sprintf("log length %d != expected %d", len(loop.Log), len(expected)))

	// Verify each entry.
	for i, exp := range expected {
		assert(loop.Log[i] == exp, fmt.Sprintf("log[%d] = %s, expected %s", i, loop.Log[i], exp))
	}

	// Edge case: Ensure that after the loop finishes, both queues are empty.
	assert(len(loop.taskQueue) == 0, "taskQueue not empty after Run")
	assert(len(loop.microTaskQueue) == 0, "microTaskQueue not empty after Run")

	// Demonstrate that scheduling a microtask after the loop has finished works
	// by creating a new loop.
	loop2 := NewEventLoop()
	loop2.EnqueueTask("X", func(el *EventLoop) {
		// No microtasks.
	})
	loop2.Run()
	// After Run, schedule a microtask; it should be executed immediately when we call drainMicrotasks.
	loop2.EnqueueMicroTask("Y", func(el *EventLoop) {})
	loop2.drainMicrotasks()
	assert(len(loop2.Log) == 2, "log length after second microtask")
	assert(loop2.Log[0] == "X", "first entry should be X")
	assert(loop2.Log[1] == "Y", "second entry should be Y")

	// If we reach this point, all assertions passed.
	fmt.Println("All assertions passed.")
}
