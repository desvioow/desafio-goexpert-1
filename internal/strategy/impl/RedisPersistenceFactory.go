package impl

import "desafio-goexpert-1/internal/strategy"

type RedisPersistenceFactory struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func (factory *RedisPersistenceFactory) CreateStrategy() strategy.PersistenceStrategyInterface {
	return &RedisPersistenceStrategy{
		Host:     factory.Host,
		Port:     factory.Port,
		Password: factory.Password,
		DB:       factory.DB,
	}
}
