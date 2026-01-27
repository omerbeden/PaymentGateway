package provider

import (
	"context"

	"github.com/omerbeden/paymentgateway/internal/domain/entity"
)

type PaymentProvider interface {
	CreatePayment(ctx context.Context, payment *entity.Payment) (*CreatePaymentResult, error)
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
