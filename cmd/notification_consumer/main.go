package consumer

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/omerbeden/paymentgateway/internal/adapter/handler/messaging"
	"github.com/omerbeden/paymentgateway/internal/adapter/handler/messaging/consumer"
	"github.com/omerbeden/paymentgateway/internal/domain/event"
	dnotification "github.com/omerbeden/paymentgateway/internal/domain/notification"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/logger"
	infrakafka "github.com/omerbeden/paymentgateway/internal/infrastructure/queue/kafka"
	"github.com/omerbeden/paymentgateway/internal/usecase/notificaiton"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	appConfig := config.Load()
	var log logger.Logger

	if appConfig.Environment == "development" {
		log = logger.NewDevelopment()
	} else {
		log = logger.New(appConfig.LogLevel)
	}

	log.Info("Starting Notification Consumer Service...")

	infraConsumer, err := infrakafka.NewConsumer(*appConfig.Kafka)
	if err != nil {
		log.Fatal("kafka consumer: %v", err)
	}
	defer infraConsumer.Close()

	kafkaConsumer := messaging.NewKafkaConsumer(infraConsumer)

	var sender dnotification.Sender //change actually to a real sender implementation
	sendNotificationUC := notificaiton.NewSendPaymentNotificationUseCase(sender)
	notificationConsumer := consumer.NewNotificationEventConsumer(sendNotificationUC, log)

	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"notification-consumer"}`))
	})

	healthServer := &http.Server{Addr: ":8081", Handler: healthMux}
	go func() {
		log.Info("Health check server listening on :8081")
		if err := healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("health server error: %v", err)
		}
	}()

	log.Info("Subscribing to notification topics...")
	topics := []string{event.TopicNotificationPaymentCompleted}

	go func() {
		if err := kafkaConsumer.Subscribe(ctx, topics, notificationConsumer.Handle); err != nil {
			log.Error("kafka consumer error: %v", err)
		}
	}()

	log.Info("Notification consumer started, listening on topics: %v", topics)

	<-ctx.Done()
	log.Info("Shutting down notification consumer...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	healthServer.Shutdown(shutdownCtx)

	log.Info("Notification consumer stopped")

}
