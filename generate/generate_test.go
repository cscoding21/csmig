package generate

import (
	"testing"

	"github.com/cscoding21/csmig/shared"
)

func TestWriteCatalogFile(t *testing.T) {
	manifest := shared.LoadManifest("migrations/.csmig.yaml")
	err := writeCatalogFile(manifest)
	if err != nil {
		t.Error(err)
	}
}

func TestWriteMigrationFile(t *testing.T) {
	//---create a catalog
	migration := shared.Migration{
		Package:     "migrations",
		Name:        "m123",
		Description: "This is a test migration",
	}

	contents := getMigrationFileContents(migration)
	if len(contents) == 0 {
		t.Error("contents should have returned a value")
	}
}

func TestNewMigration(t *testing.T) {
	manifest := shared.LoadManifest("migrations/.csmig.yaml")
	mig, err := NewMigration(manifest, "This is a test migration")
	if err != nil {
		t.Error("TestNewMigration failed: ", err)
	}

	t.Log("migration name: ", mig.Name)

	// err = RemoveMigration(manifest, mig.Name)
	// if err != nil {
	// 	t.Error("TestNewMigration failed: ", err)
	// }
}

func TestRemoveMigration(t *testing.T) {
	manifest := shared.LoadManifest("migrations/.csmig.yaml")

	//---ensure there is a migration to remove
	mig, err := NewMigration(manifest, "This is a test migration for integration testing")
	if err != nil {
		t.Error("unable to add new migration to rest removal")
	}

	err = RemoveMigration(manifest, mig.Name)
	if err != nil {
		t.Error("migration name should have been set: ", err)
	}
}

func TestInit(t *testing.T) {
	err := Init()
	if err != nil {
		t.Error(err)
	}
}
