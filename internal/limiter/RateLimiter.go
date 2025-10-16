package limiter

import (
	"desafio-goexpert-1/internal/config"
	"desafio-goexpert-1/internal/strategy"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

func (rateLimiter *RateLimiter) NewRateLimiter(strategy strategy.PersistenceStrategyInterface) *RateLimiter {
	rateLimiter.strategy = strategy
	return rateLimiter
}

func (rateLimiter *RateLimiter) CheckLimit(w http.ResponseWriter, r *http.Request) (bool, error) {

	token := r.Header.Get("API_KEY")
	var key string
	var limit int

	if token != "" {
		key = token
		limit = config.AppConfig.GetTokenLimit(token)
	} else {
		key = strings.Split(r.RemoteAddr, ":")[0]
		limit = config.AppConfig.IPLimitPerSecond
	}

	currentCount, err := getCurrentCount(key, rateLimiter.strategy)
	if err != nil {
		log.Printf("Error getting current count: %s", err)
		return false, err
	}

	if currentCount >= limit {
		limited, err := hasKeyBeenLimited(key, rateLimiter.strategy)
		if err != nil {
			log.Printf("Error checking if key has been limited: %s", err)
			return false, err
		}
		if !limited {
			_, err := rateLimiter.strategy.ExpiresAt(key, time.Now().Add(time.Second*time.Duration(config.AppConfig.RetryAfterSeconds)))
			if err != nil {
				log.Printf("Error setting expiration: %s", err)
				return false, err
			}
			go func() {
				time.Sleep(time.Second * time.Duration(config.AppConfig.RetryAfterSeconds))
				log.Printf("GO Routine Deleting key: %s", key)
				_, err := rateLimiter.strategy.Delete(key)
				if err != nil {
					log.Printf("Error deleting key: %s", err)
				}
			}()
		}
		return true, nil
	}

	if key != "" {
		_, err = rateLimiter.strategy.Persist(key)
		if err != nil {
			return false, err
		}
	}

	return false, nil
}

func getCurrentCount(key string, strategy strategy.PersistenceStrategyInterface) (int, error) {

	count, err := strategy.Get(key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	if countStr, ok := count.(string); ok {
		countInt, err := strconv.Atoi(countStr)
		if err != nil {
			return 0, err
		}
		return countInt, err
	}

	return 0, nil
}

func hasKeyBeenLimited(key string, strategy strategy.PersistenceStrategyInterface) (bool, error) {

	exp, err := strategy.GetExpiration(key)
	if err != nil {
		return false, err
	}
	if exp > time.Duration(config.AppConfig.RateWindow)*time.Second {
		return true, nil
	} else {
		return false, nil
	}
}
