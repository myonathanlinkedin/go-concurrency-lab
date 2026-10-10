package main

import (
	"fmt"
	"time"
)

func assert(cond bool, msg string) {
	if !cond {
		panic(fmt.Sprintf("assertion failed: %s", msg))
	}
}

// simulateTime advances the logical clock and triggers detector checks.
func simulateTime(det *Detector, start time.Time, steps int, stepDur time.Duration, heartbeatIDs []string) time.Time {
	now := start
	for i := 0; i < steps; i++ {
		// Send heartbeats for the specified IDs at this tick.
		for _, id := range heartbeatIDs {
			det.ReceiveHeartbeat(id, now)
		}
		// Perform failure detection.
		det.CheckFailures(now)
		now = now.Add(stepDur)
	}
	return now
}

func main() {
	// ---------- Test 1: All nodes stay alive ----------
	fmt.Println("Running Test 1: All nodes remain alive")
	start := time.Unix(0, 0)
	det1 := NewDetector()
	det1.AddNode("A", start)
	det1.AddNode("B", start)
	det1.AddNode("C", start)

	// Simulate 10 seconds with heartbeats every second from all nodes.
	simulateTime(det1, start, 10, 1*time.Second, []string{"A", "B", "C"})

	assert(det1.GetStatus("A") == "Alive", "Node A should be Alive")
	assert(det1.GetStatus("B") == "Alive", "Node B should be Alive")
	assert(det1.GetStatus("C") == "Alive", "Node C should be Alive")
	fmt.Println("Test 1 passed")

	// ---------- Test 2: Node B stops heartbeating ----------
	fmt.Println("Running Test 2: Node B failure detection")
	det2 := NewDetector()
	det2.AddNode("A", start)
	det2.AddNode("B", start)
	det2.AddNode("C", start)

	// First 3 seconds: all nodes heartbeat.
	simNow := simulateTime(det2, start, 3, 1*time.Second, []string{"A", "B", "C"})

	// Next 2 seconds: B stops heartbeating.
	simNow = simulateTime(det2, simNow, 2, 1*time.Second, []string{"A", "C"})

	// At this point, B should be Suspect (3s elapsed since last heartbeat).
	assert(det2.GetStatus("B") == "Suspect", "Node B should be Suspect after 5 seconds total")

	// Advance another 2 seconds without B heartbeat.
	simNow = simulateTime(det2, simNow, 2, 1*time.Second, []string{"A", "C"})

	// Now B should be Dead (5s since last heartbeat).
	assert(det2.GetStatus("B") == "Dead", "Node B should be Dead after 7 seconds total")
	assert(det2.GetStatus("A") == "Alive", "Node A should remain Alive")
	assert(det2.GetStatus("C") == "Alive", "Node C should remain Alive")
	fmt.Println("Test 2 passed")

	// ---------- Test 3: Gossip merging between two detectors ----------
	fmt.Println("Running Test 3: Gossip state merging")
	// Detector for node X
	detX := NewDetector()
	detX.AddNode("X", start)
	detX.AddNode("Y", start)

	// Detector for node Y
	detY := NewDetector()
	detY.AddNode("X", start)
	detY.AddNode("Y", start)

	// Simulate heartbeats: X sends, Y sends only for first 2 seconds.
	now := start
	for i := 0; i < 4; i++ {
		if i < 2 {
			detX.ReceiveHeartbeat("Y", now)
			detY.ReceiveHeartbeat("X", now)
		} else {
			// After i>=2, Y stops heartbeating.
			detX.ReceiveHeartbeat("X", now) // X continues.
		}
		detX.CheckFailures(now)
		detY.CheckFailures(now)

		// Perform gossip exchange.
		detX.Merge(detY)
		detY.Merge(detX)

		now = now.Add(1 * time.Second)
	}

	// After 4 seconds, Y should be Suspect in both detectors.
	assert(detX.GetStatus("Y") == "Suspect", "Node Y should be Suspect in detX")
	assert(detY.GetStatus("Y") == "Suspect", "Node Y should be Suspect in detY")

	// Advance two more seconds without Y heartbeat.
	for i := 0; i < 2; i++ {
		detX.ReceiveHeartbeat("X", now)
		detX.CheckFailures(now)
		detY.CheckFailures(now)
		detX.Merge(detY)
		detY.Merge(detX)
		now = now.Add(1 * time.Second)
	}

	// Y should now be Dead.
	assert(detX.GetStatus("Y") == "Dead", "Node Y should be Dead in detX")
	assert(detY.GetStatus("Y") == "Dead", "Node Y should be Dead in detY")
	fmt.Println("Test 3 passed")

	fmt.Println("All tests passed successfully.")
}
