package repository

import (
	"context"

	"github.com/omerbeden/paymentgateway/internal/domain/entity"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	GetByProviderPaymentID(ctx context.Context, providerPaymentID, providerID string) (*entity.Payment, error)
	UpdatePayment(ctx context.Context, payment *entity.Payment) error
}
