package middleware

import (
	"desafio-goexpert-1/internal/config"
	"desafio-goexpert-1/internal/limiter"
	"desafio-goexpert-1/internal/strategy"
	"log"
	"net/http"
	"strconv"
)

func RateLimiterMiddleware(strategy strategy.PersistenceStrategyInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limited, err := checkRequestLimit(w, r, strategy)
			if err != nil {
				log.Printf("Error checking request limit: %s", err)
				return
			}
			if limited {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func checkRequestLimit(w http.ResponseWriter, r *http.Request, strategy strategy.PersistenceStrategyInterface) (bool, error) {
	l := limiter.RateLimiter{}
	l.NewRateLimiter(strategy)

	limited, err := l.CheckLimit(w, r)
	if err != nil {
		return false, err
	}
	if limited {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", strconv.Itoa(config.AppConfig.RetryAfterSeconds))
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"}`))
		return true, nil
	}

	return false, nil
}
