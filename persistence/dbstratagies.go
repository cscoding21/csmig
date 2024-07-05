package persistence

import (
	"errors"

	"github.com/cscoding21/csmig/shared"
)

const (
	VersionTableName = "csmig_versions"
)

var persistenceStrategies = map[string]shared.DatabaseStrategy{
	"surrealdb": SurrealDBStrategy,
}

// GetPersistenceStrategy returns the persistence strategy for the given name.
func GetPersistenceStrategy(name string) (shared.DatabaseStrategy, error) {
	strategy, ok := persistenceStrategies[name]
	if !ok {
		return shared.DatabaseStrategy{}, errors.New("unknown persistence strategy")
	}

	return strategy, nil
}
