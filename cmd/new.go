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

	"github.com/cscoding21/csmig/generate"
	"github.com/cscoding21/csmig/shared"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new migration in the configured directory",
	Long: `The "new" command generated the scaffold for a new migration version and writes it
	to the configured directory.  It accepts an optional description to help developers understand
	what the migration is intended to do.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating new migration...")

		message, _ := cmd.Flags().GetString("message")
		config := shared.GetTestConfig()

		mig, err := generate.NewMigration(config, message)
		if err != nil {
			panic(err)
		}

		fmt.Println("Migration created: ", mig.Name)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	newCmd.Flags().StringP("message", "m", "", "A description of the migration's general purpose.")
}
