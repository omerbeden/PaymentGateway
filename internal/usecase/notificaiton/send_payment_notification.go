package notificaiton

import (
	"context"
	"fmt"
	"time"

	dnotification "github.com/omerbeden/paymentgateway/internal/domain/notification"
)

type SendPaymentNotificationInput struct {
	PaymentID     string
	CustomerID    string
	CustomerEmail string
	CustomerPhone string
	Amount        float64
	Currency      string
	Provider      string
	CompletedAt   time.Time
}

type SendPaymentNotificationUseCase struct {
	sender dnotification.Sender
}

func NewSendPaymentNotificationUseCase(sender dnotification.Sender) *SendPaymentNotificationUseCase {
	return &SendPaymentNotificationUseCase{
		sender: sender,
	}
}

func (uc *SendPaymentNotificationUseCase) Execute(ctx context.Context, input SendPaymentNotificationInput) error {
	n := dnotification.PaymentCompletedNotification{
		NotificationID: fmt.Sprintf("%s-%s", input.PaymentID, dnotification.ChannelEmail),
		CustomerID:     input.CustomerID,
		CustomerEmail:  input.CustomerEmail,
		CustomerPhone:  input.CustomerPhone,
		PaymentID:      input.PaymentID,
		Amount:         input.Amount,
		Currency:       input.Currency,
		Provider:       input.Provider,
		CompletedAt:    input.CompletedAt,
		Channel:        dnotification.ChannelEmail, // for now, we only support email notifications. In the future, we can add more channels and determine the channel based on user preferences or other factors.
	}

	//optionally, persist intent before send in case of failure

	if err := uc.sender.Send(ctx, n); err != nil {
		return fmt.Errorf("failed to send payment notification: %w", err)
	}

	return nil

}
