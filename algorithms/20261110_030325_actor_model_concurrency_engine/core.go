package main

import (
	"sync"
)

// ActorFunc defines the function signature for handling a message.
type ActorFunc func(msg interface{})

// ActorRef represents a reference to an actor with a mailbox.
type ActorRef struct {
	inbox chan interface{}
	done  chan struct{}
	wg    sync.WaitGroup
}

// NewActor creates a new actor with the given mailbox buffer size and message handler.
// The returned ActorRef can be used to Send messages and Stop the actor.
func NewActor(bufferSize int, handler ActorFunc) *ActorRef {
	ar := &ActorRef{
		inbox: make(chan interface{}, bufferSize),
		done:  make(chan struct{}),
	}
	ar.wg.Add(1)
	go ar.loop(handler)
	return ar
}

// loop runs in a dedicated goroutine, pulling messages from the inbox and invoking the handler.
func (ar *ActorRef) loop(handler ActorFunc) {
	defer ar.wg.Done()
	for {
		select {
		case msg := <-ar.inbox:
			handler(msg)
		case <-ar.done:
			// Drain remaining messages before exiting.
			for {
				select {
				case msg := <-ar.inbox:
					handler(msg)
				default:
					return
				}
			}
		}
	}
}

// Send enqueues a message into the actor's mailbox. It blocks if the mailbox is full.
func (ar *ActorRef) Send(msg interface{}) {
	ar.inbox <- msg
}

// Stop gracefully shuts down the actor, waiting for all queued messages to be processed.
func (ar *ActorRef) Stop() {
	close(ar.done)
	ar.wg.Wait()
	close(ar.inbox)
}
