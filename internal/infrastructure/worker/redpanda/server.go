package redpanda

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go-modular-monolith/internal/infrastructure/worker"

	"github.com/segmentio/kafka-go"
)

// TaskMetadata holds retry and tracking information for tasks
type TaskMetadata struct {
	RetryCount      int       `json:"retry_count"`
	OriginalOffset  int64     `json:"original_offset"`
	OriginalTime    time.Time `json:"original_time"`
	LastError       string    `json:"last_error"`
	CorrelationID   string    `json:"correlation_id"`
	ProcessingSteps []string  `json:"processing_steps"`
}

// RedpandaServer is a Redpanda/Kafka-based implementation of the worker.Server interface
type RedpandaServer struct {
	reader        *kafka.Reader
	handlers      map[string]worker.TaskHandler
	done          chan struct{}
	taskMetadata  map[string]*TaskMetadata // Track metadata for failed tasks
	metadataMutex sync.RWMutex
	maxRetries    int
	dlqWriter     *kafka.Writer
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

	dlqWriter := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic + "-dlq",
		Balancer: &kafka.LeastBytes{},
	}

	return &RedpandaServer{
		reader:       reader,
		handlers:     make(map[string]worker.TaskHandler),
		done:         make(chan struct{}),
		taskMetadata: make(map[string]*TaskMetadata),
		maxRetries:   3,
		dlqWriter:    dlqWriter,
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
			// Payload is invalid, send to DLQ
			s.sendToDeadLetterTopic(ctx, taskName, msg, fmt.Errorf("invalid payload: %w", err), nil)
			continue
		}

		// Get or create metadata for tracking
		taskID := s.getTaskID(msg)
		metadata := s.getTaskMetadata(taskID)
		metadata.ProcessingSteps = append(metadata.ProcessingSteps, fmt.Sprintf("attempt_%d_at_%s", metadata.RetryCount+1, time.Now().Format(time.RFC3339)))

		// Process the task
		if err := handler(ctx, payload); err != nil {
			metadata.LastError = err.Error()
			metadata.RetryCount++

			// Check if we should retry or send to DLQ
			if metadata.RetryCount >= s.maxRetries {
				s.sendToDeadLetterTopic(ctx, taskName, msg, err, metadata)
				s.removeTaskMetadata(taskID)
			}
			// Continue processing other messages
			continue
		}

		// Task succeeded, clean up metadata
		s.removeTaskMetadata(taskID)
	}
}

// Stop gracefully stops the Redpanda worker server
func (s *RedpandaServer) Stop(ctx context.Context) error {
	close(s.done)
	if s.dlqWriter != nil {
		s.dlqWriter.Close()
	}
	return s.reader.Close()
}

// sendToDeadLetterTopic sends failed tasks to a dead-letter topic with full metadata
func (s *RedpandaServer) sendToDeadLetterTopic(ctx context.Context, taskName string, msg kafka.Message, err error, metadata *TaskMetadata) error {
	// Build comprehensive metadata for DLQ
	dlqMetadata := TaskMetadata{
		RetryCount:      metadata.RetryCount,
		OriginalOffset:  msg.Offset,
		OriginalTime:    time.Now(),
		LastError:       err.Error(),
		CorrelationID:   s.getCorrelationID(msg),
		ProcessingSteps: metadata.ProcessingSteps,
	}

	// Marshal metadata as JSON
	metadataJSON, _ := json.Marshal(dlqMetadata)

	// Add comprehensive error information to headers
	errorMsg := kafka.Message{
		Key:   msg.Key,
		Value: msg.Value,
		Headers: append(msg.Headers,
			kafka.Header{Key: "error", Value: []byte(err.Error())},
			kafka.Header{Key: "retry_count", Value: []byte(fmt.Sprintf("%d", dlqMetadata.RetryCount))},
			kafka.Header{Key: "original_offset", Value: []byte(fmt.Sprintf("%d", msg.Offset))},
			kafka.Header{Key: "original_topic", Value: []byte(msg.Topic)},
			kafka.Header{Key: "correlation_id", Value: []byte(dlqMetadata.CorrelationID)},
			kafka.Header{Key: "metadata", Value: metadataJSON},
			kafka.Header{Key: "dlq_timestamp", Value: []byte(time.Now().Format(time.RFC3339))},
		),
	}

	return s.dlqWriter.WriteMessages(ctx, errorMsg)
}

// getTaskID generates a unique identifier for task tracking
func (s *RedpandaServer) getTaskID(msg kafka.Message) string {
	return fmt.Sprintf("%s-%d-%d", msg.Topic, msg.Partition, msg.Offset)
}

// getCorrelationID extracts or generates a correlation ID from message headers
func (s *RedpandaServer) getCorrelationID(msg kafka.Message) string {
	for _, h := range msg.Headers {
		if h.Key == "correlation_id" {
			return string(h.Value)
		}
	}
	// Generate new correlation ID if not present
	return fmt.Sprintf("%s-%d", time.Now().Format(time.RFC3339Nano), msg.Offset)
}

// getTaskMetadata retrieves or creates metadata for a task
func (s *RedpandaServer) getTaskMetadata(taskID string) *TaskMetadata {
	s.metadataMutex.Lock()
	defer s.metadataMutex.Unlock()

	if metadata, ok := s.taskMetadata[taskID]; ok {
		return metadata
	}

	metadata := &TaskMetadata{
		RetryCount:      0,
		OriginalTime:    time.Now(),
		ProcessingSteps: []string{},
	}
	s.taskMetadata[taskID] = metadata
	return metadata
}

// removeTaskMetadata cleans up metadata for a task
func (s *RedpandaServer) removeTaskMetadata(taskID string) {
	s.metadataMutex.Lock()
	defer s.metadataMutex.Unlock()
	delete(s.taskMetadata, taskID)
}
