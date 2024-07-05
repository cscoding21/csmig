/*
Copyright © 2024 Jeff Kody <jeph@cscoding.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/cscoding21/csmig/migrate"
	"github.com/cscoding21/csmig/persistence"
	"github.com/cscoding21/csmig/shared"
	"github.com/cscoding21/csmig/version"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Outputs a status of the migration configuration for this project",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Migration status...")

		manifest := shared.LoadManifest()
		strategy, _ := persistence.GetPersistenceStrategy(manifest.VersionStrategy)

		discoverd := migrate.FindDiscoveredMigrationFiles(manifest)
		applied, _ := migrate.FindAppliedMigrations(strategy)

		fmt.Println("--------------- CSMig Status ---------------")
		fmt.Println("CSMig Version: ", version.Version)
		fmt.Println("Migrations Directory: ", manifest.GetMigrationPath())
		fmt.Println("Persistence Strategy: ", manifest.VersionStrategy)
		fmt.Println("---")
		fmt.Println("Discovered Migrations: ")
		for _, d := range discoverd {
			fmt.Println("  - ", d.Name)
		}

		fmt.Println("---")
		fmt.Println("Applied Migrations: ")
		for _, a := range applied {
			fmt.Printf("  - %s (%s) : %s \n", a.Name, a.AppliedOn, a.Description)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
