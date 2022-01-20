package cmd

/*
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
	"log"

	"github.com/Redpill-Linpro/rtfminion/pkg/anypointclient"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// deploymentCmd represents the deployment command
var fabricListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Runtime Fabric instances in the organisation",
	Long:  `List all Runtime Fabric instances available in the specified organisation.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := anypointclient.NewAnypointClientWithCredentials(viper.GetString("anypoint.region"), viper.GetString("anypoint.user"), viper.GetString("anypoint.password"))

		org, err := c.ResolveOrganisation(viper.GetString("anypoint.group"))
		if err != nil {
			log.Fatalf("Failed to resolve organistation %s. %s", viper.GetString("anypoint.group"), err.Error())
		}
		payload, err := c.ListFabrics(org)
		if err != nil {
			log.Fatalf("Failed to get list of fabric %s", err.Error())
		}
		printer.Print(payload)
	},
}

func init() {
	fabricCmd.AddCommand(fabricListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deploymentListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deploymentListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
