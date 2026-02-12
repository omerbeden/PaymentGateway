package payment

import (
	"context"
	"fmt"

	"github.com/omerbeden/paymentgateway/internal/adapter/provider"
	"github.com/omerbeden/paymentgateway/internal/domain/entity"
	"github.com/omerbeden/paymentgateway/internal/domain/repository"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/logger"
)

type CreatePaymentUseCase struct {
	paymentRepo     repository.PaymentRepository
	providerFactory *provider.Factory
	log             logger.Logger
}

func NewCreatePaymentUseCase(
	paymentRepo repository.PaymentRepository,
	providerFactory *provider.Factory,
	log logger.Logger,
) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{
		paymentRepo:     paymentRepo,
		providerFactory: providerFactory,
		log:             log,
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

	requestID := getRequestID(ctx)
	log := uc.log.With("request_id", requestID)

	log.Info("Creating payment",
		"amount", input.Amount,
		"currency", input.Currency,
		"provider", input.ProviderID,
	)
	payment := &entity.Payment{
		Amount:     input.Amount,
		Currency:   input.Currency,
		Metadata:   input.Metadata,
		Status:     entity.PaymentStatusPending,
		ProviderID: input.ProviderID,
	}
	if err := uc.paymentRepo.CreatePayment(ctx, payment); err != nil {
		log.Error("Failed to create payment while saving to database",
			"error", err,
			"payment_id", payment.ID,
			"provider", input.ProviderID,
		)
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	result, err := provider.CreatePayment(ctx, payment)
	if err != nil {
		payment.Status = entity.PaymentStatusFailed
		if err := uc.paymentRepo.UpdatePayment(ctx, payment); err != nil {
			log.Error("Failed to create payment while updating database",
				"error", err,
				"payment_id", payment.ID,
				"provider", input.ProviderID,
			)
			return nil, fmt.Errorf("failed to update payment after provider failure: %w", err)
		}
		log.Error("Failed to create payment , provider error",
			"error", err,
			"payment_id", payment.ID,
			"provider", input.ProviderID,
		)
		return nil, fmt.Errorf("provider failed to create payment: %w", err)
	}

	payment.Status = result.Status
	payment.Metadata = result.Metadata

	if err := uc.paymentRepo.UpdatePayment(ctx, payment); err != nil {
		log.Error("Failed to create payment while updating database after provider call",
			"error", err,
			"payment_id", payment.ID,
			"provider", input.ProviderID,
		)
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	log.Info("Payment created successfully",
		"payment_id", payment.ID,
		"status", payment.Status,
		"provider", input.ProviderID,
	)

	return payment, nil
}

func getRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	return "unknown"
}
