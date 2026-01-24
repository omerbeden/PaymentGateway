package provider

import (
	"context"

	"github.com/omerbeden/paymentgateway/internal/domain/entity"
)

type PaymentProvider interface {
	CreatePayment(ctx context.Context, payment *entity.Payment) (PaymentResult, error)
}

type PaymentResult struct {
	ProviderPaymentID string
	Status            entity.PaymentStatus
	Amount            int64
	Currency          string
	ProviderFee       int64
	PaymentURL        string
	Metadata          map[string]string
	ErrorCode         string
	ErrorMessage      string
}
