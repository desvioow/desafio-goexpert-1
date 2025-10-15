package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	IPLimitPerSecond    int
	TokenLimitPerSecond int
	RetryAfterSeconds   int
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	ipLimit, err := strconv.Atoi(getEnv("IP_LIMIT_PER_SECOND", "10"))
	if err != nil {
		log.Printf("Invalid IP_LIMIT_PER_SECOND, using default: %v", err)
		ipLimit = 10
	}

	tokenLimit, err := strconv.Atoi(getEnv("TOKEN_LIMIT_PER_SECOND", "100"))
	if err != nil {
		log.Printf("Invalid TOKEN_LIMIT_PER_SECOND, using default: %v", err)
		tokenLimit = 100
	}

	retryAfter, err := strconv.Atoi(getEnv("RETRY_AFTER_SECONDS", "60"))
	if err != nil {
		log.Printf("Invalid RETRY_AFTER_SECONDS, using default: %v", err)
		retryAfter = 60
	}

	AppConfig = &Config{
		IPLimitPerSecond:    ipLimit,
		TokenLimitPerSecond: tokenLimit,
		RetryAfterSeconds:   retryAfter,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
