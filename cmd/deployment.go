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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deploymentCmd represents the deployment command
var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Handle deployments in a specific environment.",
	Long: `Handle deployments in a specific environment.

This command does nothing instead use one of the sub-commands listed below.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(deploymentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	deploymentCmd.PersistentFlags().StringP("environment", "e", "", "Environment id.")
	viper.BindPFlag("anypoint.environment", deploymentCmd.PersistentFlags().Lookup("environment"))
	deploymentCmd.Flags().StringP("fabric", "f", "", "Only execute operation on deployments on this fabric")
	viper.BindPFlag("fabric", deploymentCmd.Flags().Lookup("fabric"))

}
