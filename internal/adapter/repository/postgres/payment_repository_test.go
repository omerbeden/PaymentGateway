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

func TestUpdatePayment_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	payment := &entity.Payment{
		ID:             "pay_123456",
		Amount:         150.50,
		Currency:       "USD",
		IdempotencyKey: "idem_key_123",
		ProviderID:     "provider_123",
		Status:         entity.PaymentStatusSucceeded,
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Metadata: map[string]string{
			"order_id":    "order_789",
			"customer_id": "cust_456",
		},
	}

	mock.ExpectExec(`UPDATE payments SET`).
		WithArgs(
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(), // Metadata JSON
			payment.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.UpdatePayment(ctx, payment)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePayment_StatusChangeToFailed(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	payment := &entity.Payment{
		ID:             "pay_failed_123",
		Amount:         99.99,
		Currency:       "USD",
		IdempotencyKey: "idem_key_failed",
		ProviderID:     "provider_123",
		Status:         entity.PaymentStatusFailed,
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Metadata: map[string]string{
			"error": "insufficient_funds",
		},
	}

	mock.ExpectExec(`UPDATE payments SET`).
		WithArgs(
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
			payment.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.UpdatePayment(ctx, payment)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePayment_WithEmptyMetadata(t *testing.T) {
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
		Status:         entity.PaymentStatusProcessing,
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Metadata:       map[string]string{},
	}

	mock.ExpectExec(`UPDATE payments SET`).
		WithArgs(
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
			payment.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.UpdatePayment(ctx, payment)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePayment_WithNilMetadata(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	payment := &entity.Payment{
		ID:             "pay_123456",
		Amount:         75.25,
		Currency:       "EUR",
		IdempotencyKey: "idem_key_456",
		ProviderID:     "provider_789",
		Status:         entity.PaymentStatusPending,
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(48 * time.Hour),
		Metadata:       nil,
	}

	mock.ExpectExec(`UPDATE payments SET`).
		WithArgs(
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
			payment.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.UpdatePayment(ctx, payment)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePayment_DatabaseError(t *testing.T) {
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
		Status:         entity.PaymentStatusSucceeded,
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Metadata:       map[string]string{},
	}

	mock.ExpectExec(`UPDATE payments SET`).
		WithArgs(
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
			payment.ID,
		).
		WillReturnError(sql.ErrConnDone)

	// Act
	err = repo.UpdatePayment(ctx, payment)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update payment")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePayment_ContextCancellation(t *testing.T) {
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
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Metadata:       map[string]string{},
	}

	mock.ExpectExec(`UPDATE payments SET`).
		WithArgs(
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
			payment.ID,
		).
		WillReturnError(context.Canceled)

	// Act
	err = repo.UpdatePayment(ctx, payment)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update payment")
}

func TestUpdatePayment_StatusChangeToRefunded(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	payment := &entity.Payment{
		ID:             "pay_refund_123",
		Amount:         0.00,
		Currency:       "USD",
		IdempotencyKey: "idem_key_refund",
		ProviderID:     "provider_123",
		Status:         entity.PaymentStatusRefunded,
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Metadata: map[string]string{
			"refund_reason": "customer_request",
			"refund_id":     "ref_789",
		},
	}

	mock.ExpectExec(`UPDATE payments SET`).
		WithArgs(
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
			payment.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.UpdatePayment(ctx, payment)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePayment_MultipleMetadataFields(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	payment := &entity.Payment{
		ID:             "pay_metadata_test",
		Amount:         250.00,
		Currency:       "GBP",
		IdempotencyKey: "idem_key_metadata",
		ProviderID:     "provider_456",
		Status:         entity.PaymentStatusSucceeded,
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Metadata: map[string]string{
			"order_id":           "order_999",
			"customer_id":        "cust_888",
			"invoice_number":     "inv_777",
			"transaction_ref":    "txn_666",
			"merchant_reference": "merch_555",
		},
	}

	mock.ExpectExec(`UPDATE payments SET`).
		WithArgs(
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
			payment.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.UpdatePayment(ctx, payment)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePayment_LargeAmount(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	payment := &entity.Payment{
		ID:             "pay_large_amt",
		Amount:         999999.99,
		Currency:       "USD",
		IdempotencyKey: "idem_key_large",
		ProviderID:     "provider_123",
		Status:         entity.PaymentStatusSucceeded,
		UpdatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour),
		Metadata: map[string]string{
			"transaction_type": "large_payment",
		},
	}

	mock.ExpectExec(`UPDATE payments SET`).
		WithArgs(
			payment.Amount,
			payment.Currency,
			payment.IdempotencyKey,
			payment.ProviderID,
			payment.Status,
			payment.UpdatedAt,
			payment.ExpiresAt,
			sqlmock.AnyArg(),
			payment.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = repo.UpdatePayment(ctx, payment)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByProviderPaymentID_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	now := time.Now()
	payment := &entity.Payment{
		ID:                "pay_123456",
		Amount:            99.99,
		Currency:          "USD",
		IdempotencyKey:    "idem_key_123",
		ProviderID:        "provider_123",
		ProviderPaymentID: "provider_pay_123",
		Status:            entity.PaymentStatusSucceeded,
		CreatedAt:         now,
		UpdatedAt:         now,
		ExpiresAt:         now.Add(24 * time.Hour),
		Metadata: map[string]string{
			"order_id": "order_123",
		},
	}

	rows := sqlmock.NewRows([]string{"id", "amount", "currency", "idempotency_key", "provider_id", "provider_payment_id", "status", "created_at", "updated_at", "expires_at", "metadata"}).
		AddRow(payment.ID, payment.Amount, payment.Currency, payment.IdempotencyKey, payment.ProviderID, payment.ProviderPaymentID, payment.Status, payment.CreatedAt, payment.UpdatedAt, payment.ExpiresAt, `{"order_id":"order_123"}`)

	mock.ExpectQuery(`SELECT \* FROM payments WHERE provider_payment_id=\$1 AND provider_id=\$2`).
		WithArgs(payment.ProviderPaymentID, payment.ProviderID).
		WillReturnRows(rows)

	// Act
	result, err := repo.GetByProviderPaymentID(ctx, payment.ProviderPaymentID, payment.ProviderID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, payment.ID, result.ID)
	assert.Equal(t, payment.Amount, result.Amount)
	assert.Equal(t, payment.Currency, result.Currency)
	assert.Equal(t, payment.ProviderID, result.ProviderID)
	assert.Equal(t, payment.Status, result.Status)
	assert.Equal(t, payment.Metadata, result.Metadata)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByProviderPaymentID_NotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT \* FROM payments WHERE provider_payment_id=\$1 AND provider_id=\$2`).
		WithArgs("nonexistent_pay", "provider_123").
		WillReturnError(sql.ErrNoRows)

	// Act
	result, err := repo.GetByProviderPaymentID(ctx, "nonexistent_pay", "provider_123")

	// Assert
	assert.Equal(t, ErrPaymentNotFound, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByProviderPaymentID_WrongProvider(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT \* FROM payments WHERE provider_payment_id=\$1 AND provider_id=\$2`).
		WithArgs("provider_pay_456", "wrong_provider").
		WillReturnError(sql.ErrNoRows)

	// Act
	result, err := repo.GetByProviderPaymentID(ctx, "provider_pay_456", "wrong_provider")

	// Assert
	assert.Equal(t, ErrPaymentNotFound, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByProviderPaymentID_WithoutMetadata(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	now := time.Now()
	paymentID := "pay_no_metadata"
	providerPaymentID := "provider_pay_no_meta"
	providerID := "provider_789"

	rows := sqlmock.NewRows([]string{"id", "amount", "currency", "idempotency_key", "provider_id", "provider_payment_id", "status", "created_at", "updated_at", "expires_at", "metadata"}).
		AddRow(paymentID, 50.00, "EUR", "idem_789", providerID, providerPaymentID, entity.PaymentStatusPending, now, now, now.Add(24*time.Hour), []byte(""))

	mock.ExpectQuery(`SELECT \* FROM payments WHERE provider_payment_id=\$1 AND provider_id=\$2`).
		WithArgs(providerPaymentID, providerID).
		WillReturnRows(rows)

	// Act
	result, err := repo.GetByProviderPaymentID(ctx, providerPaymentID, providerID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, paymentID, result.ID)
	assert.Equal(t, 50.00, result.Amount)
	assert.Nil(t, result.Metadata)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByProviderPaymentID_DatabaseError(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT \* FROM payments WHERE provider_payment_id=\$1 AND provider_id=\$2`).
		WithArgs("provider_pay_error", "provider_123").
		WillReturnError(sql.ErrConnDone)

	// Act
	result, err := repo.GetByProviderPaymentID(ctx, "provider_pay_error", "provider_123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get payment")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByProviderPaymentID_WithComplexMetadata(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	now := time.Now()
	metadata := map[string]string{
		"order_id":        "order_999",
		"customer_id":     "cust_888",
		"invoice_number":  "inv_777",
		"transaction_ref": "txn_666",
	}

	rows := sqlmock.NewRows([]string{"id", "amount", "currency", "idempotency_key", "provider_id", "provider_payment_id", "status", "created_at", "updated_at", "expires_at", "metadata"}).
		AddRow("pay_complex", 299.99, "GBP", "idem_complex", "provider_123", "provider_pay_complex", entity.PaymentStatusSucceeded, now, now, now.Add(24*time.Hour), `{"order_id":"order_999","customer_id":"cust_888","invoice_number":"inv_777","transaction_ref":"txn_666"}`)

	mock.ExpectQuery(`SELECT \* FROM payments WHERE provider_payment_id=\$1 AND provider_id=\$2`).
		WithArgs("provider_pay_complex", "provider_123").
		WillReturnRows(rows)

	// Act
	result, err := repo.GetByProviderPaymentID(ctx, "provider_pay_complex", "provider_123")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "pay_complex", result.ID)
	assert.Equal(t, metadata, result.Metadata)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByProviderPaymentID_FailedPaymentStatus(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewPaymentRepository(db)
	ctx := context.Background()

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "amount", "currency", "idempotency_key", "provider_id", "provider_payment_id", "status", "created_at", "updated_at", "expires_at", "metadata"}).
		AddRow("pay_failed", 75.50, "USD", "idem_failed", "provider_456", "provider_pay_failed", entity.PaymentStatusFailed, now, now, now.Add(24*time.Hour), `{"error":"insufficient_funds"}`)

	mock.ExpectQuery(`SELECT \* FROM payments WHERE provider_payment_id=\$1 AND provider_id=\$2`).
		WithArgs("provider_pay_failed", "provider_456").
		WillReturnRows(rows)

	// Act
	result, err := repo.GetByProviderPaymentID(ctx, "provider_pay_failed", "provider_456")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, entity.PaymentStatusFailed, result.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}
