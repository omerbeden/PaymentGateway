package provider

import (
	"context"

	"github.com/omerbeden/paymentgateway/internal/domain/entity"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
)

type PaypalProvider struct {
}

func NewPaypalProvider(cfg config.Paypal) *PaypalProvider {
	return &PaypalProvider{}
}

func (p *PaypalProvider) CreatePayment(ctx context.Context, payment *entity.Payment) (PaymentResult, error) {
	// Implement PayPal payment creation logic here
	return PaymentResult{}, nil
}
