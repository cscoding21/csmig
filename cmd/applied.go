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
	"github.com/spf13/cobra"
)

// appliedCmd represents the applied command
var appliedCmd = &cobra.Command{
	Use:   "applied",
	Short: "Output a list of applied migrations",
	Long: `Applied migrations are those which have been run against the data source.  The "applied"
	command will output a list of migrations which have run against the target data source.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Finding applied migrations...")

		config := shared.GetTestConfig()
		strategy, err := persistence.GetPersistenceStrategy(config.DatabaseStrategyName)
		if err != nil {
			panic(err)
		}

		applied, err := migrate.FindAppliedMigrations(strategy)
		if err != nil {
			panic(err)
		}

		for _, a := range applied {
			fmt.Println(a.Name, a.Description, a.AppliedOn)
		}
	},
}

func init() {
	lsCmd.AddCommand(appliedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// appliedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// appliedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
