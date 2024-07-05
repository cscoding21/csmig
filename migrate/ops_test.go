package migrate

import (
	"fmt"
	"testing"

	"github.com/cscoding21/csmig/persistence"
	"github.com/cscoding21/csmig/shared"
)

func TestFindDiscoveredMigrations(t *testing.T) {
	manifest := shared.LoadManifest("../migrations/.csmig.yaml")
	migrations := FindDiscoveredMigrationFiles(manifest)

	if len(migrations) == 0 {
		t.Error("migrations should have returned a positive length")
	}
}

func TestEnsureInfrastructure(t *testing.T) {
	manifest := shared.LoadManifest("../migrations/.csmig.yaml")
	strategy, err := persistence.GetPersistenceStrategy(manifest.VersionStrategy)
	if err != nil {
		t.Error(err)
	}

	err = EnsureInfrastructure(strategy)
	if err != nil {
		t.Error(err)
	}
}

func TestApplyMigration(t *testing.T) {
	manifest := shared.LoadManifest("../migrations/.csmig.yaml")
	strategy, err := persistence.GetPersistenceStrategy(manifest.VersionStrategy)
	if err != nil {
		t.Error(err)
	}

	err = ApplyMigration(strategy, "m123", "unit test migration")
	if err != nil {
		t.Error(err)
	}
}

func TestRollbackMigration(t *testing.T) {
	manifest := shared.LoadManifest("../migrations/.csmig.yaml")
	strategy, err := persistence.GetPersistenceStrategy(manifest.VersionStrategy)
	if err != nil {
		t.Error(err)
	}

	err = RollbackMigration(strategy, "m123")
	if err != nil {
		t.Error(err)
	}
}

func TestFindApplyMigrations(t *testing.T) {
	manifest := shared.LoadManifest(".../migrations/csmig.yaml")
	strategy, err := persistence.GetPersistenceStrategy(manifest.VersionStrategy)
	if err != nil {
		t.Error(err)
	}

	am, err := FindAppliedMigrations(strategy)
	if err != nil {
		t.Error(err)
	}

	for _, m := range am {
		fmt.Println(m.Name, m.AppliedOn)
	}
}
