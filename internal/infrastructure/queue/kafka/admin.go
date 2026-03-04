package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type TopicConfig struct {
	Name              string
	NumPartitions     int
	ReplicationFactor int
	RetentionMs       int64
}

type AdminClient struct {
	a *ckafka.AdminClient
}

func NewAdminClient(brokers string) (*AdminClient, error) {
	a, err := ckafka.NewAdminClient(&ckafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return nil, fmt.Errorf("confluent admin: create: %w", err)
	}
	return &AdminClient{a: a}, nil
}

func (a *AdminClient) EnsureTopics(ctx context.Context, topics []TopicConfig) error {
	specs := make([]ckafka.TopicSpecification, 0, len(topics))
	for _, t := range topics {
		spec := ckafka.TopicSpecification{
			Topic:             t.Name,
			NumPartitions:     t.NumPartitions,
			ReplicationFactor: t.ReplicationFactor,
		}

		if t.RetentionMs != 0 {
			spec.Config = map[string]string{
				"retention.ms": fmt.Sprint(t.RetentionMs),
			}
		}
		specs = append(specs, spec)
	}
	results, err := a.a.CreateTopics(ctx, specs, ckafka.SetAdminOperationTimeout(15*time.Second))
	if err != nil {
		return fmt.Errorf("confluent admin: CreateTopics RPC failed: %w", err)
	}

	for _, r := range results {
		if r.Error.Code() != kafka.ErrNoError && r.Error.Code() != kafka.ErrTopicAlreadyExists {
			return fmt.Errorf("confluent admin: create topic %q: %w", r.Topic, r.Error)
		}
	}
	return nil
}

func (a *AdminClient) Close() {
	a.a.Close()
}

func DefaultTopics() []TopicConfig {
	const (
		partitions  = 6
		replicas    = 1
		sevenDaysMs = int64(7 * 24 * time.Hour / time.Millisecond)
	)
	return []TopicConfig{
		{Name: "payment.created", NumPartitions: partitions, ReplicationFactor: replicas, RetentionMs: sevenDaysMs},
		{Name: "payment.processed", NumPartitions: partitions, ReplicationFactor: replicas, RetentionMs: sevenDaysMs},
		{Name: "payment.failed", NumPartitions: partitions, ReplicationFactor: replicas, RetentionMs: sevenDaysMs},
		{Name: "payment.refunded", NumPartitions: partitions, ReplicationFactor: replicas, RetentionMs: sevenDaysMs},
		{Name: "webhook.received", NumPartitions: partitions, ReplicationFactor: replicas, RetentionMs: sevenDaysMs},
		{Name: "notification.payment_completed", NumPartitions: partitions, ReplicationFactor: replicas, RetentionMs: sevenDaysMs},
		{Name: "payment.dlq", NumPartitions: 1, ReplicationFactor: replicas, RetentionMs: -1},
	}
}
