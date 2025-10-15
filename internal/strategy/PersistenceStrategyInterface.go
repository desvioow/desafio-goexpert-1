package strategy

type PersistenceStrategyInterface interface {
	Connect() error
	Disconnect() error
	Persist(key string) (bool, error)
	Get(key string) (interface{}, error)
	Delete(key string) (interface{}, error)
}
