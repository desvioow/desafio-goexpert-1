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

type RateLimiterConfig struct {
	RateWindow                  int
	IPLimitPerSecond            int
	TokenConfigs                []TokenConfig
	FallbackTokenLimitPerSecond int
	RetryAfterSeconds           int
}

type DBConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

var AppConfig *RateLimiterConfig
var RedisConfig *DBConfig

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	loadRateLimiterConfig()
	loadRedisConfig()
}

func loadRedisConfig() {

	host := getEnv("REDIS_HOST", "localhost")
	port := safeParseInt(getEnv("REDIS_PORT", "6379"), 6379, "Invalid REDIS_PORT")
	password := getEnv("REDIS_PASSWORD", "")
	db := safeParseInt(getEnv("REDIS_DB", "0"), 0, "Invalid REDIS_DB")

	RedisConfig = &DBConfig{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
	}
}

func loadRateLimiterConfig() {
	rateWindow := safeParseInt(
		getEnv("RATE_WINDOW", "1"),
		1, "Invalid RATE_WINDOW")

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

	fallbackTokenLimit := safeParseInt(
		getEnv("UNKNOWN_TOKEN_RATE_PER_SECOND", "20"),
		20,
		"Invalid UNKNOWN_TOKEN_RATE_PER_SECOND",
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

	AppConfig = &RateLimiterConfig{
		RateWindow:                  rateWindow,
		IPLimitPerSecond:            ipLimit,
		FallbackTokenLimitPerSecond: fallbackTokenLimit,
		RetryAfterSeconds:           retryAfter,
		TokenConfigs:                tokens,
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

func (c *RateLimiterConfig) GetTokenLimit(token string) int {

	if token == "" {
		return -1
	}
	for _, tc := range c.TokenConfigs {
		if tc.Token == token {
			return tc.Limit
		}
	}
	return c.FallbackTokenLimitPerSecond
}
