package entity

import (
	"time"
)

type Payment struct {
	ID             string            `json:"id"`
	Amount         float64           `json:"amount"`
	Currency       string            `json:"currency"`
	IdempotencyKey string            `json:"idempotency_key"`
	ProviderID     string            `json:"provider_id"`
	Status         PaymentStatus     `json:"status"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	CompletedAt    time.Time         `json:"completed_at,omitempty"`
	ExpiresAt      time.Time         `json:"expires_at,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

type PaymentStatus string

const (
	PaymentStatusPending       PaymentStatus = "pending"
	PaymentStatusProcessing    PaymentStatus = "processing"
	PaymentStatusSucceeded     PaymentStatus = "succeeded"
	PaymentStatusFailed        PaymentStatus = "failed"
	PaymentStatusCancelled     PaymentStatus = "cancelled"
	PaymentStatusRefunded      PaymentStatus = "refunded"
	PaymentStatusPartialRefund PaymentStatus = "partial_refund"
)
