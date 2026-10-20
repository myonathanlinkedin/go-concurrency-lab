package main

import (
	"fmt"
	"sync"
)

// Message types used in the protocol.
type MsgType int

const (
	MsgPrepare MsgType = iota
	MsgVoteCommit
	MsgVoteAbort
	MsgCommit
	MsgAbort
)

// Msg represents a protocol message.
type Msg struct {
	Type MsgType
	From int // participant id, -1 for coordinator
}

// CoordinatorState enumerates possible coordinator states.
type CoordinatorState int

const (
	CoordInit CoordinatorState = iota
	CoordWaiting
	CoordCommitted
	CoordAborted
)

// ParticipantState enumerates possible participant states.
type ParticipantState int

const (
	PartInit ParticipantState = iota
	PartPrepared
	PartCommitted
	PartAborted
)

// Coordinator orchestrates the two‑phase commit.
type Coordinator struct {
	toParts []chan Msg   // outbound channels to participants
	voteCh  chan Msg     // inbound channel for votes
	state   CoordinatorState
}

// NewCoordinator creates a coordinator given participant outbound channels.
func NewCoordinator(toParts []chan Msg, voteCh chan Msg) *Coordinator {
	return &Coordinator{
		toParts: toParts,
		voteCh:  voteCh,
		state:   CoordInit,
	}
}

// Run executes the two‑phase commit protocol.
// It must be launched as a goroutine; wg signals completion.
func (c *Coordinator) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	// Phase 1: send Prepare to all participants.
	for _, ch := range c.toParts {
		ch <- Msg{Type: MsgPrepare, From: -1}
	}
	c.state = CoordWaiting

	// Collect votes.
	commit := true
	for i := 0; i < len(c.toParts); i++ {
		vote := <-c.voteCh
		if vote.Type != MsgVoteCommit {
			commit = false
		}
	}

	// Phase 2: broadcast decision.
	decision := MsgCommit
	if !commit {
		decision = MsgAbort
	}
	for _, ch := range c.toParts {
		ch <- Msg{Type: decision, From: -1}
	}
	if commit {
		c.state = CoordCommitted
	} else {
		c.state = CoordAborted
	}
}

// Participant implements the participant side of 2PC.
type Participant struct {
	id          int
	inbound     chan Msg // messages from coordinator
	outbound    chan Msg // votes back to coordinator
	voteCommit  bool
	state       ParticipantState
	decisionLog []MsgType // optional log for debugging
}

// NewParticipant creates a participant with deterministic vote behavior.
func NewParticipant(id int, voteCommit bool, inbound chan Msg, outbound chan Msg) *Participant {
	return &Participant{
		id:         id,
		inbound:    inbound,
		outbound:   outbound,
		voteCommit: voteCommit,
		state:      PartInit,
	}
}

// Run processes incoming messages until a final decision is received.
// It must be launched as a goroutine; wg signals completion.
func (p *Participant) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		msg := <-p.inbound
		switch msg.Type {
		case MsgPrepare:
			// Vote according to predetermined behavior.
			if p.voteCommit {
				p.outbound <- Msg{Type: MsgVoteCommit, From: p.id}
			} else {
				p.outbound <- Msg{Type: MsgVoteAbort, From: p.id}
			}
			p.state = PartPrepared
		case MsgCommit:
			p.state = PartCommitted
			return
		case MsgAbort:
			p.state = PartAborted
			return
		default:
			// Unexpected message; panic for visibility.
			panic(fmt.Sprintf("participant %d received unknown message type %d", p.id, msg.Type))
		}
	}
}

// State returns the current participant state (read‑only).
func (p *Participant) State() ParticipantState {
	return p.state
}

// State returns the current coordinator state (read‑only).
func (c *Coordinator) State() CoordinatorState {
	return c.state
}
