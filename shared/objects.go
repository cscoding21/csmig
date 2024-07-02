package shared

import (
	"path"
	"time"
)

type Migration struct {
	FilePath    string `yaml:"file_path"`
	Package     string `yaml:"package"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Up          func() error
	Down        func() error
}

type AppliedMigration struct {
	Name      string    `yaml:"name"`
	AppliedOn time.Time `yaml:"applied_on"`
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

// GetMigrationPath return the path that migrations will be stored based on properties in the manifest object.
func (manifest *Manifest) GetMigrationPath() string {
	return path.Join(manifest.ProjectRoot, manifest.GeneratorPath)
}
