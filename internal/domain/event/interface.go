package event

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, event DomainEvent) error
}

type Consumer interface {
	Subscribe(ctx context.Context, topics []string, handler ConsumerHandler) error
	Close() error
}

type ConsumerHandler func(ctx context.Context, msg Message) error

type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Partition int32
	Offset    int64
	Headers   map[string]string
}

const (
	TopicNotificationPaymentCompleted = "notification.payment_completed"
)
