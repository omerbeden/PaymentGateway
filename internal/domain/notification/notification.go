package notification

import (
	"context"
	"time"
)

type Channel string

const (
	ChannelEmail Channel = "email"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusSent    Status = "sent"
	StatusFailed  Status = "failed"
)

type PaymentCompletedNotification struct {
	NotificationID string
	CustomerID     string
	CustomerEmail  string
	CustomerPhone  string
	PaymentID      string
	TransactionID  string
	Amount         float64
	Currency       string
	Provider       string
	CompletedAt    time.Time
	Channel        Channel
}

// don't need to abstract the notification struct, since we only have one type of notification for now. If we add more types in the future, we can refactor this to use an interface and multiple structs.
type Sender interface {
	Send(ctx context.Context, notification PaymentCompletedNotification) error
}
