package generate

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/cscoding21/csgen"
	"github.com/cscoding21/csmig/migrate"
	"github.com/cscoding21/csmig/shared"

	"gopkg.in/yaml.v3"
)

// NewMigration creates a new migration
func NewMigration(manifestPath string, description string) error {
	mp := shared.GetManifestPath(manifestPath)
	manifest := LoadManifest(mp)

	builder := csgen.NewCSGenBuilderForOneOffFile("csmig", manifest.GeneratorPackage)

	migrationName := getMigrationName()
	migration := shared.Migration{
		Package:     manifest.GeneratorPackage,
		Name:        migrationName,
		Description: description,
	}

	builder.WriteString(writeMigrationFile(migration))
	migrationFileName := fmt.Sprintf("%s_gen.go", migrationName)
	migrationFilePath := path.Join(manifest.ProjectRoot, manifest.GeneratorPath, migrationFileName)

	err := csgen.WriteGeneratedGoFile(migrationFilePath, builder.String())
	if err != nil {
		log.Fatal(err)
	}

	writeCatalogFile(manifest)

	return nil
}

func writeCatalogFile(manifest shared.Manifest) error {
	migrations := migrate.FindDiscoveredMigrations(manifest)
	contents := csgen.ExecuteTemplate("catalog", catalogTemplateString, migrations)

	builder := csgen.NewCSGenBuilderForFile("csmig", manifest.GeneratorPackage)
	builder.WriteString(contents)

	catalogPath := path.Join(manifest.ProjectRoot, manifest.GeneratorPath, "catalog.gen.go")
	err := csgen.WriteGeneratedGoFile(catalogPath, builder.String())
	if err != nil {
		return err
	}

	return nil
}

func writeMigrationFile(migration shared.Migration) string {
	contents := csgen.ExecuteTemplate("migration", migrationTemplateString, migration)

	return contents
}

// LoadManifest loads the manifest file and returns a slice of ObjectMap structs.
func LoadManifest(path ...string) shared.Manifest {
	mp := shared.GetManifestPath(path...)
	log.Printf("Loading manifest file: %s\n", path)
	yfile, err := os.ReadFile(mp)
	if err != nil {
		log.Fatal(err)
	}

	var manifest shared.Manifest
	err = yaml.Unmarshal(yfile, &manifest)
	if err != nil {
		log.Fatal(err)
	}

	if len(manifest.GeneratorPackage) == 0 {
		manifest.GeneratorPackage = csgen.InferPackageFromOutputPath(manifest.GeneratorPath)
	}

	return manifest
}

var catalogTemplateString = `
import (
	"github.com/cscoding21/csmig/shared"
)

func FindDiscoveredMigrations() []shared.Migration {
	out := []shared.Migration{}

	//---Add created migrations here{{range .}}    
	out = append(out, {{ .Name }}){{end}}

	return out
}
`

var migrationTemplateString = `

import "github.com/cscoding21/csmig/shared"

var {{ .Name }} = shared.Migration{
	Name:        "{{.Name}}",
	Description: "{{.Description}}",
	Up: func() error {
		//---your code here
		panic("migration up not implemented")
	},
	Down: func() error {
		// your code here
		panic("migration down not implemented")
	},
}
`
