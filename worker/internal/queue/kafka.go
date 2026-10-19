package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Job represents an oracle request flowing through Kafka.
type Job struct {
	ID   string `json:"id"`
	Pair string `json:"pair"`
}

// Producer publishes jobs to Kafka.
type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{writer: w}
}

func (p *Producer) Publish(ctx context.Context, job Job) error {
	body, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal job: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(job.Pair),
		Value: body,
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error { return p.writer.Close() }

// Consumer reads jobs from Kafka.
type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:               brokers,
		Topic:                 topic,
		GroupID:               groupID,
		WatchPartitionChanges: true,
	})
	return &Consumer{reader: r}
}

// Consume blocks, calling handler for each job. Returns when the context is cancelled.
func (c *Consumer) Consume(ctx context.Context, handler func(context.Context, Job) error) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return err
		}

		var job Job
		if err := json.Unmarshal(msg.Value, &job); err != nil {
			// Invalid message — log and commit to skip
			if cerr := c.reader.CommitMessages(ctx, msg); cerr != nil {
				return cerr
			}
			continue
		}

		if err := handler(ctx, job); err != nil {
			// Handler failed — don't commit, message will be retried
			// In production, add a dead letter queue for repeated failures
			return err
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			return err
		}
	}
}

func (c *Consumer) Close() error { return c.reader.Close() }
