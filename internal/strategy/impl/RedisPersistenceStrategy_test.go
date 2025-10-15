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

	testKey := "test_key"
	ctx := context.Background()
	strategy.client.Del(ctx, testKey)
	defer strategy.client.Del(ctx, testKey)

	success, err := strategy.Persist(testKey)
	assert.True(t, success)
	assert.NoError(t, err)

	val, err := strategy.client.Get(ctx, testKey).Result()
	assert.NoError(t, err)
	assert.Equal(t, "1", val)

	success, err = strategy.Persist(testKey)
	assert.True(t, success)
	assert.NoError(t, err)

	val, err = strategy.client.Get(ctx, testKey).Result()
	assert.NoError(t, err)
	assert.Equal(t, "2", val)
}

func TestGet(t *testing.T) {
	strategy := setupRedis(t)
	ctx := context.Background()

	tests := []struct {
		name          string
		key           string
		setupValue    string
		expectedValue interface{}
		expectError   bool
	}{
		{
			name:          "existing_key",
			key:           "test_get_key1",
			setupValue:    "test_value1",
			expectedValue: "test_value1",
			expectError:   false,
		},
		{
			name:          "existing_numeric_key",
			key:           "test_get_key2",
			setupValue:    "42",
			expectedValue: "42",
			expectError:   false,
		},
		{
			name:          "non_existent_key",
			key:           "non_existent_key",
			setupValue:    "",
			expectedValue: nil,
			expectError:   true,
		},
	}

	defer func() {
		for _, tc := range tests {
			strategy.client.Del(ctx, tc.key)
		}
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupValue != "" {
				strategy.client.Set(ctx, tc.key, tc.setupValue, time.Minute)
			}

			result, err := strategy.Get(tc.key)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, result)
			}
		})
	}
}
