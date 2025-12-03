package worker

import (
	"fmt"

	"go-modular-monolith/internal/app/core"
)

// ModuleTaskRegistrar defines the interface for modules to register their tasks
type ModuleTaskRegistrar interface {
	RegisterTasks(registry *TaskRegistry, container *core.Container, featureFlags *core.FeatureFlag) error
}

// ModuleRegistry manages all module task registrations
type ModuleRegistry struct {
	modules []ModuleTaskRegistrar
}

// NewModuleRegistry creates a new module registry
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: make([]ModuleTaskRegistrar, 0),
	}
}

// Register adds a module task registrar
func (r *ModuleRegistry) Register(registrar ModuleTaskRegistrar) *ModuleRegistry {
	r.modules = append(r.modules, registrar)
	return r
}

// RegisterAll registers tasks from all modules
func (r *ModuleRegistry) RegisterAll(
	taskRegistry *TaskRegistry,
	container *core.Container,
	featureFlags *core.FeatureFlag,
) error {
	for i, module := range r.modules {
		fmt.Printf("[INFO] Registering tasks from module %d...\n", i+1)
		if err := module.RegisterTasks(taskRegistry, container, featureFlags); err != nil {
			return fmt.Errorf("failed to register module tasks: %w", err)
		}
	}
	return nil
}
