package payment

import (
	"context"
	"fmt"

	"github.com/omerbeden/paymentgateway/internal/adapter/provider"
	"github.com/omerbeden/paymentgateway/internal/domain/entity"
	"github.com/omerbeden/paymentgateway/internal/domain/repository"
)

type CreatePaymentUseCase struct {
	paymentRepo     repository.PaymentRepository
	providerFactory *provider.Factory
}

func NewCreatePaymentUseCase(paymentRepo repository.PaymentRepository, providerFactory *provider.Factory) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{
		paymentRepo:     paymentRepo,
		providerFactory: providerFactory,
	}
}

type CreatePaymentInput struct {
	IdempotencyKey string
	Amount         float64
	Currency       string
	ProviderID     string
	Metadata       map[string]string
}

func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (*entity.Payment, error) {
	provider, err := uc.providerFactory.GetProvider(input.ProviderID)
	if provider == nil {
		return nil, fmt.Errorf("invalid provider: %w", err)
	}

	payment := &entity.Payment{
		Amount:     input.Amount,
		Currency:   input.Currency,
		Metadata:   input.Metadata,
		Status:     entity.PaymentStatusPending,
		ProviderID: input.ProviderID,
	}
	if err := uc.paymentRepo.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	result, err := provider.CreatePayment(ctx, payment)
	if err != nil {
		payment.Status = entity.PaymentStatusFailed
		if err := uc.paymentRepo.UpdatePayment(ctx, payment); err != nil {
			return nil, fmt.Errorf("failed to update payment after provider failure: %w", err)
		}
		return nil, fmt.Errorf("provider failed to create payment: %w", err)
	}

	payment.Status = result.Status
	payment.Metadata = result.Metadata

	if err := uc.paymentRepo.UpdatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	return payment, nil
}
