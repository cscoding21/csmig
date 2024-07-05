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
	"github.com/cscoding21/csmig/shared"
	"github.com/spf13/cobra"
)

// discoveredCmd represents the discovered command
var discoveredCmd = &cobra.Command{
	Use:   "discovered",
	Short: "Output a list of discovered migrations",
	Long: `Discovered migrations are files that have been created which contain the logic for
	a migration version, including its "Up" and "Down" functions.  This command outputs all migrations
	that have been created regardless of whether or not they have been run.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Finding discovered migrations...")

		manifest := shared.LoadManifest()
		discovered := migrate.FindDiscoveredMigrationFiles(manifest)

		for _, d := range discovered {
			fmt.Println(d.Name, d.Description)
		}
	},
}

func init() {
	lsCmd.AddCommand(discoveredCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// discoveredCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// discoveredCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
