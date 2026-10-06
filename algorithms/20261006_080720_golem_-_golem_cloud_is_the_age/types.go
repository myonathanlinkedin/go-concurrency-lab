package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// AgentStatus represents the lifecycle state of an AI agent.
type AgentStatus int

const (
	AgentStatusIdle AgentStatus = iota
	AgentStatusRunning
	AgentStatusPaused
	AgentStatusTerminated
)

func (s AgentStatus) String() string {
	switch s {
	case AgentStatusIdle:
		return "Idle"
	case AgentStatusRunning:
		return "Running"
	case AgentStatusPaused:
		return "Paused"
	case AgentStatusTerminated:
		return "Terminated"
	default:
		return "Unknown"
	}
}

// TaskPriority determines the scheduling weight of a task.
type TaskPriority int

const (
	PriorityLow TaskPriority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// Task represents a unit of work to be executed by an agent.
type Task struct {
	ID        string
	Priority  TaskPriority
	Payload   []byte
	CreatedAt time.Time
	Deadline  time.Time
}

// Agent represents an autonomous AI agent entity.
type Agent struct {
	ID        string
	Name      string
	Status    AgentStatus
	CreatedAt time.Time
	LastSeen  time.Time
	// Mutex to protect status changes
	mu sync.RWMutex
}

// NewAgent creates a new Agent instance.
func NewAgent(id, name string) *Agent {
	now := time.Now()
	return &Agent{
		ID:        id,
		Name:      name,
		Status:    AgentStatusIdle,
		CreatedAt: now,
		LastSeen:  now,
	}
}

// SetStatus atomically updates the agent's status.
func (a *Agent) SetStatus(status AgentStatus) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Status = status
	a.LastSeen = time.Now()
}

// GetStatus atomically retrieves the agent's status.
func (a *Agent) GetStatus() AgentStatus {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Status
}

// TaskResult encapsulates the outcome of a task execution.
type TaskResult struct {
	TaskID    string
	AgentID   string
	Success   bool
	Error     error
	Duration  time.Duration
	Output    []byte
	Timestamp time.Time
}

// GolemError represents a specific error type for the Golem platform.
type GolemError struct {
	Code    int
	Message string
}

func (e *GolemError) Error() string {
	return fmt.Sprintf("Golem Error [%d]: %s", e.Code, e.Message)
}

// Common error codes
const (
	ErrCodeAgentNotFound = 1001
	ErrCodeTaskNotFound  = 1002
	ErrCodeInvalidState  = 1003
	ErrCodeTimeout       = 1004
)

// ErrAgentNotFound is returned when an agent ID is not found.
var ErrAgentNotFound = &GolemError{Code: ErrCodeAgentNotFound, Message: "agent not found"}

// ErrTaskNotFound is returned when a task ID is not found.
var ErrTaskNotFound = &GolemError{Code: ErrCodeTaskNotFound, Message: "task not found"}

// ErrInvalidState is returned when an operation is invalid for the current state.
var ErrInvalidState = &GolemError{Code: ErrCodeInvalidState, Message: "invalid state transition"}

// ErrTimeout is returned when an operation exceeds its deadline.
var ErrTimeout = &GolemError{Code: ErrCodeTimeout, Message: "operation timed out"}

// AgentHandler is the interface for processing tasks.
type AgentHandler interface {
	// HandleTask processes a task and returns a result.
	HandleTask(task *Task) (*TaskResult, error)
	// Shutdown gracefully stops the agent.
	Shutdown() error
}

// TaskQueue is an interface for task management.
type TaskQueue interface {
	// Enqueue adds a task to the queue.
	Enqueue(task *Task) error
	// Dequeue removes and returns the next task.
	Dequeue() (*Task, error)
	// Peek returns the next task without removing it.
	Peek() (*Task, error)
	// Len returns the number of tasks in the queue.
	Len() int
	// Clear removes all tasks from the queue.
	Clear()
}

// GolemEngine is the core orchestration engine.
type GolemEngine struct {
	agents    map[string]*Agent
	tasks     map[string]*Task
	queue     TaskQueue
	handlers  map[string]AgentHandler
	mu        sync.RWMutex
	running   bool
	stopCh    chan struct{}
}

// NewGolemEngine creates a new GolemEngine instance.
func NewGolemEngine() *GolemEngine {
	return &GolemEngine{
		agents:   make(map[string]*Agent),
		tasks:    make(map[string]*Task),
		queue:    NewPriorityQueue(),
		handlers: make(map[string]AgentHandler),
		stopCh:   make(chan struct{}),
	}
}

// PriorityQueue is a thread-safe priority queue implementation.
type PriorityQueue struct {
	mu     sync.Mutex
	items  []*Task
	length int
}

// NewPriorityQueue creates a new PriorityQueue.
func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		items: make([]*Task, 0),
	}
}

// Enqueue adds a task to the priority queue.
func (pq *PriorityQueue) Enqueue(task *Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}
	pq.mu.Lock()
	defer pq.mu.Unlock()
	pq.items = append(pq.items, task)
	pq.length++
	return nil
}

// Dequeue removes and returns the highest priority task.
func (pq *PriorityQueue) Dequeue() (*Task, error) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.length == 0 {
		return nil, ErrTaskNotFound
	}
	// Find the highest priority task
	maxIdx := 0
	for i := 1; i < pq.length; i++ {
		if pq.items[i].Priority > pq.items[maxIdx].Priority {
			maxIdx = i
		}
	}
	task := pq.items[maxIdx]
	// Remove the task
	pq.items[maxIdx] = pq.items[pq.length-1]
	pq.items = pq.items[:pq.length-1]
	pq.length--
	return task, nil
}

// Peek returns the highest priority task without removing it.
func (pq *PriorityQueue) Peek() (*Task, error) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.length == 0 {
		return nil, ErrTaskNotFound
	}
	maxIdx := 0
	for i := 1; i < pq.length; i++ {
		if pq.items[i].Priority > pq.items[maxIdx].Priority {
			maxIdx = i
		}
	}
	return pq.items[maxIdx], nil
}

// Len returns the number of tasks in the queue.
func (pq *PriorityQueue) Len() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return pq.length
}

// Clear removes all tasks from the queue.
func (pq *PriorityQueue) Clear() {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	pq.items = make([]*Task, 0)
	pq.length = 0
}
