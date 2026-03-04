package event

import "time"

type EventType string

const (
	PaymentCompleted EventType = "payment.completed"
)

type DomainEvent interface {
	EventType() EventType
	AggregateID() string
	OccurredAt() time.Time
}

type BaseEvent struct {
	Type        EventType `json:"type"`
	AggregateId string    `json:"aggregate_id"`
	OccurredOn  time.Time `json:"occurred_at"`
}

func (b BaseEvent) EventType() EventType  { return b.Type }
func (b BaseEvent) AggregateID() string   { return b.AggregateId }
func (b BaseEvent) OccurredAt() time.Time { return b.OccurredOn }

type PaymentCompletedEvent struct {
	BaseEvent
	PaymentID string `json:"payment_id"`
	//CustomerID  string  `json:"customer_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Provider    string  `json:"provider"`
	Description string  `json:"description"`
}

func NewPaymentCompletedEvent(paymentID, currency, provider, description string, amount float64) PaymentCompletedEvent {
	return PaymentCompletedEvent{
		BaseEvent: BaseEvent{Type: PaymentCompleted, AggregateId: paymentID, OccurredOn: time.Now().UTC()},
		PaymentID: paymentID,
		//CustomerID:  customerID,
		Amount:      amount,
		Currency:    currency,
		Provider:    provider,
		Description: description,
	}
}
