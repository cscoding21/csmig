package generate

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/cscoding21/csgen"
	"github.com/cscoding21/csmig/migrate"
	"github.com/cscoding21/csmig/shared"
)

func Init() error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//---create an initial manifest it doesn't exist
	manifest := shared.Manifest{
		ProjectRoot:      pwd,
		GeneratorPath:    "migrations",
		GeneratorPackage: "migrations",
		VersionStrategy:  "surrealdb",
	}
	migrationsDir := path.Join(manifest.ProjectRoot, manifest.GeneratorPath)

	//---create the initial manifest if it doesn't exist
	err = writeInitialManifefst(manifest, migrationsDir)
	if err != nil {
		return err
	}

	//---create or overwrite the runner file
	err = writeInitialRunner(manifest, migrationsDir)
	if err != nil {
		return err
	}

	//---create an initial catalog file
	err = writeCatalogFile(manifest)
	if err != nil {
		return err
	}

	return nil
}

// NewMigration creates a new migration
func NewMigration(manifestPath string, description string) error {
	manifest := shared.LoadManifest(manifestPath)

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
	migrations := migrate.FindDiscoveredMigrationFiles(manifest)
	contents := csgen.ExecuteTemplate[[]shared.Migration]("catalog", catalogTemplateString, migrations)

	builder := csgen.NewCSGenBuilderForFile("csmig", manifest.GeneratorPackage)
	builder.WriteString(contents)

	catalogPath := path.Join(manifest.ProjectRoot, manifest.GeneratorPath, "catalog.gen.go")
	err := csgen.WriteGeneratedGoFile(catalogPath, builder.String())
	if err != nil {
		return err
	}

	return nil
}

func writeInitialManifefst(manifest shared.Manifest, outputPath string) error {
	file := path.Join(outputPath, "manifest.yaml")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		builder := strings.Builder{}
		builder.WriteString(csgen.ExecuteTemplate("csmig_manifest", manifestTemplateString, manifest))

		err = os.WriteFile(file, []byte(builder.String()), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeInitialRunner(manifest shared.Manifest, outputPath string) error {
	builder := csgen.NewCSGenBuilderForFile("csmig", manifest.GeneratorPackage)
	builder.WriteString(runFileTemplateString)

	file := path.Join(outputPath, "runner.gen.go")
	return csgen.WriteGeneratedGoFile(file, builder.String())
}

func writeMigrationFile(migration shared.Migration) string {
	contents := csgen.ExecuteTemplate("migration", migrationTemplateString, migration)

	return contents
}

// // LoadManifest loads the manifest file and returns a slice of ObjectMap structs.
// func LoadManifest(path ...string) shared.Manifest {
// 	mp := shared.GetManifestPath(path...)
// 	log.Printf("Loading manifest file: %s\n", path)
// 	yfile, err := os.ReadFile(mp)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var manifest shared.Manifest
// 	err = yaml.Unmarshal(yfile, &manifest)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if len(manifest.GeneratorPackage) == 0 {
// 		manifest.GeneratorPackage = csgen.InferPackageFromOutputPath(manifest.GeneratorPath)
// 	}

// 	return manifest
// }

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

var runFileTemplateString = `
import (
	"github.com/cscoding21/csmig/migrate"
	"github.com/cscoding21/csmig/persistence"
	"github.com/cscoding21/csmig/shared"
)

// ApplyMigrations run any migrations that have not been applied yet.
func ApplyMigrations(manifest shared.Manifest) error {
	discoveredMigrations := FindDiscoveredMigrations()

	strategy, err := persistence.GetPersistenceStrategy(manifest.VersionStrategy)
	if err != nil {
		return err
	}

	appliedMigrations, err := strategy.FindAppliedMigrations(strategy.DBConfig)
	if err != nil {
		return err
	}

	for _, dm := range discoveredMigrations {
		//---skip if migration is already applied
		if migrationIsApplied(dm.Name, appliedMigrations) {
			continue
		}

		err = dm.Up()
		if err != nil {
			return err
		}

		err = strategy.ApplyMigration(strategy.DBConfig, dm.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// FindAppliedMigrations return a list of all migrations that have been applied
func FindAppliedMigrations(manifest shared.Manifest) ([]shared.AppliedMigration, error) {
	strategy, err := persistence.GetPersistenceStrategy(manifest.VersionStrategy)
	if err != nil {
		return nil, err
	}

	return strategy.FindAppliedMigrations(strategy.DBConfig)
}

// FindUnappliedMigrations return a list of migrations that have not been applied yet.
func FindUnappliedMigrations(manifest shared.Manifest) ([]shared.Migration, error) {
	strategy, err := persistence.GetPersistenceStrategy(manifest.VersionStrategy)
	if err != nil {
		return nil, err
	}

	discoveredMigrations := FindDiscoveredMigrations()
	appliedMigrations, err := migrate.FindAppliedMigrations(strategy)
	if err != nil {
		return nil, err
	}

	out := []shared.Migration{}

	for _, dm := range discoveredMigrations {
		if !migrationIsApplied(dm.Name, appliedMigrations) {
			out = append(out, dm)
		}
	}

	return out, nil
}

func migrationIsApplied(name string, appliedMigrations []shared.AppliedMigration) bool {
	for _, appliedMigration := range appliedMigrations {
		if appliedMigration.Name == name {
			return true
		}
	}

	return false
}

`

var manifestTemplateString = `
#####################################################################
# Common Sense Coding: CSMig Manifest File
# https://github.com/cscoding21/csmig

# The project root.  This will be used when necessary for determining file locations
project_root: /home/jeph/projects/cscoding21/csmig

# The path where migrations functionality will live
generator_path: migrations

# The name of the implementation that is used for naming and comments for generated files
# implementation_name: csmig

# The package name that will be used for generated migration files
generator_package: migrations

# The database (or other persistence strategy) that will be used to track migration state
version_strategy: surrealdb
`
