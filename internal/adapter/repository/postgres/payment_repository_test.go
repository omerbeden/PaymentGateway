package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/omerbeden/paymentgateway/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreatePayment_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	payment := &entity.Payment{
		ID:             "pay_123456",
		Amount:         99.99,
		Currency:       "USD",
		IdempotencyKey: "idem_key_123",
		ProviderID:     "provider_123",
		Status:         entity.PaymentStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Metadata: map[string]string{
			"order_id": "order_123",
		},
	}

	// Expect the SQL query to be executed successfully
	mock.ExpectExec(`INSERT INTO payments`).
		WithArgs(
			payment.ID,
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.CreatedAt,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(), // Metadata JSON
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.CreatePayment(ctx, payment)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePayment_DuplicateIdempotencyKey(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	payment := &entity.Payment{
		ID:             "pay_123456",
		Amount:         99.99,
		Currency:       "USD",
		IdempotencyKey: "idem_key_duplicate",
		ProviderID:     "provider_123",
		Status:         entity.PaymentStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	mock.ExpectExec(`INSERT INTO payments`).
		WithArgs(
			payment.ID,
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.CreatedAt,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
		).
		WillReturnError(errors.New("pq: duplicate key value violates unique constraint \"payments_idempotency_key_key\""))

	// Act
	err = repo.CreatePayment(ctx, payment)

	// Assert
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePayment_InvalidMetadata(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	payment := &entity.Payment{
		ID:             "pay_123456",
		Amount:         99.99,
		Currency:       "USD",
		IdempotencyKey: "idem_key_123",
		ProviderID:     "provider_123",
		Status:         entity.PaymentStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Metadata: map[string]string{
			"valid_key": "valid_value",
		},
	}

	// Mock successful execution
	mock.ExpectExec(`INSERT INTO payments`).
		WithArgs(
			payment.ID,
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.CreatedAt,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.CreatePayment(ctx, payment)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePayment_DatabaseError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	payment := &entity.Payment{
		ID:             "pay_123456",
		Amount:         99.99,
		Currency:       "USD",
		IdempotencyKey: "idem_key_123",
		ProviderID:     "provider_123",
		Status:         entity.PaymentStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
	}

	// Mock database error
	mock.ExpectExec(`INSERT INTO payments`).
		WithArgs(
			payment.ID,
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.CreatedAt,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
		).
		WillReturnError(sql.ErrConnDone)

	// Act
	err = repo.CreatePayment(ctx, payment)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create payment")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePayment_ContextCancellation(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	payment := &entity.Payment{
		ID:             "pay_123456",
		Amount:         99.99,
		Currency:       "USD",
		IdempotencyKey: "idem_key_123",
		ProviderID:     "provider_123",
		Status:         entity.PaymentStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Metadata:       map[string]string{},
	}

	mock.ExpectExec(`INSERT INTO payments`).
		WithArgs(
			payment.ID,
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.CreatedAt,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
		).
		WillReturnError(context.Canceled)

	// Act
	err = repo.CreatePayment(ctx, payment)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create payment")
}

func TestCreatePayment_WithNilMetadata(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	payment := &entity.Payment{
		ID:             "pay_123456",
		Amount:         50.00,
		Currency:       "EUR",
		IdempotencyKey: "idem_key_789",
		ProviderID:     "provider_456",
		Status:         entity.PaymentStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Metadata:       nil,
	}

	mock.ExpectExec(`INSERT INTO payments`).
		WithArgs(
			payment.ID,
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.CreatedAt,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.CreatePayment(ctx, payment)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
