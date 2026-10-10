package main

import (
	"fmt"
	"time"
)

// Status represents the health state of a node.
type Status int

const (
	Alive Status = iota
	Suspect
	Dead
)

func (s Status) String() string {
	switch s {
	case Alive:
		return "Alive"
	case Suspect:
		return "Suspect"
	case Dead:
		return "Dead"
	default:
		return "Unknown"
	}
}

// Configuration constants for timeouts.
const (
	suspectTimeout = 3 * time.Second
	deadTimeout    = 5 * time.Second
)

// nodeInfo holds per‑node metadata.
type nodeInfo struct {
	id        string
	counter   uint64
	lastSeen  time.Time
	status    Status
}

// Detector implements a simple gossip‑based failure detector.
type Detector struct {
	nodes map[string]*nodeInfo
}

// NewDetector creates an empty detector.
func NewDetector() *Detector {
	return &Detector{
		nodes: make(map[string]*nodeInfo),
	}
}

// AddNode registers a new node with the detector.
func (d *Detector) AddNode(id string, now time.Time) {
	if _, exists := d.nodes[id]; exists {
		return
	}
	d.nodes[id] = &nodeInfo{
		id:       id,
		counter:  0,
		lastSeen: now,
		status:   Alive,
	}
}

// ReceiveHeartbeat records a heartbeat from the given node.
func (d *Detector) ReceiveHeartbeat(id string, now time.Time) {
	ni, ok := d.nodes[id]
	if !ok {
		// Auto‑add unknown node.
		ni = &nodeInfo{
			id: id,
		}
		d.nodes[id] = ni
	}
	ni.counter++
	ni.lastSeen = now
	ni.status = Alive
}

// Merge incorporates state from another detector (gossip).
func (d *Detector) Merge(other *Detector) {
	for id, otherInfo := range other.nodes {
		if localInfo, ok := d.nodes[id]; ok {
			// Keep the higher counter; if equal, keep the later timestamp.
			if otherInfo.counter > localInfo.counter ||
				(otherInfo.counter == localInfo.counter && otherInfo.lastSeen.After(localInfo.lastSeen)) {
				localInfo.counter = otherInfo.counter
				localInfo.lastSeen = otherInfo.lastSeen
				// Preserve status based on timestamps later in CheckFailures.
			}
		} else {
			// Clone the remote info.
			clone := *otherInfo
			d.nodes[id] = &clone
		}
	}
}

// CheckFailures updates node statuses based on elapsed time since last heartbeat.
func (d *Detector) CheckFailures(now time.Time) {
	for _, ni := range d.nodes {
		elapsed := now.Sub(ni.lastSeen)
		switch {
		case elapsed >= deadTimeout:
			ni.status = Dead
		case elapsed >= suspectTimeout:
			if ni.status != Dead {
				ni.status = Suspect
			}
		default:
			ni.status = Alive
		}
	}
}

// GetStatus returns the current status string of the node.
func (d *Detector) GetStatus(id string) string {
	if ni, ok := d.nodes[id]; ok {
		return ni.status.String()
	}
	return "Unknown"
}

// DebugString returns a formatted snapshot of the detector (for debugging).
func (d *Detector) DebugString() string {
	s := "Detector State:\n"
	for _, ni := range d.nodes {
		s += fmt.Sprintf("  ID:%s Counter:%d LastSeen:%s Status:%s\n",
			ni.id, ni.counter, ni.lastSeen.Format(time.RFC3339Nano), ni.status)
	}
	return s
}
