package repository

import (
	"context"

	"github.com/omerbeden/paymentgateway/internal/domain/entity"
)

type WebhookEventRepository interface {
	Save(ctx context.Context, event *entity.WebhookEvent) error
}
