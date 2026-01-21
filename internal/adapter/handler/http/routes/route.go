package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/omerbeden/paymentgateway/internal/adapter/handler/http/handler"
)

func SetupRoutes() *gin.Engine {
	r := gin.New()

	healthHandler := handler.NewHealthHandler()

	r.GET("/health", healthHandler.Health)
	return r
}
