package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-modular-monolith/internal/app/core"
	"go-modular-monolith/internal/app/worker"
	infraMongo "go-modular-monolith/internal/infrastructure/db/mongo"
	infraSQL "go-modular-monolith/internal/infrastructure/db/sql"
	userworker "go-modular-monolith/internal/modules/user/worker"
)

// RunWorker initializes and starts the worker server
func RunWorker() error {
	cfg, err := core.LoadConfig("config/config.yaml")
	if err != nil {
		return err
	}

	featureFlag, err := core.LoadFeatureFlags("config/featureflags.yaml")
	if err != nil {
		return err
	}

	// Check if workers are enabled
	if !featureFlag.Worker.Enabled || featureFlag.Worker.Backend == "disable" {
		return fmt.Errorf("workers are not enabled (feature flags: enabled=%v, backend=%s)", featureFlag.Worker.Enabled, featureFlag.Worker.Backend)
	}

	// Initialize databases
	db, err := infraSQL.Open(cfg.App.Database.SQL.DBUrl)
	if err != nil {
		if featureFlag.Repository.User == "postgres" || featureFlag.Repository.Product == "postgres" || featureFlag.Repository.Authentication == "postgres" {
			return err
		}
		fmt.Println("[WARN] PostgreSQL not loaded:", err)
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	mongo, err := infraMongo.OpenMongo(cfg.App.Database.Mongo.MongoURL)
	if err != nil {
		if featureFlag.Repository.Product == "mongo" || featureFlag.Repository.Authentication == "mongo" {
			return err
		}
		fmt.Println("[WARN] MongoDB not loaded:", err)
	}
	defer func() {
		if mongo != nil {
			infraMongo.CloseMongo(mongo)
		}
	}()

	// Create container with all dependencies
	container := core.NewContainer(*featureFlag, cfg, db, mongo)
	if container == nil {
		return errors.New("failed to create container")
	}

	// Initialize worker manager
	workerManager := worker.NewWorkerManager(container)

	// Setup module task registrations
	fmt.Println("[INFO] Setting up task registrations...")
	moduleRegistry := worker.NewModuleRegistry()

	// Register user module tasks
	moduleRegistry.Register(userworker.NewUserModuleRegistrar())

	// Register all module tasks
	if err := moduleRegistry.RegisterAll(workerManager.GetRegistry(), container, featureFlag); err != nil {
		return fmt.Errorf("failed to register module tasks: %w", err)
	}

	// Register all collected tasks with the worker server
	fmt.Println("[INFO] Registering all tasks with worker server...")
	if err := workerManager.RegisterTasks(); err != nil {
		return fmt.Errorf("failed to register tasks: %w", err)
	}

	// Create a context that can be canceled by signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start the worker server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		fmt.Printf("[INFO] Worker server running (backend: %s)\n", featureFlag.Worker.Backend)
		errChan <- workerManager.Start(ctx)
	}()

	// Wait for either a signal or an error
	select {
	case sig := <-sigChan:
		fmt.Printf("\n[INFO] Received signal %v, initiating graceful shutdown...\n", sig)
		cancel()
		// Give the server a moment to shut down cleanly
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()
		if err := workerManager.Stop(shutdownCtx); err != nil {
			fmt.Printf("[WARN] Error during worker shutdown: %v\n", err)
		}
		fmt.Println("[INFO] Worker server stopped")
		return nil
	case err := <-errChan:
		fmt.Printf("[ERROR] Worker server error: %v\n", err)
		return err
	}
}
