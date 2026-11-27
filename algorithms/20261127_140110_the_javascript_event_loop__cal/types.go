package main

// ScheduledTask represents a unit of work that can be executed by the EventLoop.
// The function receives a pointer to the EventLoop so it can schedule additional
// tasks or microtasks during its execution.
type ScheduledTask struct {
	Name string
	Fn   func(*EventLoop)
}

// EventLoop models a simplified JavaScript event loop with a call stack,
// a task queue (macrotasks) and a microtask queue.
type EventLoop struct {
	// callStack holds the names of functions currently being executed.
	callStack []string
	// taskQueue holds pending macrotasks (FIFO order).
	taskQueue []ScheduledTask
	// microTaskQueue holds pending microtasks (FIFO order, drained after each macrotask).
	microTaskQueue []ScheduledTask
	// Log records the order in which tasks and microtasks were executed.
	Log []string
}

// NewEventLoop creates an empty EventLoop ready for use.
func NewEventLoop() *EventLoop {
	return &EventLoop{
		callStack:      []string{},
		taskQueue:      []ScheduledTask{},
		microTaskQueue: []ScheduledTask{},
		Log:            []string{},
	}
}
