package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/metrics"
)

func Metrics(m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		m.HTTPRequestTotal.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		m.HTTPRequestDuration.WithLabelValues(
			c.Request.Method,
			path,
		).Observe(duration)
	}
}
