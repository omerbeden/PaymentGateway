package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
)

type MessageHandler func(ctx context.Context, msg *kafka.Message) error

type Consumer struct {
	c   *kafka.Consumer
	cfg config.Kafka
}

func NewConsumer(cfg config.Kafka) (*Consumer, error) {
	autoOffset := cfg.AutoOffsetReset
	if autoOffset == "" {
		autoOffset = "earliest"
	}

	cm := kafka.ConfigMap{
		"bootstrap.servers":    cfg.Brokers,
		"auto.offset.reset":    autoOffset,
		"enable.auto.commit":   false, // manual commit for at-least-once guarantee
		"session.timeout.ms":   30000,
		"max.poll.interval.ms": 300000,
		"fetch.min.bytes":      1,
		"fetch.max.bytes":      10485760,
	}

	if cfg.SASLUsername != "" {
		cm["security.protocol"] = saslProtocol(cfg.TLSEnabled)
		cm["sasl.mechanisms"] = cfg.SASLMechanism
		cm["sasl.username"] = cfg.SASLUsername
		cm["sasl.password"] = cfg.SASLPassword
	} else if cfg.TLSEnabled {
		cm["security.protocol"] = "ssl"
	}

	c, err := kafka.NewConsumer(&cm)
	if err != nil {
		return nil, fmt.Errorf("confluent consumer: create: %w", err)
	}

	return &Consumer{c: c, cfg: cfg}, nil

}

func (c *Consumer) Subscribe(ctx context.Context, topics []string, handler MessageHandler) error {
	if err := c.c.SubscribeTopics(topics, nil); err != nil {
		return fmt.Errorf("confluent consumer: subscribe %v: %w", topics, err)
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		ev := c.c.Poll(200)
		if ev == nil {
			continue
		}
		switch e := ev.(type) {
		case *kafka.Message:
			if err := handler(ctx, e); err != nil {
				fmt.Printf("message handler error: %v\n", err)
				continue
			}
			if _, err := c.c.CommitMessage(e); err != nil {
				return fmt.Errorf("confluent consumer: commit offset: %w", err)
			}
		case kafka.Error:
			if e.IsFatal() {
				return fmt.Errorf("confluent consumer: fatal error: %w", e)
			}

		default:
			fmt.Printf("ignored event: %v\n", e)
		}
	}
}

func (c *Consumer) Close() error {
	return c.c.Close()
}
