package worker

import (
	"context"
	"time"
)

// TaskPayload represents the data passed to a task
type TaskPayload map[string]interface{}

// TaskHandler processes a task and returns an error if processing fails
type TaskHandler func(ctx context.Context, payload TaskPayload) error

// Client is responsible for enqueueing tasks
type Client interface {
	// Enqueue adds a task to the queue
	Enqueue(ctx context.Context, taskName string, payload TaskPayload, options ...Option) error

	// EnqueueDelayed adds a task to the queue with a delay before processing
	EnqueueDelayed(ctx context.Context, taskName string, payload TaskPayload, delay time.Duration, options ...Option) error

	// Close closes the client connection
	Close() error
}

// Server is responsible for processing tasks from the queue
type Server interface {
	// RegisterHandler registers a handler for a specific task type
	RegisterHandler(taskName string, handler TaskHandler) error

	// Start starts the worker server and begins processing tasks
	Start(ctx context.Context) error

	// Stop gracefully stops the worker server
	Stop(ctx context.Context) error
}

// Option is used to specify task options like priority, retry count, timeout, etc.
type Option interface{}

// PriorityOption sets the priority of a task (higher = more important)
type PriorityOption struct {
	Priority int
}

// MaxRetriesOption sets the maximum number of retries for a task
type MaxRetriesOption struct {
	MaxRetries int
}

// TimeoutOption sets the maximum time a task can run
type TimeoutOption struct {
	Timeout time.Duration
}

// QueueOption specifies which queue to use
type QueueOption struct {
	Queue string
}

// NewPriorityOption creates a new priority option
func NewPriorityOption(priority int) *PriorityOption {
	return &PriorityOption{Priority: priority}
}

// NewMaxRetriesOption creates a new max retries option
func NewMaxRetriesOption(maxRetries int) *MaxRetriesOption {
	return &MaxRetriesOption{MaxRetries: maxRetries}
}

// NewTimeoutOption creates a new timeout option
func NewTimeoutOption(timeout time.Duration) *TimeoutOption {
	return &TimeoutOption{Timeout: timeout}
}

// NewQueueOption creates a new queue option
func NewQueueOption(queue string) *QueueOption {
	return &QueueOption{Queue: queue}
}
