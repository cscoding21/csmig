package generate

import (
	"testing"

	"github.com/cscoding21/csmig/shared"
)

func TestNewMigration(t *testing.T) {
	manifest := shared.GetManifestPath(".csmig.yaml")
	err := NewMigration(manifest, "test migration")
	if err != nil {
		t.Error(err)
	}
}

func TestWriteCatalogFile(t *testing.T) {
	manifest := shared.LoadManifest(".csmig.yaml")
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

	contents := writeMigrationFile(migration)
	if len(contents) == 0 {
		t.Error("contents should have returned a value")
	}
}
