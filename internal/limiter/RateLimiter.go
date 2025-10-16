package limiter

import (
	"desafio-goexpert-1/internal/config"
	"desafio-goexpert-1/internal/strategy"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

func (rateLimiter *RateLimiter) NewRateLimiter(strategy strategy.PersistenceStrategyInterface) *RateLimiter {
	rateLimiter.strategy = strategy
	return rateLimiter
}

func (rateLimiter *RateLimiter) CheckLimit(w http.ResponseWriter, r *http.Request) (bool, error) {

	token := r.Header.Get("API_KEY")
	if token != "" {
		return checkTokenLimit(r, rateLimiter)
	}

	return checkIpLimit(r, rateLimiter)
}

func checkTokenLimit(r *http.Request, rateLimiter *RateLimiter) (bool, error) {

	token := r.Header.Get("API_KEY")
	tokenLimit := config.AppConfig.GetTokenLimit(token)

	if tokenLimit > 0 {
		currentCount, err := getCurrentTries(token, rateLimiter.strategy)
		if err != nil {
			return false, err
		}

		if currentCount >= tokenLimit {
			return true, nil
		}

		_, err = rateLimiter.strategy.Persist(token)
		if err != nil {
			return false, err

		}
	}

	return false, nil
}

func checkIpLimit(r *http.Request, rateLimiter *RateLimiter) (bool, error) {

	ipLimit := config.AppConfig.IPLimitPerSecond
	requestIp := strings.Split(r.RemoteAddr, ":")[0]

	currentCount, err := getCurrentTries(requestIp, rateLimiter.strategy)
	if err != nil {
		return false, err
	}

	if currentCount >= ipLimit {
		return true, nil
	}

	if requestIp != "" {
		_, err := rateLimiter.strategy.Persist(requestIp)
		if err != nil {
			return false, err
		}
	}

	return false, nil
}

func getCurrentTries(key string, strategy strategy.PersistenceStrategyInterface) (int, error) {

	tries, err := strategy.Get(key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	if triesStr, ok := tries.(string); ok {
		triesInt, err := strconv.Atoi(triesStr)
		if err != nil {
			return 0, err
		}
		return triesInt, err
	}

	return 0, nil
}
