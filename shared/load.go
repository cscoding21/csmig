package shared

import (
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
	log.Printf("Loading manifest file: %s\n", path)
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
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	mp := ".csmap.yaml"
	if len(manifestPath) > 0 {
		//---user passed in a manifestPath.  Use that instead of the default
		mp = manifestPath[0]

		//---ensure the manifest file has a valid yaml extension
		if !strings.HasSuffix(mp, ".yaml") && !strings.HasSuffix(mp, ".yml") {
			log.Fatalf("Manifest file must be a YAML file. %s is not a YAML file.", mp)
		}

		//---if the manifest file is a relative path, don't prepend the present working directory
		if path.IsAbs(mp) {
			pwd = ""
		}
	}

	return path.Join(pwd, mp)
}
