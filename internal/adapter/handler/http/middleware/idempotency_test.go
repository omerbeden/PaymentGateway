package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMockRedis(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	t.Cleanup(func() {
		client.Close()
		mr.Close()
	})

	return client
}
func TestIdempotencyMW_First_Request_Creates_New(t *testing.T) {
	gin.SetMode(gin.TestMode)
	redis := setupMockRedis(t)
	mw := NewIdempotancyMiddleware(redis)

	router := gin.New()

	router.POST("/test", mw.Check(), func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"id": "pay_123"})
	})

	body := []byte(`{"amount:1000}`)
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req.Header.Set("X-Idempotency-Key", "test-key-123")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "pay_123", response["id"])
}
func TestIdempotencyMW_Return_Cached(t *testing.T) {

	gin.SetMode(gin.TestMode)
	redis := setupMockRedis(t)
	mw := NewIdempotancyMiddleware(redis)

	router := gin.New()

	callCount := 0
	router.POST("/test", mw.Check(), func(c *gin.Context) {
		callCount++
		c.JSON(http.StatusCreated, gin.H{
			"id":     fmt.Sprintf("pay_%d", callCount),
			"amount": 10000,
		})
	})

	body := []byte(`{"amount:1000}`)
	req1 := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req1.Header.Set("X-Idempotency-Key", "test-key-123")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	req2 := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req2.Header.Set("X-Idempotency-Key", "test-key-123")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusCreated, w1.Code)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, w1.Body.String(), w2.Body.String())
}

func TestIdempotencyMW_DifferentKeys_Same_Body_Should_Create_New(t *testing.T) {
	gin.SetMode(gin.TestMode)
	redis := setupMockRedis(t)
	mw := NewIdempotancyMiddleware(redis)

	router := gin.New()
	callCount := 0
	router.POST("/test", mw.Check(), func(c *gin.Context) {
		callCount++
		c.JSON(http.StatusCreated, gin.H{
			"id":     fmt.Sprintf("pay_%d", callCount),
			"amount": 10000,
		})
	})

	body := []byte(`{"amount:1000}`)
	req1 := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req1.Header.Set("X-Idempotency-Key", "test-key-1")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	req2 := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req2.Header.Set("X-Idempotency-Key", "test-key-2")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusCreated, w1.Code)
	assert.Equal(t, http.StatusCreated, w2.Code)
	assert.NotEqual(t, w1.Body.String(), w2.Body.String())
}
func TestIdempotencyMW_Generate_Key_Different_Body_Creates_New(t *testing.T) {

	gin.SetMode(gin.TestMode)
	redis := setupMockRedis(t)
	mw := NewIdempotancyMiddleware(redis)

	router := gin.New()

	callCount := 0
	router.POST("/test", mw.Check(), func(c *gin.Context) {
		callCount++
		c.JSON(http.StatusCreated, gin.H{
			"id":     fmt.Sprintf("pay_%d", callCount),
			"amount": 10000,
		})
	})

	body := []byte(`{"amount:1}`)

	req1 := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	w1 := httptest.NewRecorder()

	router.ServeHTTP(w1, req1)

	body = []byte(`{"amount:2}`)
	req2 := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusCreated, w1.Code)
	assert.Equal(t, http.StatusCreated, w2.Code)

}

func TestIdempotencyMW_Generate_Key_Should_Return_Cache_For_Second_Request(t *testing.T) {

	gin.SetMode(gin.TestMode)
	redis := setupMockRedis(t)
	mw := NewIdempotancyMiddleware(redis)

	router := gin.New()

	callCount := 0
	router.POST("/test", mw.Check(), func(c *gin.Context) {
		callCount++
		c.JSON(http.StatusCreated, gin.H{
			"id":     fmt.Sprintf("pay_%d", callCount),
			"amount": 10000,
		})
	})

	body := []byte(`{"amount:1000}`)
	req1 := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	w1 := httptest.NewRecorder()

	router.ServeHTTP(w1, req1)

	req2 := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusCreated, w1.Code)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, w1.Body.String(), w2.Body.String())

}
