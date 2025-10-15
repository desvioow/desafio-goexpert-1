package middleware

import (
	"desafio-goexpert-1/internal/config"
	"desafio-goexpert-1/internal/limiter"
	"desafio-goexpert-1/internal/strategy"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func RateLimiterMiddleware(strategy strategy.PersistenceStrategyInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Hello from RateLimiterMiddleware")
			limited, err := checkRequestLimit(w, r, strategy)
			if err != nil {
				return
			}
			if limited {
				log.Printf("Rate limit exceeded for request: %s", r.URL.Path)
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
		w.Header().Set("Retry-After", strconv.Itoa(config.AppConfig.RetryAfter))
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"}`))
		return true, nil
	}

	return false, nil
}
