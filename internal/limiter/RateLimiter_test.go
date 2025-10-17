package limiter

import (
	"desafio-goexpert-1/internal/config"
	"desafio-goexpert-1/internal/strategy/impl"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupRedisStrategy - Inicializa redisStrategy para testes
func setupRedisStrategy(t *testing.T) *impl.RedisPersistenceStrategy {
	host := "localhost"
	port := 6379
	password := ""
	db := 0

	redisStrategy := &impl.RedisPersistenceStrategy{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       db,
	}
	err := redisStrategy.Connect()
	require.NoError(t, err, "Failed to connect to Redis")
	return redisStrategy
}

// cleanupKeys - Limpa chaves do Redis
func cleanupKeys(strategy *impl.RedisPersistenceStrategy, keys ...string) {
	for _, key := range keys {
		_, err := strategy.Delete(key)
		if err != nil {
			return
		}
	}
}

// TestIPRateLimit - Teste de eficácia: Verifica se o limite por IP é respeitado
func TestIPRateLimit(t *testing.T) {
	config.Load()
	redisStrategy := setupRedisStrategy(t)
	defer redisStrategy.Disconnect()

	rateLimiter := RateLimiter{strategy: redisStrategy}

	// Limpa qualquer estado anterior
	testKey := "127.0.0.1"
	cleanupKeys(redisStrategy, testKey)
	defer cleanupKeys(redisStrategy, testKey)

	// Primeira requisição deve passar
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1:8080"
	recorder := httptest.NewRecorder()

	limited, err := rateLimiter.CheckLimit(recorder, req)
	assert.NoError(t, err)
	assert.False(t, limited, "First request should not be limited")

	// Requisições até o limite devem passar
	for i := 1; i < config.AppConfig.IPLimitPerSecond; i++ {
		limited, err = rateLimiter.CheckLimit(recorder, req)
		assert.NoError(t, err)
		assert.False(t, limited, fmt.Sprintf("Request %d should not be limited", i+1))
	}

	// Próxima requisição deve ser limitada
	limited, err = rateLimiter.CheckLimit(recorder, req)
	assert.NoError(t, err)
	assert.True(t, limited, "Request exceeding limit should be limited")
}

// TestTokenRateLimit - Teste de eficácia: Verifica se o limite por token é respeitado
func TestTokenRateLimit(t *testing.T) {
	config.Load()
	redisStrategy := setupRedisStrategy(t)
	defer redisStrategy.Disconnect()

	rateLimiter := RateLimiter{strategy: redisStrategy}

	// Teste com NORMAL_TOKEN
	normalToken := config.NORMAL_TOKEN
	cleanupKeys(redisStrategy, normalToken)
	defer cleanupKeys(redisStrategy, normalToken)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", normalToken)
	req.RemoteAddr = "192.168.1.1:8080"
	recorder := httptest.NewRecorder()

	normalLimit := config.AppConfig.GetTokenLimit(normalToken)

	// Requisições até o limite do token devem passar
	for i := 0; i < normalLimit; i++ {
		limited, err := rateLimiter.CheckLimit(recorder, req)
		assert.NoError(t, err)
		assert.False(t, limited, fmt.Sprintf("Token request %d should not be limited", i+1))
	}

	// Próxima requisição deve ser limitada
	limited, err := rateLimiter.CheckLimit(recorder, req)
	assert.NoError(t, err)
	assert.True(t, limited, "Request exceeding token limit should be limited")
}

// TestTokenPriorityOverIP - Teste de eficácia: Verifica se token tem prioridade sobre IP
func TestTokenPriorityOverIP(t *testing.T) {
	config.Load()
	redisStrategy := setupRedisStrategy(t)
	defer redisStrategy.Disconnect()

	rateLimiter := RateLimiter{strategy: redisStrategy}

	ultraToken := config.ULTRA_TOKEN
	testIP := "10.0.0.1"

	// Limpa qualquer estado anterior
	cleanupKeys(redisStrategy, ultraToken, testIP)
	defer cleanupKeys(redisStrategy, ultraToken, testIP)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", ultraToken)
	req.RemoteAddr = testIP + ":8080"
	recorder := httptest.NewRecorder()

	ultraLimit := config.AppConfig.GetTokenLimit(ultraToken)
	ipLimit := config.AppConfig.IPLimitPerSecond

	// ULTRA_TOKEN tem limite maior que IP, deve permitir mais requisições
	assert.Greater(t, ultraLimit, ipLimit, "Ultra token should have higher limit than IP")

	// Faz requisições até exceder o limite de IP mas não o de token
	requestsToMake := ipLimit + 10

	for i := 0; i < requestsToMake; i++ {
		limited, err := rateLimiter.CheckLimit(recorder, req)
		assert.NoError(t, err)
		assert.False(t, limited, fmt.Sprintf("Request %d with ultra token should not be limited", i+1))
	}
}
