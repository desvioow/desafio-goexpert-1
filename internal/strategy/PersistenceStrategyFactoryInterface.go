package strategy

type PersistenceStrategyFactory interface {
	CreateStrategy() PersistenceStrategyInterface
}
