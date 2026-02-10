package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type IdempotencyMiddleware struct {
	redis *redis.Client
	ttl   time.Duration
}

func NewIdempotancyMiddleware(redis *redis.Client) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{
		redis: redis,
		ttl:   24 * time.Hour,
	}
}

func (im *IdempotencyMiddleware) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.Next()
			return
		}

		idempotencyKey := c.GetHeader("X-Idempotency-Key")

		if idempotencyKey == "" {
			body, _ := io.ReadAll(c.Request.Body)
			idempotencyKey = im.generateKey(body)
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
		}

		ctx := context.Background()
		key := fmt.Sprintf("idempotency:%s", idempotencyKey)

		cachedResponse, err := im.redis.Get(ctx, key).Result()

		if err != nil {
			var response map[string]interface{}
			if err := json.Unmarshal([]byte(cachedResponse), &response); err == nil {
				c.JSON(http.StatusOK, response)
				c.Abort()
				return
			}
		}

		c.Set("idempotency_key", idempotencyKey)

		blw := &bodyLogWriter{body: []byte{}, ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		if c.Writer.Status() == http.StatusCreated || c.Writer.Status() == http.StatusOK {
			im.redis.Set(ctx, key, blw.body, im.ttl)
		}

	}
}

func (im *IdempotencyMiddleware) generateKey(body []byte) string {
	hash := sha256.Sum256(body)
	return hex.EncodeToString(hash[:])
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body []byte
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}
