package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/omerbeden/paymentgateway/internal/adapter/handler/http/handler"
	"github.com/omerbeden/paymentgateway/internal/adapter/provider"
	"github.com/omerbeden/paymentgateway/internal/adapter/provider/paypal"
	"github.com/omerbeden/paymentgateway/internal/adapter/repository/postgres"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
	"github.com/omerbeden/paymentgateway/internal/usecase/payment"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(db *sql.DB, redis *redis.Client, cfg *config.Config) *gin.Engine {
	r := gin.New()

	paymentRepository := postgres.NewPaymentRepository(db)

	providerFactory := provider.NewProviderFactory()
	if cfg.Paypal.Enabled {
		providerFactory.RegisterProvider("paypal", paypal.NewProvider(cfg.Paypal))
	}
	createPaymentUC := payment.NewCreatePaymentUseCase(paymentRepository, providerFactory)

	healthHandler := handler.NewHealthHandler(db, redis)
	paymentHandler := handler.NewPaymentHandler(createPaymentUC)

	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)

	v1 := r.Group("/api/v1")
	{
		payments := v1.Group("/payments")
		{
			payments.POST("/payments", paymentHandler.CreatePayment)
		}
	}

	return r
}
