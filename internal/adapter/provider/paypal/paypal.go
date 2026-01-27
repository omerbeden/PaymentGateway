package paypal

import (
	"context"

	"github.com/omerbeden/paymentgateway/internal/adapter/provider"
	"github.com/omerbeden/paymentgateway/internal/domain/entity"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
)

type Provider struct {
}

func NewProvider(cfg config.Paypal) *Provider {
	return &Provider{}
}

func (p *Provider) CreatePayment(ctx context.Context, payment *entity.Payment) (provider.PaymentResult, error) {
	// Implementation for creating a payment with PayPal
	return provider.PaymentResult{}, nil

}
