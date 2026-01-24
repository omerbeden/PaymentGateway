package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/omerbeden/paymentgateway/internal/domain/entity"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	query := `INSERT INTO payments (id, amount, currency, idempotency_key, provider_id, status, created_at, updated_at, expires_at,metadata)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	jsonMetadata, err := json.Marshal(payment.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		payment.ID,
		payment.Amount,
		payment.Currency,
		payment.IdempotencyKey,
		payment.ProviderID,
		payment.Status,
		payment.CreatedAt,
		payment.UpdatedAt,
		payment.ExpiresAt,
		jsonMetadata)

	if err != nil {
		// Check for unique constraint violation (idempotency key)
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return ErrDuplicateIdempotencyKey
			}
		}
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil

}

func (r *PaymentRepository) GetPaymentStatus(ctx context.Context, paymentID string) (entity.PaymentStatus, error) {
	return entity.PaymentStatusPending, nil
}

var (
	ErrPaymentNotFound         = errors.New("payment not found")
	ErrDuplicateIdempotencyKey = errors.New("duplicate idempotency key")
)
