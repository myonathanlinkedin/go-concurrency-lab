package main

import (
	"fmt"
	"sync"
)

func assert(condition bool, msg string) {
	if !condition {
		panic("assertion failed: " + msg)
	}
}

func testAllCommit() {
	const participantCount = 3
	voteCh := make(chan Msg, participantCount)
	coordChans := make([]chan Msg, participantCount)
	participants := make([]*Participant, participantCount)

	// Create participants that will vote to commit.
	for i := 0; i < participantCount; i++ {
		in := make(chan Msg, 2)
		coordChans[i] = in
		participants[i] = NewParticipant(i, true, in, voteCh)
	}

	coord := NewCoordinator(coordChans, voteCh)

	var wg sync.WaitGroup
	wg.Add(participantCount + 1)

	// Launch participants.
	for _, p := range participants {
		go p.Run(&wg)
	}
	// Launch coordinator.
	go coord.Run(&wg)

	// Wait for protocol to finish.
	wg.Wait()

	// Verify final states.
	assert(coord.State() == CoordCommitted, "coordinator should be committed")
	for _, p := range participants {
		assert(p.State() == PartCommitted, fmt.Sprintf("participant %d should be committed", p.id))
	}
	fmt.Println("testAllCommit passed")
}

func testAbortOnOne() {
	const participantCount = 3
	voteCh := make(chan Msg, participantCount)
	coordChans := make([]chan Msg, participantCount)
	participants := make([]*Participant, participantCount)

	// Create participants: first votes abort, others vote commit.
	for i := 0; i < participantCount; i++ {
		in := make(chan Msg, 2)
		coordChans[i] = in
		vote := true
		if i == 0 {
			vote = false
		}
		participants[i] = NewParticipant(i, vote, in, voteCh)
	}

	coord := NewCoordinator(coordChans, voteCh)

	var wg sync.WaitGroup
	wg.Add(participantCount + 1)

	// Launch participants.
	for _, p := range participants {
		go p.Run(&wg)
	}
	// Launch coordinator.
	go coord.Run(&wg)

	// Wait for protocol to finish.
	wg.Wait()

	// Verify final states.
	assert(coord.State() == CoordAborted, "coordinator should be aborted")
	for _, p := range participants {
		assert(p.State() == PartAborted, fmt.Sprintf("participant %d should be aborted", p.id))
	}
	fmt.Println("testAbortOnOne passed")
}

func main() {
	testAllCommit()
	testAbortOnOne()
	fmt.Println("All tests passed")
}
