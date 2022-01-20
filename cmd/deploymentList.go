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
	"log"

	"github.com/Redpill-Linpro/rtfminion/pkg/anypointclient"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// deploymentCmd represents the deployment command
var deploymentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployments",
	Long: `Lists deploymanets in a specific organization and environment. Optionaly other criteria can be added, see flags.

This command requires you to specifiy an organisation and environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := anypointclient.NewAnypointClientWithCredentials(viper.GetString("anypoint.region"), viper.GetString("anypoint.user"), viper.GetString("anypoint.password"))

		org, err := c.ResolveOrganisation(viper.GetString("anypoint.group"))
		if err != nil {
			log.Fatalf("Failed to resolve organisation %+v", err)
		}
		environment, err := c.ResolveEnvironment(org, viper.GetString("anypoint.environment"))
		if err != nil {
			log.Fatalf("Failed to resolve environment %+v", err)
		}
		fabric := anypointclient.Fabric{}
		if viper.GetString("fabric") != "" {
			fabric, err = c.ResolveFabricByName(org, viper.GetString("fabric"))
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
		payload, _, err := c.ListDeployments(org, environment, fabric)
		if err != nil {
			log.Fatalf("Failed to list deployments %s", err.Error())
		}
		printer.Print(payload)
	},
}

func init() {
	deploymentCmd.AddCommand(deploymentListCmd)
}
