package main

import (
	"sync"
)

// Message represents a simple message exchanged between processes.
type Message struct {
	From int
	To   int
	Body string
}

// Process simulates a distributed process with its own inbox and state.
type Process struct {
	ID      int
	inbox   chan Message
	counter int
	mu      sync.Mutex
	handler func(p *Process, msg Message)
}

// NewProcess creates a Process with the given ID and handler.
// If handler is nil, a default handler that increments a counter on "inc" messages is used.
func NewProcess(id int, handler func(p *Process, msg Message)) *Process {
	p := &Process{
		ID:    id,
		inbox: make(chan Message, 64),
	}
	if handler != nil {
		p.handler = handler
	} else {
		p.handler = defaultHandler
	}
	return p
}

// defaultHandler increments the process's counter when the message body is "inc".
func defaultHandler(p *Process, msg Message) {
	if msg.Body == "inc" {
		p.mu.Lock()
		p.counter++
		p.mu.Unlock()
	}
}

// Start launches the process's message loop. The provided WaitGroup is signaled when the loop exits.
func (p *Process) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range p.inbox {
			p.handler(p, msg)
		}
	}()
}

// Send enqueues a message into the process's inbox.
func (p *Process) Send(msg Message) {
	p.inbox <- msg
}

// Close shuts down the process's inbox channel.
func (p *Process) Close() {
	close(p.inbox)
}

// Counter returns the current value of the process's internal counter.
func (p *Process) Counter() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.counter
}

// Network simulates a simple in-memory network connecting multiple processes.
type Network struct {
	processes map[int]*Process
	mu        sync.RWMutex
}

// NewNetwork creates an empty Network.
func NewNetwork() *Network {
	return &Network{
		processes: make(map[int]*Process),
	}
}

// AddProcess registers a process with the network.
func (n *Network) AddProcess(p *Process) {
	n.mu.Lock()
	n.processes[p.ID] = p
	n.mu.Unlock()
}

// Send delivers a message to a specific process identified by its ID.
func (n *Network) Send(to int, msg Message) {
	n.mu.RLock()
	p, ok := n.processes[to]
	n.mu.RUnlock()
	if ok {
		p.Send(msg)
	}
}

// Broadcast sends a message to all registered processes.
func (n *Network) Broadcast(msg Message) {
	n.mu.RLock()
	for _, p := range n.processes {
		p.Send(msg)
	}
	n.mu.RUnlock()
}

// Close shuts down all processes in the network.
func (n *Network) Close() {
	n.mu.RLock()
	for _, p := range n.processes {
		p.Close()
	}
	n.mu.RUnlock()
}
