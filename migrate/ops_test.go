package migrate

import (
	"fmt"
	"testing"

	"github.com/cscoding21/csmig/persistence"
	"github.com/cscoding21/csmig/shared"
)

func TestFindDiscoveredMigrations(t *testing.T) {
	manifest := shared.LoadManifest("../generate/migrations/.csmig.yaml")
	migrations := FindDiscoveredMigrationFiles(manifest)

	if len(migrations) == 0 {
		t.Log("no discovered migrations found.  this may be an error, but not necessarily")
	} else {
		t.Log("found discovered migrations ", len(migrations))
	}
}

func TestEnsureInfrastructure(t *testing.T) {
	manifest := shared.LoadManifest("../generate/migrations/.csmig.yaml")
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
	manifest := shared.LoadManifest("../generate/migrations/.csmig.yaml")
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
	manifest := shared.LoadManifest("../generate/migrations/.csmig.yaml")
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
	manifest := shared.LoadManifest("../generate/migrations/.csmig.yaml")
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
