package generate

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/cscoding21/csgen"
	"github.com/cscoding21/csmig/migrate"
	"github.com/cscoding21/csmig/persistence"
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
	err = writeRunner(manifest, migrationsDir)
	if err != nil {
		return err
	}

	err = writeRunnerTest(manifest, migrationsDir)
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
func NewMigration(manifest shared.Manifest, description string) (shared.Migration, error) {
	builder := csgen.NewCSGenBuilderForOneOffFile("csmig", manifest.GeneratorPackage)

	migrationName := getMigrationName()
	migration := shared.Migration{
		Package:     manifest.GeneratorPackage,
		Name:        migrationName,
		Description: description,
	}

	builder.WriteString(getMigrationFileContents(migration))
	migrationFileName := fmt.Sprintf("%s_gen.go", migrationName)
	migrationFilePath := path.Join(manifest.ProjectRoot, manifest.GeneratorPath, migrationFileName)

	err := csgen.WriteGeneratedGoFile(migrationFilePath, builder.String())
	if err != nil {
		log.Fatal(err)
	}

	writeCatalogFile(manifest)

	return migration, nil
}

func RemoveMigration(manifest shared.Manifest, name string) error {
	//---remove the migration file
	strategy, err := persistence.GetPersistenceStrategy(manifest.VersionStrategy)
	if err != nil {
		return err
	}

	am, err := strategy.FindAppliedMigrations(strategy.DBConfig)
	if err != nil {
		return err
	}

	//---stop if the migration has already been applied as this will taint the version sequence
	for _, m := range am {
		if m.Name == name {
			return fmt.Errorf("cannot remove migration %s, it has already been applied", name)
		}
	}

	//---remove the migration file
	migrationFilePath := path.Join(manifest.ProjectRoot, manifest.GeneratorPath, fmt.Sprintf("%s_gen.go", name))
	err = os.Remove(migrationFilePath)
	if err != nil {
		return err
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
	file := path.Join(outputPath, ".csmig.yaml")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		builder := strings.Builder{}
		builder.WriteString(csgen.ExecuteTemplate("csmig_manifest", manifestTemplateString, manifest))

		err = os.MkdirAll(outputPath, 0755)
		if err != nil {
			return err
		}

		err = os.WriteFile(file, []byte(builder.String()), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func writeRunner(manifest shared.Manifest, outputPath string) error {
	builder := csgen.NewCSGenBuilderForFile("csmig", manifest.GeneratorPackage)
	builder.WriteString(runFileTemplateString)

	file := path.Join(outputPath, "runner.gen.go")
	return csgen.WriteGeneratedGoFile(file, builder.String())
}

func writeRunnerTest(manifest shared.Manifest, outputPath string) error {
	builder := csgen.NewCSGenBuilderForFile("csmig", manifest.GeneratorPackage)
	builder.WriteString(runFileTestTemplateString)

	file := path.Join(outputPath, "runner_test.go")
	return csgen.WriteGeneratedGoFile(file, builder.String())
}

func getMigrationFileContents(migration shared.Migration) string {
	contents := csgen.ExecuteTemplate("migration", migrationTemplateString, migration)

	return contents
}

var catalogTemplateString = `
import (
	"github.com/cscoding21/csmig/shared"
)

func FindDiscoveredMigrations() []shared.Migration {
	out := []shared.Migration{}

	//---Generated migrations will be appended here via code generation{{range .}}    
	out = append(out, {{ .Name }}){{end}}

	return out
}
`

var migrationTemplateString = `
import (
	"fmt"
	"github.com/cscoding21/csmig/shared"
)

var {{ .Name }} = shared.Migration{
	Name:        "{{.Name}}",
	Description: "{{.Description}}",
	Up: func(ds shared.DatabaseStrategy) error {
		//---your code here
		fmt.Printf("migration up for {{ .Name }} not implemented")

		return nil
	},
	Down: func(ds shared.DatabaseStrategy) error {
		// your code here
		fmt.Printf("migration down for {{ .Name }} not implemented")

		return nil
	},
}
`

var runFileTemplateString = `
import (
	"github.com/cscoding21/csmig/migrate"
	"github.com/cscoding21/csmig/persistence"
	"github.com/cscoding21/csmig/shared"
)

// Apply run any migrations that have not been applied yet.
func Apply(manifest shared.Manifest) error {
	discoveredMigrations := FindDiscoveredMigrations()

	//---get the persistence strategy as defined in the manifest
	strategy, err := persistence.GetPersistenceStrategy(manifest.VersionStrategy)
	if err != nil {
		return err
	}

	//---make sure the required support tables have been created
	err = strategy.EnsureInfrastructure(strategy.DBConfig)
	if err != nil {
		return err
	}

	//---retrieve a list of migrations that have already been applied
	appliedMigrations, err := strategy.FindAppliedMigrations(strategy.DBConfig)
	if err != nil {
		return err
	}

	//---iterate over the migrations that have been created and apply any that have not been applied yet
	for _, dm := range discoveredMigrations {
		//---skip if migration is already applied
		if migrationIsApplied(dm.Name, appliedMigrations) {
			continue
		}

		err = dm.Up(strategy)
		if err != nil {
			return err
		}

		err = strategy.ApplyMigration(strategy.DBConfig, dm.Name, dm.Description)
		if err != nil {
			return err
		}
	}

	return nil
}

// Rollback call the "Down" method of the most recently applied migration
func Rollback(manifest shared.Manifest) error {
	strategy, err := persistence.GetPersistenceStrategy(manifest.VersionStrategy)
	if err != nil {
		return err
	}

	appliedMigrations, err := strategy.FindAppliedMigrations(strategy.DBConfig)
	if err != nil {
		return err
	}

	latestMigration := getLatestMigration(appliedMigrations)

	if latestMigration == nil {
		return nil
	}

	for _, dm := range FindDiscoveredMigrations() {
		if latestMigration.Name == dm.Name {
			err = dm.Down(strategy)
			if err != nil {
				return err
			}
		}
	}

	return strategy.RollbackMigration(strategy.DBConfig, latestMigration.Name)
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

func getLatestMigration(appliedMigrations []shared.AppliedMigration) *shared.AppliedMigration {
	if len(appliedMigrations) == 0 {
		return nil
	}

	out := appliedMigrations[len(appliedMigrations)-1]

	for _, am := range appliedMigrations {
		if am.Name > out.Name {
			out = am
		}
	}

	return &out
}
`

var runFileTestTemplateString = `
import (
	"fmt"
	"testing"

	"github.com/cscoding21/csmig/shared"
)

func TestApply(t *testing.T) {
	manifest := shared.LoadManifest(".csmig.yaml")

	err := Apply(manifest)
	if err != nil {
		t.Error(err)
	}
}

func TestRollback(t *testing.T) {
	manifest := shared.LoadManifest(".csmig.yaml")

	err := Rollback(manifest)
	if err != nil {
		t.Error(err)
	}
}

func TestFindApplyMigrations(t *testing.T) {
	manifest := shared.LoadManifest(".csmig.yaml")

	appliedMigrations, err := FindAppliedMigrations(manifest)
	if err != nil {
		t.Error(err)
	}

	for _, am := range appliedMigrations {
		fmt.Println(am.Name, am.AppliedOn)
	}
}

func TestFindUnapplyMigrations(t *testing.T) {
	manifest := shared.LoadManifest(".csmig.yaml")

	unappliedMigrations, err := FindUnappliedMigrations(manifest)
	if err != nil {
		t.Error(err)
	}

	if len(unappliedMigrations) == 0 {
		fmt.Println("No unapplied migrations were found")
	}

	for _, um := range unappliedMigrations {
		fmt.Println(um.Name, um.FilePath)
	}
}
 `

var manifestTemplateString = `
#####################################################################
# Common Sense Coding: CSMig Manifest File
# https://github.com/cscoding21/csmig

# The project root.  This will be used when necessary for determining file locations
project_root: {{.ProjectRoot}}

# The path where migrations functionality will live
generator_path: {{.GeneratorPath}}

# The name of the implementation that is used for naming and comments for generated files
# implementation_name: {{.ImplementationName}}

# The package name that will be used for generated migration files
generator_package: {{.GeneratorPackage}}

# The database (or other persistence strategy) that will be used to track migration state
version_strategy: {{.VersionStrategy}}
`
