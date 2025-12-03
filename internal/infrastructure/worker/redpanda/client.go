package redpanda

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-modular-monolith/internal/infrastructure/worker"

	"github.com/segmentio/kafka-go"
)

// RedpandaClient is a Redpanda/Kafka-based implementation of the worker.Client interface
type RedpandaClient struct {
	writer *kafka.Writer
	topic  string
}

// NewRedpandaClient creates a new Redpanda client
func NewRedpandaClient(brokers []string, topic string) *RedpandaClient {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &RedpandaClient{
		writer: writer,
		topic:  topic,
	}
}

// Enqueue enqueues a task immediately
func (c *RedpandaClient) Enqueue(
	ctx context.Context,
	taskName string,
	payload worker.TaskPayload,
	options ...worker.Option,
) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	return c.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(taskName),
		Value: data,
	})
}

// EnqueueDelayed enqueues a task with a delay
// Note: Kafka/Redpanda doesn't natively support delayed delivery
// For delayed tasks, consider using a separate delayed topic + scheduler
func (c *RedpandaClient) EnqueueDelayed(
	ctx context.Context,
	taskName string,
	payload worker.TaskPayload,
	delay time.Duration,
	options ...worker.Option,
) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// For now, just enqueue immediately
	// In production, you'd implement a separate delayed topic scheduler
	return c.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(taskName),
		Value: data,
	})
}

// Close closes the Redpanda client
func (c *RedpandaClient) Close() error {
	return c.writer.Close()
}
