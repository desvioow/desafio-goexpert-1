package strategy

type PersistenceStrategyInterface interface {
	Connect() error
	Disconnect() error
	Persist(key string, value interface{}) (bool, error)
	Get(key string) (interface{}, error)
}
