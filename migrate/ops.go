package migrate

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/cscoding21/csmig/persistence"
	"github.com/cscoding21/csmig/shared"
)

// EnsureInfrastructure create the migration table in the target DB if it doesn't exist.
func EnsureInfrastructure(strategy persistence.DatabaseStrategy) error {
	return strategy.EnsureInfrastructure(strategy.DBConfig)
}

// ApplyMigration record a migration as being applied in the database
func ApplyMigration(strategy persistence.DatabaseStrategy, name string) error {
	return strategy.ApplyMigration(strategy.DBConfig, name)
}

func FindAppliedMigrations(strategy persistence.DatabaseStrategy) ([]shared.AppliedMigration, error) {
	return strategy.FindAppliedMigrations(strategy.DBConfig)
}

func RollbackMigration(strategy persistence.DatabaseStrategy, name string) error {
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
