package limiter

import (
	"desafio-goexpert-1/internal/config"
	"desafio-goexpert-1/internal/strategy"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	if limitedByToken {
		return true, nil
	}

	limitedByIP, err := checkIpLimit(r, rateLimiter)
	if err != nil {
		log.Printf("Error checking IP limit: %s", err)
		return false, err
	}

	return limitedByIP, err
}

func checkTokenLimit(r *http.Request, rateLimiter *RateLimiter) (bool, error) {

	tokenLimit := config.AppConfig.TokenLimitPerSecond

	if token := r.Header.Get("API_KEY"); token != "" {
		success, err := rateLimiter.strategy.Persist(token)
		if err != nil {
			return false, err
		}
		if success {
			tries, err := rateLimiter.strategy.Get(token)
			var triesInt int

			if triesStr, ok := tries.(string); ok {
				triesInt, err = strconv.Atoi(triesStr)
				if err != nil {
					return false, fmt.Errorf("failed to convert tries to int: %v", err)
				}
			} else {
				return false, fmt.Errorf("unexpected type for tries: %T", tries)
			}

			if triesInt >= tokenLimit {
				log.Printf("Request with token: %s exceeded limit", token)
				return true, nil
			}
		}
	} else {
		log.Printf("Token not found in request header.")
	}

	return false, nil
}

func checkIpLimit(r *http.Request, rateLimiter *RateLimiter) (bool, error) {

	ipLimit := config.AppConfig.IPLimitPerSecond
	requestIp := r.RemoteAddr

	if requestIp != "" {
		success, err := rateLimiter.strategy.Persist(requestIp)
		if err != nil {
			return false, err
		}
		if success {
			tries, err := rateLimiter.strategy.Get(requestIp)
			if err != nil {
				return false, err
			}
			var triesInt int

			if triesStr, ok := tries.(string); ok {
				triesInt, err = strconv.Atoi(triesStr)
				if err != nil {
					return false, fmt.Errorf("failed to convert tries to int: %v", err)
				}
			} else {
				return false, fmt.Errorf("unexpected type for tries: %T", tries)
			}

			if triesInt >= ipLimit {
				log.Printf("Request with ip: %s exceeded limit", requestIp)
				return true, nil
			}
		}
	} else {
		log.Printf("IP not found in request.")
	}

	return false, nil
}
