package repository

import (
	"context"

	"github.com/omerbeden/paymentgateway/internal/domain/entity"
)

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	GetPaymentStatus(ctx context.Context, paymentID string) (entity.PaymentStatus, error)
}
