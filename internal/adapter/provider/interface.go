package provider

import (
	"context"
	"time"

	"github.com/omerbeden/paymentgateway/internal/domain/entity"
)

type PaymentProvider interface {
	CreatePayment(ctx context.Context, payment *entity.Payment) (*CreatePaymentResult, error)
	VerifyWebhook(payload []byte, signature string) error
	ParseWebhook(payload []byte) (*WebhookEvent, error)
}

type CreatePaymentResult struct {
	ProviderPaymentID string
	Status            entity.PaymentStatus
	Amount            float64
	Currency          string
	ProviderFee       int64
	PaymentURL        string
	Metadata          map[string]string
	ErrorCode         string
	ErrorMessage      string
}

type WebhookEvent struct {
	ProviderID        string
	EventType         string // "payment.succeeded", "payment.failed", "refund.succeeded"
	ProviderPaymentID string
	Status            entity.PaymentStatus
	Amount            int64
	Currency          string
	Timestamp         time.Time
	RawPayload        string
}
