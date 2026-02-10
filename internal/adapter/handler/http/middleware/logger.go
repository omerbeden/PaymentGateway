package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/logger"
)

func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		requestID, _ := c.Get("request_id")

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Info("HTTP Request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", path,
			"status", statusCode,
			"latency", latency.String(),
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				log.Error("Request error",
					"request_id", requestID,
					"error", err.Error(),
				)
			}
		}
	}
}
