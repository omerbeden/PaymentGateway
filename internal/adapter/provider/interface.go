package provider

import (
	"context"
	"net/http"
	"time"

	"github.com/omerbeden/paymentgateway/internal/domain/entity"
)

type PaymentProvider interface {
	CreatePayment(ctx context.Context, payment *entity.Payment) (*CreatePaymentResult, error)
	VerifyWebhook(ctx context.Context, webhookCtx *WebhookContext) error
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

type WebhookContext struct {
	Payload   []byte
	Headers   http.Header
	Signature string
}

type WebhookEvent struct {
	ProviderID        string
	EventType         string
	ProviderPaymentID string
	Status            entity.PaymentStatus
	Amount            float64
	Currency          string
	CreateTime        time.Time
	RawPayload        string
}
