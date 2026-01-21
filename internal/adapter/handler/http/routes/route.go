package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/omerbeden/paymentgateway/internal/adapter/handler/http/handler"
	"github.com/redis/go-redis/v9"
)

func SetupRoutes(db *sql.DB, redis *redis.Client) *gin.Engine {
	r := gin.New()

	healthHandler := handler.NewHealthHandler(db, redis)

	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)
	return r
}
