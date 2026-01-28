package entity

import "time"

type WebhookEvent struct {
	ID                string    `json:"id"`
	ProviderID        string    `json:"provider_id,omitempty"`
	ProviderPaymentID string    `json:"provider_payment_id,omitempty"`
	EventType         string    `json:"event_type,omitempty"`
	Signature         string    `json:"signature,omitempty"`
	Payload           string    `json:"payload,omitempty"`
	IsVerified        bool      `json:"is_verified,omitempty"`
	IsProcessed       bool      `json:"is_processed,omitempty"`
	ProcessingError   string    `json:"processing_error,omitempty"`
	ReceivedAt        time.Time `json:"received_at,omitempty"`
	ProcessedAt       time.Time `json:"processed_at,omitempty"`
}
