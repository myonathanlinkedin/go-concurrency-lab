package main

import (
	"math/rand"
)

// State represents the role of a Raft node.
type State int

const (
	Follower State = iota
	Candidate
	Leader
)

// RequestVoteArgs models the arguments of a RequestVote RPC.
type RequestVoteArgs struct {
	Term        int
	CandidateID int
}

// RequestVoteReply models the reply of a RequestVote RPC.
type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

// AppendEntriesArgs models the arguments of an AppendEntries (heartbeat) RPC.
type AppendEntriesArgs struct {
	Term     int
	LeaderID int
}

// AppendEntriesReply models the reply of an AppendEntries RPC.
type AppendEntriesReply struct {
	Term    int
	Success bool
}

// Node encapsulates the Raft election state machine.
type Node struct {
	ID               int
	Peers            []*Node // other nodes in the cluster (excluding self)
	State            State
	CurrentTerm      int
	VotedFor         *int // nil if not voted in current term
	electionTimeout  int // randomized timeout in ticks
	heartbeatTimeout int // fixed heartbeat interval in ticks
	electionTimer    int // counts up to electionTimeout
	heartbeatTimer   int // counts up to heartbeatTimeout
	randSrc          *rand.Rand
}

// NewNode creates a new Raft node with deterministic randomness.
func NewNode(id int, peers []*Node) *Node {
	src := rand.New(rand.NewSource(int64(id) + 1)) // deterministic per ID
	n := &Node{
		ID:               id,
		Peers:            peers,
		State:            Follower,
		CurrentTerm:      0,
		VotedFor:         nil,
		heartbeatTimeout: 50, // fixed heartbeat interval
		randSrc:          src,
	}
	n.resetElectionTimer()
	return n
}

// resetElectionTimer picks a new randomized election timeout.
func (n *Node) resetElectionTimer() {
	// Raft suggests 150-300ms; we use ticks: 150-300.
	n.electionTimeout = 150 + n.randSrc.Intn(151) // [150,300]
	n.electionTimer = 0
}

// Tick advances the node's internal timers by one tick.
func (n *Node) Tick() {
	switch n.State {
	case Leader:
		n.heartbeatTimer++
		if n.heartbeatTimer >= n.heartbeatTimeout {
			n.sendHeartbeats()
			n.heartbeatTimer = 0
		}
	default: // Follower or Candidate
		n.electionTimer++
		if n.electionTimer >= n.electionTimeout {
			n.startElection()
			n.electionTimer = 0
		}
	}
}

// startElection transitions the node to Candidate and initiates voting.
func (n *Node) startElection() {
	n.State = Candidate
	n.CurrentTerm++
	termAtStart := n.CurrentTerm
	n.VotedFor = &n.ID
	votes := 1 // vote for self

	// Send RequestVote RPCs to all peers.
	for _, peer := range n.Peers {
		args := RequestVoteArgs{
			Term:        termAtStart,
			CandidateID: n.ID,
		}
		reply := peer.HandleRequestVote(args)
		if reply.VoteGranted {
			votes++
		} else if reply.Term > n.CurrentTerm {
			// Found higher term; step down.
			n.becomeFollower(reply.Term)
			return
		}
	}

	majority := len(n.Peers)/2 + 1
	if votes >= majority {
		n.becomeLeader()
	} else {
		// Election failed; remain follower for next timeout.
		n.becomeFollower(n.CurrentTerm)
	}
}

// becomeFollower transitions to Follower state with given term.
func (n *Node) becomeFollower(term int) {
	if term > n.CurrentTerm {
		n.CurrentTerm = term
	}
	n.State = Follower
	n.VotedFor = nil
	n.resetElectionTimer()
}

// becomeLeader transitions to Leader state.
func (n *Node) becomeLeader() {
	n.State = Leader
	n.heartbeatTimer = 0
	n.sendHeartbeats()
}

// sendHeartbeats issues AppendEntries RPCs (heartbeat) to all peers.
func (n *Node) sendHeartbeats() {
	args := AppendEntriesArgs{
		Term:     n.CurrentTerm,
		LeaderID: n.ID,
	}
	for _, peer := range n.Peers {
		peer.HandleAppendEntries(args)
	}
}

// HandleRequestVote processes an incoming RequestVote RPC.
func (n *Node) HandleRequestVote(args RequestVoteArgs) RequestVoteReply {
	// If term is stale, reject.
	if args.Term < n.CurrentTerm {
		return RequestVoteReply{Term: n.CurrentTerm, VoteGranted: false}
	}
	// If term is newer, step down.
	if args.Term > n.CurrentTerm {
		n.becomeFollower(args.Term)
	}
	// Grant vote if not voted yet in this term.
	if n.VotedFor == nil || *n.VotedFor == args.CandidateID {
		n.VotedFor = &args.CandidateID
		n.resetElectionTimer() // grant resets election timer
		return RequestVoteReply{Term: n.CurrentTerm, VoteGranted: true}
	}
	return RequestVoteReply{Term: n.CurrentTerm, VoteGranted: false}
}

// HandleAppendEntries processes an incoming AppendEntries (heartbeat) RPC.
func (n *Node) HandleAppendEntries(args AppendEntriesArgs) AppendEntriesReply {
	if args.Term < n.CurrentTerm {
		return AppendEntriesReply{Term: n.CurrentTerm, Success: false}
	}
	// If term is newer, update term and become follower.
	if args.Term > n.CurrentTerm {
		n.becomeFollower(args.Term)
	}
	// Valid heartbeat: reset election timer and stay follower.
	n.State = Follower
	n.resetElectionTimer()
	return AppendEntriesReply{Term: n.CurrentTerm, Success: true}
}

// GetState returns a snapshot of the node's observable state.
func (n *Node) GetState() (state State, term int, leaderID int) {
	leaderID = -1
	if n.State == Leader {
		leaderID = n.ID
	}
	return n.State, n.CurrentTerm, leaderID
}
