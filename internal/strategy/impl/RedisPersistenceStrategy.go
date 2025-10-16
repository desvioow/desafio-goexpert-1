package impl

import (
	"context"
	"desafio-goexpert-1/internal/config"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisPersistenceStrategy struct {
	client   *redis.Client
	Host     string
	Port     int
	Password string
	DB       int
}

func (strategy *RedisPersistenceStrategy) NewRedisPersistenceStrategy() {
	strategy.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", strategy.Host, strategy.Port),
		Password: strategy.Password,
		DB:       strategy.DB,
	})
}

func (strategy *RedisPersistenceStrategy) Connect() error {
	strategy.NewRedisPersistenceStrategy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := strategy.client.Ping(ctx).Result()
	return err
}

func (strategy *RedisPersistenceStrategy) Disconnect() error {
	if strategy.client != nil {
		return strategy.client.Close()
	}
	return nil
}

func (strategy *RedisPersistenceStrategy) Persist(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := strategy.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if count == 1 {
		expiresAt := time.Now().Add(time.Duration(config.AppConfig.RateWindow) * time.Second)
		expirationErr := strategy.client.ExpireAt(ctx, key, expiresAt).Err()
		if expirationErr != nil {
			return false, expirationErr
		}
	}

	return true, err
}

func (strategy *RedisPersistenceStrategy) ExpiresAt(key string, expiration time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := strategy.client.ExpireAt(ctx, key, expiration).Err()
	if err != nil {
		return false, err
	}
	return true, err
}

func (strategy *RedisPersistenceStrategy) Get(key string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return strategy.client.Get(ctx, key).Result()
}

func (strategy *RedisPersistenceStrategy) GetExpiration(key string) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return strategy.client.TTL(ctx, key).Result()
}

func (strategy *RedisPersistenceStrategy) Delete(key string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return strategy.client.Del(ctx, key).Result()
}
