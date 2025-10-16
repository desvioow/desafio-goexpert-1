package strategy

import "time"

type PersistenceStrategyInterface interface {
	Connect() error
	Disconnect() error
	Persist(key string) (bool, error)
	ExpiresAt(key string, expiresAt time.Time) (bool, error)
	Get(key string) (interface{}, error)
	GetExpiration(key string) (time.Duration, error)
	Delete(key string) (interface{}, error)
}
