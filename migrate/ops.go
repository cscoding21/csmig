package migrate

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/cscoding21/csmig/shared"
)

// EnsureInfrastructure create the migration table in the target DB if it doesn't exist.
func EnsureInfrastructure(strategy shared.DatabaseStrategy) error {
	return strategy.EnsureInfrastructure(strategy.DBConfig)
}

// ApplyMigration record a migration as being applied in the database
func ApplyMigration(strategy shared.DatabaseStrategy, name string, description string) error {
	return strategy.ApplyMigration(strategy.DBConfig, name, description)
}

func FindAppliedMigrations(strategy shared.DatabaseStrategy) ([]shared.AppliedMigration, error) {
	return strategy.FindAppliedMigrations(strategy.DBConfig)
}

func RollbackMigration(strategy shared.DatabaseStrategy, name string) error {
	return strategy.RollbackMigration(strategy.DBConfig, name)
}

// FindDiscoveredMigrationFiles iterated over files in the migratin path and return all created migrations
func FindDiscoveredMigrationFiles(manifest shared.Manifest) []shared.Migration {
	migrations := []shared.Migration{}

	files, err := filepath.Glob(path.Join(manifest.GetMigrationPath(), "/m*_gen.go"))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fn := filepath.Base(file)
		mn := strings.Replace(fn, "_gen.go", "", 1)
		migrations = append(migrations, shared.Migration{
			FilePath: file,
			Package:  manifest.GeneratorPackage,
			Name:     mn,
		})

		fmt.Println(mn)
	}

	return migrations
}
