package worker

import (
	"context"
	"fmt"

	"go-modular-monolith/internal/app/core"
	infraworker "go-modular-monolith/internal/infrastructure/worker"
)

// TaskRegistry holds a collection of task registrations
type TaskRegistry struct {
	registrations []TaskRegistration
}

// TaskRegistration defines how a task should be registered
type TaskRegistration struct {
	TaskName string
	Handler  infraworker.TaskHandler
}

// NewTaskRegistry creates a new task registry
func NewTaskRegistry() *TaskRegistry {
	return &TaskRegistry{
		registrations: make([]TaskRegistration, 0),
	}
}

// Register adds a task registration to the registry
func (r *TaskRegistry) Register(taskName string, handler infraworker.TaskHandler) *TaskRegistry {
	r.registrations = append(r.registrations, TaskRegistration{
		TaskName: taskName,
		Handler:  handler,
	})
	return r
}

// RegisterAll registers all tasks in the registry with the worker server
func (r *TaskRegistry) RegisterAll(server infraworker.Server) error {
	for _, reg := range r.registrations {
		if err := server.RegisterHandler(reg.TaskName, reg.Handler); err != nil {
			return fmt.Errorf("failed to register handler for task %s: %w", reg.TaskName, err)
		}
		fmt.Printf("[INFO] Registered handler: %s\n", reg.TaskName)
	}
	return nil
}

// WorkerManager handles worker initialization and task registration
type WorkerManager struct {
	container *core.Container
	registry  *TaskRegistry
}

// NewWorkerManager creates a new worker manager
func NewWorkerManager(container *core.Container) *WorkerManager {
	return &WorkerManager{
		container: container,
		registry:  NewTaskRegistry(),
	}
}

// GetRegistry returns the task registry for adding registrations
func (m *WorkerManager) GetRegistry() *TaskRegistry {
	return m.registry
}

// RegisterTasks registers all tasks from the registry with the worker server
func (m *WorkerManager) RegisterTasks() error {
	return m.registry.RegisterAll(m.container.WorkerServer)
}

// Start initializes and starts the worker server
func (m *WorkerManager) Start(ctx context.Context) error {
	fmt.Println("[INFO] Starting worker server...")
	fmt.Printf("[INFO] Worker server running (backend: %s)\n", "configured")
	return m.container.WorkerServer.Start(ctx)
}

// Stop gracefully stops the worker server
func (m *WorkerManager) Stop(ctx context.Context) error {
	fmt.Println("[INFO] Stopping worker server...")
	return m.container.WorkerServer.Stop(ctx)
}
