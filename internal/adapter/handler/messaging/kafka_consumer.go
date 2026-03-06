package messaging

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/omerbeden/paymentgateway/internal/domain/event"

	infrakafka "github.com/omerbeden/paymentgateway/internal/infrastructure/queue/kafka"
)

type KafkaConsumer struct {
	consumer *infrakafka.Consumer
}

func NewKafkaConsumer(consumer *infrakafka.Consumer) *KafkaConsumer {
	return &KafkaConsumer{consumer: consumer}
}

func (c *KafkaConsumer) Subscribe(ctx context.Context, topics []string, handler event.ConsumerHandler) error {
	return c.consumer.Subscribe(ctx, topics, func(ctx context.Context, msg *kafka.Message) error {
		headers := make(map[string]string, len(msg.Headers))
		for _, h := range msg.Headers {
			headers[h.Key] = string(h.Value)
		}

		domainMsg := event.Message{
			Topic:   *msg.TopicPartition.Topic,
			Key:     msg.Key,
			Value:   msg.Value,
			Headers: headers,
		}
		if msg.TopicPartition.Partition >= 0 {
			domainMsg.Partition = msg.TopicPartition.Partition
		}
		if msg.TopicPartition.Offset >= 0 {
			domainMsg.Offset = int64(msg.TopicPartition.Offset)
		}
		return handler(ctx, domainMsg)
	})
}

func (c *KafkaConsumer) Close() error {
	return c.consumer.Close()
}
