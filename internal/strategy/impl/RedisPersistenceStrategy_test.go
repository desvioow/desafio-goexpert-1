package impl

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func setupRedis(t *testing.T) *RedisPersistenceStrategy {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		t.Fatalf("Unable to connect to Redis: %v", err)
	}

	strategy := &RedisPersistenceStrategy{
		client: client,
		Host:   "localhost",
		Port:   6379,
	}

	return strategy
}

func TestConnect(t *testing.T) {
	strategy := &RedisPersistenceStrategy{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	err := strategy.Connect()
	assert.NoError(t, err)
}

func TestDisconnect(t *testing.T) {
	strategy := setupRedis(t)

	err := strategy.Disconnect()
	assert.NoError(t, err)
}

func TestPersist(t *testing.T) {
	strategy := setupRedis(t)

	tests := []struct {
		name  string
		key   string
		value interface{}
		valid bool
	}{
		{"normal string", "key1", "value1", true},
		{"empty value", "key2", "", true},
		{"int value", "key3", 123, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			success, err := strategy.Persist(tc.key, tc.value)
			if tc.valid {
				assert.True(t, success)
				assert.NoError(t, err)
			} else {
				assert.False(t, success)
				assert.Error(t, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	strategy := setupRedis(t)

	// Set up initial data
	strategy.Persist("key1", "value1")
	strategy.Persist("key2", "value2")

	tests := []struct {
		name      string
		key       string
		expected  interface{}
		shouldErr bool
	}{
		{"existing key 1", "key1", "value1", false},
		{"existing key 2", "key2", "value2", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			value, err := strategy.Get(tc.key)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, value)
			}
		})
	}
}
