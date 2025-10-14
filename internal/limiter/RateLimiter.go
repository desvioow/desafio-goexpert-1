package limiter

import (
	"desafio-goexpert-1/internal/config"
	"desafio-goexpert-1/internal/strategy"
	"log"
	"net/http"
	"strconv"
)

func (rateLimiter *RateLimiter) NewRateLimiter(strategy strategy.PersistenceStrategyInterface) *RateLimiter {
	return &RateLimiter{strategy: strategy}
}

func (rateLimiter *RateLimiter) CheckLimit(w http.ResponseWriter, r *http.Request) (bool, error) {

	limitedByToken, err := checkTokenLimit(r, rateLimiter)
	if err != nil {
		log.Printf(">>> Error checking Token limit: %s", err)
	}

	limitedByIP, err := checkIpLimit(r, rateLimiter)
	if err != nil {
		log.Printf(">>> Error checking IP limit: %s", err)
	}

	return limitedByIP || limitedByToken, err
}

func checkTokenLimit(r *http.Request, rateLimiter *RateLimiter) (bool, error) {

	tokenLimit := config.AppConfig.TokenLimitPerSecond

	if token := r.Header.Get("API_KEY"); token != "" {
		dbTokenTentatives, err := rateLimiter.strategy.Get(token)
		if err != nil {
			return false, err
		}

		tentatives, err := strconv.Atoi(dbTokenTentatives.(string))
		if err != nil {
			return false, err
		}

		return tentatives >= tokenLimit, nil
	}

	log.Printf("Token not found in request header.")
	return false, nil
}

func checkIpLimit(r *http.Request, rateLimiter *RateLimiter) (bool, error) {

	ipLimit := config.AppConfig.IPLimitPerSecond
	requestIp := r.RemoteAddr

	if requestIp != "" {
		dbIpTentatives, err := rateLimiter.strategy.Get(requestIp)
		if err != nil {
			return false, err
		}

		tentatives, err := strconv.Atoi(dbIpTentatives.(string))
		if err != nil {
			return false, err
		}

		return tentatives >= ipLimit, nil
	}

	log.Printf("IP not found in request.")
	return true, nil
}
