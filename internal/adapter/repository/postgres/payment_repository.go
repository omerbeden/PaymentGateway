package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/omerbeden/paymentgateway/internal/domain/entity"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/metrics"
)

type PaymentRepository struct {
	db      *sql.DB
	metrics *metrics.Metrics
}

func NewPaymentRepository(db *sql.DB, metrics *metrics.Metrics) *PaymentRepository {
	return &PaymentRepository{db: db, metrics: metrics}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	start := time.Now()
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

	duration := time.Since(start).Seconds()
	r.metrics.DBQueryDuration.WithLabelValues("create_payment").Observe(duration)

	return nil

}

func (r *PaymentRepository) GetByProviderPaymentID(ctx context.Context, providerPaymentID, providerID string) (*entity.Payment, error) {
	start := time.Now()
	query := `SELECT * FROM payments WHERE provider_payment_id=$1 AND provider_id=$2`
	row := r.db.QueryRowContext(ctx, query, providerPaymentID, providerID)

	var p entity.Payment
	var metadataBytes []byte

	err := row.Scan(
		&p.ID,
		&p.Amount,
		&p.Currency,
		&p.IdempotencyKey,
		&p.ProviderID,
		&p.ProviderPaymentID,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.ExpiresAt,
		&metadataBytes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPaymentNotFound
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	duration := time.Since(start).Seconds()
	r.metrics.DBQueryDuration.WithLabelValues("create_payment").Observe(duration)

	if len(metadataBytes) > 0 {
		var metadata map[string]string
		if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
		p.Metadata = metadata
	}

	return &p, nil
}

func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment *entity.Payment) error {
	start := time.Now()
	query := `UPDATE payments SET 
	amount=$1,
	currency=$2, 
	idempotency_key=$3, 
	provider_id=$4, 
	status=$5, 
	updated_at=$6, 
	expires_at=$7, 
	metadata=$8 WHERE id=$9`

	jsonMetadata, err := json.Marshal(payment.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		payment.Amount,
		payment.Currency,
		payment.IdempotencyKey,
		payment.ProviderID,
		payment.Status,
		payment.UpdatedAt,
		payment.ExpiresAt,
		jsonMetadata,
		payment.ID)

	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	duration := time.Since(start).Seconds()
	r.metrics.DBQueryDuration.WithLabelValues("create_payment").Observe(duration)
	return nil

}

var (
	ErrPaymentNotFound         = errors.New("payment not found")
	ErrDuplicateIdempotencyKey = errors.New("duplicate idempotency key")
)
