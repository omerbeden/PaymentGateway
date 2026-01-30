package config

import (
	"os"
)

type Config struct {
	Environment string
	DatabaseDSN string
	RedisAddr   string
	ServerPort  string
	Paypal      *Paypal
}

type Paypal struct {
	Enabled      bool
	BaseURL      string
	SandBoxURL   string
	ClientID     string
	ClientSecret string
}

func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		DatabaseDSN: getEnv("DATABASE_DSN", "postgresql://postgres:postgres@localhost:5432/payment_gateway?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Paypal: &Paypal{
			Enabled:      getEnv("PAYPAL_ENABLED", "true") == "true",
			BaseURL:      getEnv("PAYPAL_BASE_URL", "base_url"),
			SandBoxURL:   getEnv("PAYPAL_SANDBOX_URL", "https://api-m.sandbox.paypal.com"),
			ClientID:     getEnv("PAYPAL_CLIENT_ID", "client_id"),
			ClientSecret: getEnv("PAYPAL_CLIENT_SECRET", "client_secret"),
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
