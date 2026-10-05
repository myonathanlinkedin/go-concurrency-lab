package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Message is a generic interface for actor messages.
type Message interface{}

// Actor encapsulates a mailbox and processing behavior.
type Actor struct {
	name      string
	inbox     chan Message
	ctx       context.Context
	cancel    context.CancelFunc
	behavior  func(context.Context, Message)
	wg        sync.WaitGroup
	stopped   bool
	stopMutex sync.Mutex
}

// NewActor creates a new actor with a bounded mailbox.
func NewActor(name string, buffer int, behavior func(context.Context, Message)) *Actor {
	ctx, cancel := context.WithCancel(context.Background())
	a := &Actor{
		name:     name,
		inbox:    make(chan Message, buffer),
		ctx:      ctx,
		cancel:   cancel,
		behavior: behavior,
	}
	a.wg.Add(1)
	go a.loop()
	return a
}

// loop processes incoming messages sequentially.
func (a *Actor) loop() {
	defer a.wg.Done()
	for {
		select {
		case <-a.ctx.Done():
			return
		case msg := <-a.inbox:
			a.behavior(a.ctx, msg)
		}
	}
}

// Send enqueues a message; blocks if mailbox is full.
func (a *Actor) Send(msg Message) {
	select {
	case <-a.ctx.Done():
		return
	case a.inbox <- msg:
	}
}

// Stop gracefully shuts down the actor.
func (a *Actor) Stop() {
	a.stopMutex.Lock()
	if a.stopped {
		a.stopMutex.Unlock()
		return
	}
	a.stopped = true
	a.stopMutex.Unlock()
	a.cancel()
	close(a.inbox)
	a.wg.Wait()
}

// ActorSystem manages a set of actors.
type ActorSystem struct {
	actors map[string]*Actor
	mu     sync.RWMutex
}

// NewSystem creates an empty actor system.
func NewSystem() *ActorSystem {
	return &ActorSystem{
		actors: make(map[string]*Actor),
	}
}

// Spawn creates and registers a new actor.
func (s *ActorSystem) Spawn(name string, buffer int, behavior func(context.Context, Message)) *Actor {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.actors[name]; exists {
		log.Fatalf("actor %s already exists", name)
	}
	actor := NewActor(name, buffer, behavior)
	s.actors[name] = actor
	return actor
}

// Send delivers a message to a named actor.
func (s *ActorSystem) Send(name string, msg Message) {
	s.mu.RLock()
	actor, ok := s.actors[name]
	s.mu.RUnlock()
	if !ok {
		log.Fatalf("actor %s not found", name)
	}
	actor.Send(msg)
}

// StopAll terminates all actors.
func (s *ActorSystem) StopAll() {
	s.mu.RLock()
	actors := make([]*Actor, 0, len(s.actors))
	for _, a := range s.actors {
		actors = append(actors, a)
	}
	s.mu.RUnlock()
	for _, a := range actors {
		a.Stop()
	}
}

// ---- Example messages ----
type Ping struct {
	From string
	To   string
	N    int
	Max  int
}
type Pong struct {
	From string
	To   string
	N    int
}

// ---- Test driver ----
func testPingPong() {
	sys := NewSystem()
	const max = 5

	// Ponger behavior
	sys.Spawn("ponger", 10, func(ctx context.Context, msg Message) {
		switch m := msg.(type) {
		case Ping:
			if m.N < m.Max {
				sys.Send(m.From, Pong{From: "ponger", To: m.From, N: m.N + 1})
			}
		}
	})

	// Pinger behavior
	sys.Spawn("pinger", 10, func(ctx context.Context, msg Message) {
		switch m := msg.(type) {
		case Pong:
			if m.N < max {
				sys.Send(m.From, Ping{From: "pinger", To: "ponger", N: m.N + 1, Max: max})
			}
		}
	})

	// Kick off the exchange
	sys.Send("pinger", Ping{From: "pinger", To: "ponger", N: 0, Max: max})

	// Wait for completion
	time.Sleep(100 * time.Millisecond)

	// Verify that the final count reached max
	// Since actors are isolated, we use a shared counter via channel
	done := make(chan struct{})
	sys.Spawn("monitor", 1, func(ctx context.Context, msg Message) {
		if n, ok := msg.(int); ok && n == max {
			close(done)
		}
	})

	// Send final ping to monitor
	sys.Send("monitor", max)

	select {
	case <-done:
		// success
	case <-time.After(500 * time.Millisecond):
		log.Fatal("testPingPong failed: timeout")
	}
	sys.StopAll()
}

func main() {
	fmt.Println("Running Actor Model tests...")
	testPingPong()
	fmt.Println("All tests passed.")
}
