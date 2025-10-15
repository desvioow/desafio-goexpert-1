package limiter

import (
	"desafio-goexpert-1/internal/strategy/impl"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimited(t *testing.T) {

	redisFactory := &impl.RedisPersistenceFactory{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}
	redisStrategy := redisFactory.CreateStrategy()

	rateLimiter := RateLimiter{strategy: redisStrategy}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	limited, err := rateLimiter.CheckLimit(recorder, req)
	if err != nil {
		return
	}
	if limited {
		assert.Equal(t, true, limited)
	}
}
