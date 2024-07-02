package persistence

import (
	"errors"

	"github.com/cscoding21/csmig/shared"
)

const (
	VersionTableName = "csmig_versions"
)

// DatabaseConfig contains the configuration for the database to be used by the migration system.
type DatabaseConfig struct {
	Name      string
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	Namespace string
}

// DatabaseStrategy defines the interface for a database strategy.
type DatabaseStrategy struct {
	Name                  string
	DBConfig              DatabaseConfig
	EnsureInfrastructure  func(DatabaseConfig) error
	ApplyMigration        func(DatabaseConfig, string) error
	FindAppliedMigrations func(DatabaseConfig) ([]shared.AppliedMigration, error)
	RollbackMigration     func(DatabaseConfig, string) error
	ResetMigrations       func(DatabaseConfig) error
}

var persistenceStrategies = map[string]DatabaseStrategy{
	"surrealdb": SurrealDBStrategy,
}

// GetPersistenceStrategy returns the persistence strategy for the given name.
func GetPersistenceStrategy(name string) (DatabaseStrategy, error) {
	strategy, ok := persistenceStrategies[name]
	if !ok {
		return DatabaseStrategy{}, errors.New("unknown persistence strategy")
	}

	return strategy, nil
}
