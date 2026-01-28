package webhook

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/omerbeden/paymentgateway/internal/adapter/provider"
	"github.com/omerbeden/paymentgateway/internal/domain/entity"
	"github.com/omerbeden/paymentgateway/internal/domain/repository"
)

type ProcessWebHookUseCase struct {
	paymentRepo      repository.PaymentRepository
	webhookEventRepo repository.WebhookEventRepository
	providerFactory  *provider.Factory
}

func NewProcessWebHookUseCase(paymentRepo repository.PaymentRepository, providerFactory *provider.Factory) *ProcessWebHookUseCase {
	return &ProcessWebHookUseCase{
		paymentRepo:     paymentRepo,
		providerFactory: providerFactory,
	}
}

type ProcessWebHookInput struct {
	ProviderId string
	Payload    []byte
	Signature  string
}

func (uc *ProcessWebHookUseCase) Execute(ctx context.Context, input ProcessWebHookInput) error {

	providerAdapter, err := uc.providerFactory.GetProvider(input.ProviderId)
	if err != nil {
		return err
	}

	err = providerAdapter.VerifyWebhook(input.Payload, input.Signature)
	if err != nil {
		return err
	}

	webhookEvent, err := providerAdapter.ParseWebhook(input.Payload)
	if err != nil {
		return err
	}

	uc.webhookEventRepo.Save(ctx, &entity.WebhookEvent{
		ID:                uuid.New().String(),
		ProviderID:        input.ProviderId,
		ProviderPaymentID: webhookEvent.ProviderPaymentID,
		EventType:         webhookEvent.EventType,
		Signature:         input.Signature,
		Payload:           string(input.Payload),
		IsVerified:        true,
		IsProcessed:       false,
		ReceivedAt:        time.Now(),
	})

	payment, err := uc.paymentRepo.GetByProviderPaymentID(ctx, webhookEvent.ProviderPaymentID, input.ProviderId)
	if err != nil {
		return err
	}

	payment.Status = webhookEvent.Status
	payment.UpdatedAt = time.Now()

	if webhookEvent.Status == entity.PaymentStatusSucceeded {
		payment.CompletedAt = time.Now()
	}

	if err = uc.paymentRepo.UpdatePayment(ctx, payment); err != nil {
		return err
	}

	return nil
}
