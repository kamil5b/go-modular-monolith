package worker

import (
	"fmt"

	"go-modular-monolith/internal/app/core"
	appworker "go-modular-monolith/internal/app/worker"
)

// UserModuleRegistrar implements the ModuleTaskRegistrar interface for user module tasks
type UserModuleRegistrar struct{}

// NewUserModuleRegistrar creates a new user module task registrar
func NewUserModuleRegistrar() *UserModuleRegistrar {
	return &UserModuleRegistrar{}
}

// RegisterTasks registers all user module tasks based on feature flags
func (r *UserModuleRegistrar) RegisterTasks(
	registry *appworker.TaskRegistry,
	container *core.Container,
	featureFlags *core.FeatureFlag,
) error {
	// Create user worker handler once for all email tasks
	userHandler := NewUserWorkerHandler(
		container.UserRepository,
		container.EmailClient,
	)

	// Register email notification tasks if enabled
	if featureFlags.Worker.Tasks.EmailNotifications {
		fmt.Println("[INFO] Registering user email notification tasks...")

		// Welcome email task
		registry.Register(
			TaskSendWelcomeEmail,
			userHandler.HandleSendWelcomeEmail,
		)

		// Password reset email task
		registry.Register(
			TaskSendPasswordResetEmail,
			userHandler.HandleSendPasswordResetEmail,
		)
	}

	// Register data export task if enabled
	if featureFlags.Worker.Tasks.DataExport {
		fmt.Println("[INFO] Registering user data export task...")
		registry.Register(
			TaskExportUserData,
			userHandler.HandleExportUserData,
		)
	}

	// Register report generation task if enabled
	if featureFlags.Worker.Tasks.ReportGeneration {
		fmt.Println("[INFO] Registering user report generation task...")
		registry.Register(
			TaskGenerateUserReport,
			userHandler.HandleGenerateUserReport,
		)
	}

	return nil
}
