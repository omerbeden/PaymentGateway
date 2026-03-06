package consumer

import (
	"context"
	"encoding/json"

	"github.com/omerbeden/paymentgateway/internal/domain/event"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/logger"
	"github.com/omerbeden/paymentgateway/internal/usecase/notificaiton"
)

type NotificationEventConsumer struct {
	sendNotificaitonUC *notificaiton.SendPaymentNotificationUseCase
	log                logger.Logger
}

func NewNotificationEventConsumer(sendNotificaitonUC *notificaiton.SendPaymentNotificationUseCase, log logger.Logger) *NotificationEventConsumer {
	return &NotificationEventConsumer{
		sendNotificaitonUC: sendNotificaitonUC,
		log:                log,
	}
}

func (c *NotificationEventConsumer) Handle(ctx context.Context, msg event.Message) error {
	var e event.PaymentCompletedEvent
	if err := json.Unmarshal(msg.Value, &e); err != nil {
		c.log.Error("failed to unmarshal event message: %v", err)
		//route to DLQ
		return nil
	}

	c.log.Info("dispatching payment completion notification",
		"payment_id", e.PaymentID,
		"amount", e.Amount,
		"currency", e.Currency,
	)
	input := notificaiton.SendPaymentNotificationInput{
		PaymentID:     e.PaymentID,
		CustomerID:    "dummy-customer-id",
		CustomerEmail: "dummy@example.com",
		CustomerPhone: "dummy-phone",
		Amount:        e.Amount,
		Currency:      e.Currency,
		Provider:      e.Provider,
		CompletedAt:   e.OccurredAt(),
	}

	if err := c.sendNotificaitonUC.Execute(ctx, input); err != nil {
		c.log.Error("send failed, will retry",
			"payment_id", e.PaymentID,
			"err", err,
		)

		return nil
	}

	c.log.Info("notification dispatched",
		"payment_id", e.PaymentID,
	)
	return nil
}
