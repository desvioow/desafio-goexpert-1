package limiter

import (
	"desafio-goexpert-1/internal/config"
	"desafio-goexpert-1/internal/strategy/impl"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimited(t *testing.T) {

	config.Load()

	redisFactory := &impl.RedisPersistenceFactory{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}
	redisStrategy := redisFactory.CreateStrategy()
	err := redisStrategy.Connect()
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisStrategy.Disconnect()

	rateLimiter := RateLimiter{strategy: redisStrategy}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "test-token")
	recorder := httptest.NewRecorder()

	limited, err := rateLimiter.CheckLimit(recorder, req)
	assert.NoError(t, err, "CheckLimit should not return an error")
	assert.False(t, limited, "First request should not be limited")

	for i := 0; i < config.AppConfig.IPLimitPerSecond+5; i++ {
		limited, err = rateLimiter.CheckLimit(recorder, req)
		if err != nil {
			t.Errorf("Error on request %d: %v", i, err)
			return
		}
	}

	assert.True(t, limited, "Should be limited after many requests")
}
