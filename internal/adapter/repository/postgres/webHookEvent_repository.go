package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/omerbeden/paymentgateway/internal/domain/entity"
)

type WebHookEventRepository struct {
	db *sql.DB
}

func NewWebHookEventRepository(db *sql.DB) *WebHookEventRepository {
	return &WebHookEventRepository{db: db}
}
func (r *WebHookEventRepository) Save(ctx context.Context, event *entity.WebhookEvent) error {

	query := `INSERT INTO webhook_events (name, provider_id,
	 	provider_payment_id, event_type, signature, payload,
		is_verified, is_processed, processing_error, received_at, processed_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query, event.ProviderID,
		event.ProviderPaymentID, event.EventType, event.Signature, event.Payload,
		event.IsVerified, event.IsProcessed, event.ProcessingError, event.ReceivedAt, event.ProcessedAt)
	if err != nil {
		return fmt.Errorf("failed to save webhook event: %w", err)
	}
	return nil
}
