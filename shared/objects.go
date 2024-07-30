package shared

import (
	"time"
)

// Migration represents a single migration.
type Migration struct {
	FilePath    string `yaml:"file_path"`
	Package     string `yaml:"package"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Up          func(DatabaseStrategy) error
	Down        func(DatabaseStrategy) error
}

// AppliedMigration represents a migration that has been applied to the database.
type AppliedMigration struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AppliedOn   time.Time `json:"applied_on"`
}

// Manifest strongly typed respresentation of the manifest file.
type MigratorConfig struct {
	ManifestPath         string         `yaml:"manifest_path"`
	GeneratorPath        string         `yaml:"generator_path"`
	GeneratorPackage     string         `yaml:"generator_package"`
	ImplementationName   string         `yaml:"implementation_name"`
	DatabaseStrategyName string         `yaml:"database_strategy_name"`
	DatabaseStrategy     DatabaseConfig `yaml:"database_strategy"`

	Migrations []Migration `yaml:"migrations"`
}

// DatabaseConfig contains the configuration for the database to be used by the migration system.
type DatabaseConfig struct {
	Name      string `yaml:"name"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Database  string `yaml:"database"`
	Namespace string `yaml:"namespace"`
}

// DatabaseStrategy defines the interface for a database strategy.
type DatabaseStrategy struct {
	Name                  string
	DBConfig              DatabaseConfig
	EnsureInfrastructure  func(DatabaseConfig) error
	ApplyMigration        func(DatabaseConfig, string, string) error
	FindAppliedMigrations func(DatabaseConfig) ([]AppliedMigration, error)
	RollbackMigration     func(DatabaseConfig, string) error
	ResetMigrations       func(DatabaseConfig) error
	Exec                  func(DatabaseConfig, string, map[string]interface{}) error
}

// GetMigrationPath return the path that migrations will be stored based on properties in the manifest object.
func (manifest *MigratorConfig) GetMigrationPath() string {
	// return path.Join(manifest.ProjectRoot, manifest.GeneratorPath)
	return manifest.GeneratorPath
}

func GetTestConfig() MigratorConfig {
	config := MigratorConfig{
		GeneratorPath:        "migrations",
		GeneratorPackage:     "migrations",
		DatabaseStrategyName: "surrealdb",
	}

	return config
}
