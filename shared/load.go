package shared

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/cscoding21/csgen"
	"gopkg.in/yaml.v3"
)

// LoadManifest loads the manifest file and returns a slice of ObjectMap structs.
func LoadManifest(path ...string) Manifest {
	mp := getManifestPath(path...)

	yfile, err := os.ReadFile(mp)
	if err != nil {
		log.Fatal(err)
	}

	var manifest Manifest
	err = yaml.Unmarshal(yfile, &manifest)
	if err != nil {
		log.Fatal(err)
	}

	if len(manifest.GeneratorPackage) == 0 {
		manifest.GeneratorPackage = csgen.InferPackageFromOutputPath(manifest.GeneratorPath)
	}

	return manifest
}

func getManifestPath(manifestPath ...string) string {
	//---handle if the user explicitly passes in the manifest path
	if len(manifestPath) > 0 {
		//---user passed in a manifestPath.  Use that instead of the default
		userManifestPath := manifestPath[0]

		//---ensure the manifest file has a valid yaml extension
		if !strings.HasSuffix(userManifestPath, ".yaml") && !strings.HasSuffix(userManifestPath, ".yml") {
			panic(fmt.Sprintf("Manifest file must be a YAML file. %s is not a YAML file.", userManifestPath))
		}

		//---if the manifest file is an absolute path, don't prepend the present working directory
		if path.IsAbs(userManifestPath) {
			return userManifestPath
		}
	}

	//---if not, try to discover by searching in common locations
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	mpOpts := []string{".csmig.yaml", "migrations/.csmig.yaml"}

	for _, mo := range mpOpts {
		mp := path.Join(pwd, mo)
		if _, err := os.Stat(mp); err == nil {
			return mp
		}
	}

	panic("Unable to find manifest file.  Please create a manifest file in the current directory or in the migrations directory.")
}
