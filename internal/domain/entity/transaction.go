package entity

import "time"

type Transaction struct {
	ID        string            `json:"id"`
	PaymentID string            `json:"payment_id"`
	Amount    float64           `json:"amount"`
	Currency  string            `json:"currency"`
	Type      TransactionType   `json:"type"`
	Status    TransactionStatus `json:"status"`

	ProviderID      string `json:"provider_id"`
	ProviderTxnID   string `json:"provider_txn_id"`
	RequestPayload  string `json:"request_payload,omitempty"`
	ResponsePayload string `json:"response_payload,omitempty"`

	ProcessedAt time.Time `json:"processed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TransactionType string

const (
	TransactionTypeCharge  TransactionType = "charge"
	TransactionTypeRefund  TransactionType = "refund"
	TransactionTypeCapture TransactionType = "capture"
	TransactionTypeVoid    TransactionType = "void"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusSucceeded TransactionStatus = "succeeded"
	TransactionStatusFailed    TransactionStatus = "failed"
)
