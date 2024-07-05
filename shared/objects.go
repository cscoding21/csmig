package shared

import (
	"path"
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
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	AppliedOn   time.Time `yaml:"applied_on"`
}

// Manifest strongly typed respresentation of the manifest file.
type Manifest struct {
	ManifestPath       string `yaml:"manifest_path"`
	ProjectRoot        string `yaml:"project_root"`
	GeneratorPath      string `yaml:"generator_path"`
	GeneratorPackage   string `yaml:"generator_package"`
	ImplementationName string `yaml:"implementation_name"`
	VersionStrategy    string `yaml:"version_strategy"`

	Migrations []Migration `yaml:"migrations"`
}

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
	ApplyMigration        func(DatabaseConfig, string, string) error
	FindAppliedMigrations func(DatabaseConfig) ([]AppliedMigration, error)
	RollbackMigration     func(DatabaseConfig, string) error
	ResetMigrations       func(DatabaseConfig) error
	Exec                  func(DatabaseConfig, string, map[string]interface{}) error
}

// GetMigrationPath return the path that migrations will be stored based on properties in the manifest object.
func (manifest *Manifest) GetMigrationPath() string {
	return path.Join(manifest.ProjectRoot, manifest.GeneratorPath)
}
