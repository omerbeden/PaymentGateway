package routes

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	handler "github.com/omerbeden/paymentgateway/internal/adapter/handler/http"
	"github.com/omerbeden/paymentgateway/internal/adapter/handler/http/middleware"
	"github.com/omerbeden/paymentgateway/internal/adapter/provider"
	"github.com/omerbeden/paymentgateway/internal/adapter/provider/paypal"
	"github.com/omerbeden/paymentgateway/internal/adapter/repository/postgres"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/logger"
	"github.com/omerbeden/paymentgateway/internal/usecase/payment"
	"github.com/omerbeden/paymentgateway/internal/usecase/webhook"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(db *sql.DB, redis *redis.Client, cfg *config.Config) *gin.Engine {
	r := gin.New()
	var log logger.Logger

	if cfg.Environment == "development" {
		log = logger.NewDevelopment()
	} else {
		log = logger.New(cfg.LogLevel)
	}

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(gin.Recovery())
	r.Use(middleware.Timeout(30 * time.Second))

	paymentRepository := postgres.NewPaymentRepository(db)

	providerFactory := provider.NewProviderFactory()
	if cfg.Paypal.Enabled {
		if cfg.Environment == "developmet" {
			cfg.Paypal.BaseURL = cfg.Paypal.SandBoxURL
		}
		providerFactory.RegisterProvider("paypal", paypal.NewProvider(*cfg.Paypal))
	}

	createPaymentUC := payment.NewCreatePaymentUseCase(paymentRepository, providerFactory, log)
	webhookUC := webhook.NewProcessWebHookUseCase(paymentRepository, providerFactory)

	healthHandler := handler.NewHealthHandler(db, redis)
	paymentHandler := handler.NewPaymentHandler(createPaymentUC)
	webhookHandler := handler.NewWebhookHandler(webhookUC)

	idempotancyMW := middleware.NewIdempotancyMiddleware(redis)

	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)

	v1 := r.Group("/api/v1")
	{
		payments := v1.Group("/payments")
		{
			payments.POST("/payments", idempotancyMW.Check(), paymentHandler.CreatePayment)
		}

		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("/paypal", webhookHandler.HandlePaypal)
		}
	}

	return r
}
