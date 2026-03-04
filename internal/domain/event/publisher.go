package event

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, event DomainEvent) error
}

const (
	TopicNotificationPaymentCompleted = "notification.payment_completed"
)
