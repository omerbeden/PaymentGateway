package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/omerbeden/paymentgateway/internal/domain/event"
	infrakafka "github.com/omerbeden/paymentgateway/internal/infrastructure/queue/kafka"
)

type KafkaPublisher struct {
	producer *infrakafka.Producer
}

func NewKafkaPublisher(producer *infrakafka.Producer) *KafkaPublisher {
	return &KafkaPublisher{producer: producer}
}

func (p *KafkaPublisher) Publish(ctx context.Context, topic string, event event.DomainEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("kafka publisher: marshal %T: %w", event, err)
	}
	headers := map[string]string{
		"event_type":  string(event.EventType()),
		"occurred_at": event.OccurredAt().String(),
	}

	if err := p.producer.Produce(ctx, topic, []byte(event.AggregateID()), payload, headers); err != nil {
		return fmt.Errorf("kafka publisher: produce: %w", err)
	}

	return nil
}
