package main

// EnqueueTask adds a macrotask to the end of the task queue.
func (el *EventLoop) EnqueueTask(name string, fn func(*EventLoop)) {
	el.taskQueue = append(el.taskQueue, ScheduledTask{Name: name, Fn: fn})
}

// EnqueueMicroTask adds a microtask to the end of the microtask queue.
func (el *EventLoop) EnqueueMicroTask(name string, fn func(*EventLoop)) {
	el.microTaskQueue = append(el.microTaskQueue, ScheduledTask{Name: name, Fn: fn})
}

// runScheduled executes a single ScheduledTask, managing the call stack and logging.
func (el *EventLoop) runScheduled(st ScheduledTask) {
	// Push onto call stack.
	el.callStack = append(el.callStack, st.Name)
	// Record execution.
	el.Log = append(el.Log, st.Name)
	// Execute the task's function.
	if st.Fn != nil {
		st.Fn(el)
	}
	// Pop from call stack.
	el.callStack = el.callStack[:len(el.callStack)-1]
}

// drainMicrotasks processes all pending microtasks in FIFO order.
// It continues draining until the microtask queue is empty, handling
// microtasks that schedule further microtasks.
func (el *EventLoop) drainMicrotasks() {
	for len(el.microTaskQueue) > 0 {
		// Dequeue the first microtask.
		mt := el.microTaskQueue[0]
		el.microTaskQueue = el.microTaskQueue[1:]
		el.runScheduled(mt)
	}
}

// Run executes the event loop until both the task queue and microtask queue are empty.
// The order respects JavaScript semantics: after each macrotask, all pending
// microtasks are drained before the next macrotask.
func (el *EventLoop) Run() {
	for len(el.taskQueue) > 0 {
		// Dequeue the first macrotask.
		mt := el.taskQueue[0]
		el.taskQueue = el.taskQueue[1:]
		el.runScheduled(mt)
		// After a macrotask finishes, drain all microtasks.
		el.drainMicrotasks()
	}
	// Finally, drain any stray microtasks that might have been queued
	// after the last macrotask (unlikely but safe).
	el.drainMicrotasks()
}
