package main

import (
	"runtime"
	"sync/atomic"
	"time"
)

type Node struct {
	id int
}

var finalizerCount int64

// newEscapingNode creates a Node whose address escapes to the heap.
// A finalizer increments a counter when the object is garbage-collected.
func newEscapingNode(id int) *Node {
	n := Node{id: id}
	runtime.SetFinalizer(&n, func(p *Node) {
		atomic.AddInt64(&finalizerCount, 1)
	})
	return &n
}

// newNonEscapingNode returns a Node value that can remain on the stack.
func newNonEscapingNode(id int) Node {
	return Node{id: id}
}

// triggerGC forces a garbage collection and pauses briefly to let finalizers run.
func triggerGC() {
	runtime.GC()
	time.Sleep(10 * time.Millisecond)
}

// getFinalizerCount returns the number of times a finalizer has run.
func getFinalizerCount() int64 {
	return atomic.LoadInt64(&finalizerCount)
}
