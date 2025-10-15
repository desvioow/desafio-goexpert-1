package impl

import (
	"context"
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

	err := strategy.client.Incr(ctx, key).Err()
	if err != nil {
		return false, err
	}
	err = strategy.client.Expire(ctx, key, time.Second).Err()
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

func (strategy *RedisPersistenceStrategy) Delete(key string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return strategy.client.Del(ctx, key).Result()
}
