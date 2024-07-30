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

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a discovered migration as long as it has not been applied",
	Long:  `The "remove" command deletes a migration file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("remove called")

		name, _ := cmd.Flags().GetString("name")
		config := shared.GetTestConfig()

		err := generate.RemoveMigration(config, name)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//removeCmd.PersistentFlags().String("name", "", "The name of the migration to remove.  The name does not include the file extension.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	removeCmd.Flags().StringP("name", "n", "", "The name of the migration to remove.  The name does not include the file extension.")
}
