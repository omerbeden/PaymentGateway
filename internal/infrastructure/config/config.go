package config

import "os"

type Config struct {
	Environment string
	DatabaseDSN string
	RedisAddr   string
	ServerPort  string
}

func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		DatabaseDSN: getEnv("DATABASE_DSN", "postgresql://postgres:postgres@localhost:5432/payment_gateway?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
