package cmd

/*
Copyright © 2021 Pontus Ullgren

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

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// groupCmd represents the deployment command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Groups and environments.",
	Long: `Manage bussiness groups and environments within a given organisation. 
	
	This command does nothing instead use one of the sub-commands listed below.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deployment called for")
		viper.Debug()
	},
}

func init() {
	rootCmd.AddCommand(groupCmd)
}
