package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	NORMAL_TOKEN = "NORMAL_TOKEN"
	ULTRA_TOKEN  = "ULTRA_TOKEN"
)

var validTokens = []string{NORMAL_TOKEN, ULTRA_TOKEN}

type TokenConfig struct {
	Token string
	Limit int
}

type Config struct {
	IPLimitPerSecond           int
	TokenConfigs               []TokenConfig
	DefaultTokenLimitPerSecond int
	RetryAfterSeconds          int
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	ipLimit := safeParseInt(
		getEnv("IP_LIMIT_PER_SECOND", "10"),
		10,
		"Invalid IP_LIMIT_PER_SECOND",
	)

	retryAfter := safeParseInt(
		getEnv("RETRY_AFTER_SECONDS", "60"),
		60,
		"Invalid RETRY_AFTER_SECONDS",
	)

	defaultTokenLimit := safeParseInt(
		getEnv("DEFAULT_TOKEN_RATE_PER_SECOND", "20"),
		20,
		"Invalid DEFAULT_TOKEN_RATE_PER_SECOND",
	)

	tokens := []TokenConfig{
		{
			Token: NORMAL_TOKEN,
			Limit: safeParseInt(
				getEnv("BASIC_TOKEN_RATE_PER_SECOND", "50"),
				50,
				"Invalid BASIC_TOKEN_RATE_PER_SECOND",
			),
		},
		{
			Token: ULTRA_TOKEN,
			Limit: safeParseInt(
				getEnv("ULTRA_TOKEN_RATE_PER_SECOND", "100"),
				100,
				"Invalid ULTRA_TOKEN_RATE_PER_SECOND",
			),
		},
	}

	AppConfig = &Config{
		IPLimitPerSecond:           ipLimit,
		DefaultTokenLimitPerSecond: defaultTokenLimit,
		RetryAfterSeconds:          retryAfter,
		TokenConfigs:               tokens,
	}
}

func safeParseInt(value string, defaultValue int, errorMessage string) int {

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("%s, using default: %d", errorMessage, defaultValue)
		return defaultValue
	}
	return parsedValue
}

func getEnv(key string, defaultValue string) string {

	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetTokenLimit(token string) int {

	if token == "" {
		return -1
	}
	for _, tc := range c.TokenConfigs {
		if tc.Token == token {
			return tc.Limit
		}
	}
	return c.DefaultTokenLimitPerSecond
}
