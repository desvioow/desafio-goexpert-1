package limiter

import (
	"desafio-goexpert-1/internal/strategy"
	"net/http"
)

type RateLimiter struct {
	strategy strategy.PersistenceStrategyInterface
}

type RateLimiterInterface interface {
	CheckLimit(w http.ResponseWriter, r *http.Request) (bool, error)
	NewRateLimiter(strategy strategy.PersistenceStrategyInterface) *RateLimiter
}
