package limiter

import (
	"desafio-goexpert-1/internal/config"
	"desafio-goexpert-1/internal/strategy"
	"log"
	"net/http"
)

func (rateLimiter *RateLimiter) NewRateLimiter(strategy strategy.PersistenceStrategyInterface) *RateLimiter {
	rateLimiter.strategy = strategy
	return rateLimiter
}

func (rateLimiter *RateLimiter) CheckLimit(w http.ResponseWriter, r *http.Request) (bool, error) {

	limitedByToken, err := checkTokenLimit(r, rateLimiter)
	if err != nil {
		log.Printf("Error checking Token limit: %s", err)
		return false, err
	}

	limitedByIP, err := checkIpLimit(r, rateLimiter)
	if err != nil {
		log.Printf("Error checking IP limit: %s", err)
		return false, err
	}

	return limitedByIP || limitedByToken, err
}

func checkTokenLimit(r *http.Request, rateLimiter *RateLimiter) (bool, error) {

	tokenLimit := config.AppConfig.TokenLimitPerSecond

	if token := r.Header.Get("API_KEY"); token != "" {
		return rateLimiter.strategy.IncrAndCheckLimit(token, tokenLimit)
	}

	log.Printf("Token not found in request header.")
	return false, nil
}

func checkIpLimit(r *http.Request, rateLimiter *RateLimiter) (bool, error) {

	ipLimit := config.AppConfig.IPLimitPerSecond
	requestIp := r.RemoteAddr

	if requestIp != "" {
		return rateLimiter.strategy.IncrAndCheckLimit(requestIp, ipLimit)
	}

	log.Printf("IP not found in request.")
	return false, nil
}
