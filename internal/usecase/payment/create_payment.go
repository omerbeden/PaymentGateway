package payment

import (
	"context"
	"fmt"

	"github.com/omerbeden/paymentgateway/internal/domain/entity"
	"github.com/omerbeden/paymentgateway/internal/domain/repository"
)

type CreatePaymentUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewCreatePaymentUseCase(paymentRepo repository.PaymentRepository) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{
		paymentRepo: paymentRepo,
	}
}

type CreatePaymentInput struct {
	MerchantID     string
	IdempotencyKey string
	Amount         int64
	Currency       string
	Description    string
	CustomerEmail  string
	CustomerID     string
	Metadata       map[string]string
}

func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*entity.Payment, error) {

	payment := &entity.Payment{}
	if err := uc.paymentRepo.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return payment, nil
}
