package kafka

import (
	"context"
	"fmt"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
)

type Producer struct {
	p   *ckafka.Producer
	cfg config.Kafka
}

func NewKafkaProducer(cfg config.Kafka) (*Producer, error) {

	cm := ckafka.ConfigMap{
		"bootstrap.servers":  cfg.Brokers,
		"acks":               "all",
		"enable.idempotence": true,
		"retries":            5,
		"retry.backoff.ms":   200,
		"linger.ms":          5,
		"compression.type":   "snappy",
	}
	if cfg.SASLUsername != "" {
		cm["security.protocol"] = saslProtocol(cfg.TLSEnabled)
		cm["sasl.mechanisms"] = cfg.SASLMechanism
		cm["sasl.username"] = cfg.SASLUsername
		cm["sasl.password"] = cfg.SASLPassword
	} else if cfg.TLSEnabled {
		cm["security.protocol"] = "ssl"
	}

	p, err := ckafka.NewProducer(&cm)
	if err != nil {
		return nil, fmt.Errorf("confluent producer: create: %w", err)
	}
	if err != nil {
		panic(err)
	}

	prod := &Producer{p: p, cfg: cfg}

	go prod.drainEvents()

	return prod, nil
}

func (p *Producer) Produce(ctx context.Context, topic string, key, value []byte, headers map[string]string) error {
	deliveryCh := make(chan ckafka.Event, 1)

	kHeaders := make([]ckafka.Header, 0, len(headers))
	for k, v := range headers {
		kHeaders = append(kHeaders, ckafka.Header{Key: k, Value: []byte(v)})
	}

	if err := p.p.Produce(&ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: int32(ckafka.PartitionAny),
		},
		Key:     key,
		Value:   value,
		Headers: kHeaders,
	}, deliveryCh); err != nil {
		return fmt.Errorf("confluent producer: enqueue to %s: %w", topic, err)
	}

	select {
	case e := <-deliveryCh:
		m := e.(*ckafka.Message)
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("confluent producer: delivery to %s: %w", topic, m.TopicPartition.Error)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("confluent producer: context cancelled: %w", ctx.Err())
	case <-time.After(30 * time.Second):
		return fmt.Errorf("confluent producer: delivery timeout for topic %s", topic)
	}
}

func (p *Producer) Close() {
	ms := p.cfg.FlushTimeoutMs
	if ms == 0 {
		ms = 5000
	}
	p.p.Flush(ms)
	p.p.Close()
}

func (p *Producer) drainEvents() {
	for range p.p.Events() {
	}
}

func saslProtocol(tls bool) string {
	if tls {
		return "SASL_SSL"
	}
	return "SASL_PLAINTEXT"
}
