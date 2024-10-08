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
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init sets up csmig for a project",
	Long: `Init sets up csmig for a project.  By default, a directory named "migrations" will be 
	created where migration assets will be generated.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")

		err := generate.Init()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().StringP("dir", "d", "migrations", "The directory where the migrations framework will be installed.")
}
