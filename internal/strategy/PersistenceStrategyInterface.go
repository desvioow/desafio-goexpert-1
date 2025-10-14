package strategy

import "time"

type PersistenceStrategyInterface interface {
	Connect() error
	Disconnect() error
	Persist(key string, value interface{}, expiration time.Duration) (bool, error)
	Get(key string) (interface{}, error)
}
