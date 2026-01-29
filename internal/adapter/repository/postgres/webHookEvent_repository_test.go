package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/omerbeden/paymentgateway/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave_Success(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewWebHookEventRepository(db)
	ctx := context.Background()

	event := &entity.WebhookEvent{
		ProviderID:        "paypal",
		ProviderPaymentID: "payment123",
		EventType:         "PAYMENT.CAPTURE.COMPLETED",
		Signature:         "sig123",
		Payload:           `{"id":"payment123"}`,
		IsVerified:        true,
		IsProcessed:       false,
		ProcessingError:   "",
		ReceivedAt:        time.Now(),
		ProcessedAt:       time.Time{},
	}

	mock.ExpectExec(`INSERT INTO webhook_events`).
		WithArgs(event.ProviderID,
			event.ProviderPaymentID,
			event.EventType,
			event.Signature,
			event.Payload,
			event.IsVerified,
			event.IsProcessed,
			event.ProcessingError,
			event.ReceivedAt,
			event.ProcessedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	//Act
	err = repo.Save(ctx, event)

	//Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestSave_DatabaseError(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewWebHookEventRepository(db)
	ctx := context.Background()

	event := &entity.WebhookEvent{
		ProviderID:        "paypal",
		ProviderPaymentID: "payment123",
		EventType:         "PAYMENT.CAPTURE.COMPLETED",
		Signature:         "sig123",
		Payload:           `{"id":"payment123"}`,
		IsVerified:        true,
		IsProcessed:       false,
		ReceivedAt:        time.Now(),
	}

	mock.ExpectExec(`INSERT INTO webhook_events`).
		WithArgs(
			event.ProviderID,
			event.ProviderPaymentID,
			event.EventType,
			event.Signature,
			event.Payload,
			event.IsVerified,
			event.IsProcessed,
			event.ProcessingError,
			event.ReceivedAt,
			event.ProcessedAt,
		).
		WillReturnError(sql.ErrConnDone)

	//Act
	err = repo.Save(ctx, event)

	//Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save webhook event")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave_WithProcessingError(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewWebHookEventRepository(db)
	ctx := context.Background()

	now := time.Now()
	event := &entity.WebhookEvent{
		ProviderID:        "paypal",
		ProviderPaymentID: "payment456",
		EventType:         "PAYMENT.CAPTURE.DENIED",
		Signature:         "sig456",
		Payload:           `{"id":"payment456","status":"DENIED"}`,
		IsVerified:        true,
		IsProcessed:       true,
		ProcessingError:   "Payment rejected by issuer",
		ReceivedAt:        now,
		ProcessedAt:       now.Add(1 * time.Second),
	}

	mock.ExpectExec(`INSERT INTO webhook_events`).
		WithArgs(
			event.ProviderID,
			event.ProviderPaymentID,
			event.EventType,
			event.Signature,
			event.Payload,
			event.IsVerified,
			event.IsProcessed,
			event.ProcessingError,
			event.ReceivedAt,
			event.ProcessedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	//Act
	err = repo.Save(ctx, event)

	//Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave_UnverifiedEvent(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewWebHookEventRepository(db)
	ctx := context.Background()

	event := &entity.WebhookEvent{
		ProviderID:        "paypal",
		ProviderPaymentID: "payment789",
		EventType:         "PAYMENT.CAPTURE.COMPLETED",
		Signature:         "invalid_sig",
		Payload:           `{"id":"payment789"}`,
		IsVerified:        false,
		IsProcessed:       false,
		ReceivedAt:        time.Now(),
	}

	mock.ExpectExec(`INSERT INTO webhook_events`).
		WithArgs(
			event.ProviderID,
			event.ProviderPaymentID,
			event.EventType,
			event.Signature,
			event.Payload,
			event.IsVerified,
			event.IsProcessed,
			event.ProcessingError,
			event.ReceivedAt,
			event.ProcessedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	//Act
	err = repo.Save(ctx, event)

	//Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave_EmptyEvent(t *testing.T) {
	//Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewWebHookEventRepository(db)
	ctx := context.Background()

	event := &entity.WebhookEvent{}

	mock.ExpectExec(`INSERT INTO webhook_events`).
		WithArgs(
			event.ProviderID,
			event.ProviderPaymentID,
			event.EventType,
			event.Signature,
			event.Payload,
			event.IsVerified,
			event.IsProcessed,
			event.ProcessingError,
			event.ReceivedAt,
			event.ProcessedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	//Act
	err = repo.Save(ctx, event)

	//Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewWebHookEventRepository(t *testing.T) {
	//Arrange
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	//Act
	repo := NewWebHookEventRepository(db)

	//Assert
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}
