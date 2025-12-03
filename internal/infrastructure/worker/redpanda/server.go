package redpanda

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-modular-monolith/internal/infrastructure/worker"

	"github.com/segmentio/kafka-go"
)

// RedpandaServer is a Redpanda/Kafka-based implementation of the worker.Server interface
type RedpandaServer struct {
	reader   *kafka.Reader
	handlers map[string]worker.TaskHandler
	done     chan struct{}
}

// NewRedpandaServer creates a new Redpanda server
func NewRedpandaServer(brokers []string, topic, consumerGroup string, workerCount int) *RedpandaServer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        consumerGroup,
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
		MaxBytes:       10e6, // 10MB
	})

	return &RedpandaServer{
		reader:   reader,
		handlers: make(map[string]worker.TaskHandler),
		done:     make(chan struct{}),
	}
}

// RegisterHandler registers a handler for a task type
func (s *RedpandaServer) RegisterHandler(taskName string, handler worker.TaskHandler) error {
	s.handlers[taskName] = handler
	return nil
}

// Start starts the Redpanda worker server
func (s *RedpandaServer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.done:
			return nil
		default:
		}

		msg, err := s.reader.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				return nil
			}
			return fmt.Errorf("failed to read message: %w", err)
		}

		// Get handler for this task
		taskName := string(msg.Key)
		handler, ok := s.handlers[taskName]
		if !ok {
			// No handler registered, skip this message
			continue
		}

		// Parse payload
		var payload worker.TaskPayload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			// Payload is invalid, skip
			continue
		}

		// Process the task
		if err := handler(ctx, payload); err != nil {
			// Task failed, but we still commit (fire-and-forget semantics for now)
			// In production, you might want to implement a dead-letter topic
		}
	}
}

// Stop gracefully stops the Redpanda worker server
func (s *RedpandaServer) Stop(ctx context.Context) error {
	close(s.done)
	return s.reader.Close()
}
