package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment string
	DatabaseDSN string
	RedisAddr   string
	ServerPort  string
	Paypal      *Paypal
	LogLevel    string
	Kafka       *Kafka
	Mongo       *Mongo
}

type Paypal struct {
	Enabled      bool
	BaseURL      string
	SandBoxURL   string
	ClientID     string
	ClientSecret string
	WebhookID    string
}

type Kafka struct {
	Brokers         string
	FlushTimeoutMs  int
	SASLUsername    string
	SASLPassword    string
	SASLMechanism   string
	TLSEnabled      bool
	AutoOffsetReset string
}

type Mongo struct {
	URI      string
	Timeout  time.Duration
	Database string
}

func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		DatabaseDSN: getEnv("DATABASE_DSN", "postgresql://postgres:postgres@localhost:5432/payment_gateway?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "development"),
		Paypal: &Paypal{
			Enabled:      getEnv("PAYPAL_ENABLED", "true") == "true",
			BaseURL:      getEnv("PAYPAL_BASE_URL", "base_url"),
			SandBoxURL:   getEnv("PAYPAL_SANDBOX_URL", "https://api-m.sandbox.paypal.com"),
			ClientID:     getEnv("PAYPAL_CLIENT_ID", "client_id"),
			ClientSecret: getEnv("PAYPAL_CLIENT_SECRET", "client_secret"),
		},
		Kafka: &Kafka{
			Brokers:        getEnv("KAFKA_BROKERS", "localhost:9092"),
			FlushTimeoutMs: getEnvInt("KAFKA_FLUSH_TIMEOUT_MS", 5000),
			SASLUsername:   getEnv("KAFKA_SASL_USERNAME", ""),
			SASLPassword:   getEnv("KAFKA_SASL_PASSWORD", ""),
			SASLMechanism:  getEnv("KAFKA_SASL_MECHANISM", "PLAIN"),
			TLSEnabled:     getEnvBool("KAFKA_TLS_ENABLED", false),
		},
		Mongo: &Mongo{
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Timeout:  getEnvDuration("MONGO_TIMEOUT", 10*time.Second),
			Database: getEnv("MONGO_DATABASE", "payment_gateway"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
