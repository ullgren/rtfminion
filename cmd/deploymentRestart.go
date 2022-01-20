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
	"regexp"
	"time"

	"github.com/Redpill-Linpro/rtfminion/pkg/anypointclient"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// deploymentCmd represents the deployment command
var deploymentRestartCmd = &cobra.Command{
	Use:   "restart [deployment name]",
	Short: "Restart one or more deployments that matches the given pattern",
	Long: `Restart one or more deployments, that mathcer the given pattern, in a specific organization and environment. 
	
	The deployment pattern is given using go regexp.
	
	Optionaly other criteria can be added, see flags.
	
	This command requires you to specifiy an organisation and environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		pattern, err := regexp.Compile(args[0])
		if err != nil {
			log.Fatalf("Could not compile deployment pattern. %s", err.Error())

		}
		c := anypointclient.NewAnypointClientWithCredentials(viper.GetString("anypoint.region"), viper.GetString("anypoint.user"), viper.GetString("anypoint.password"))

		org, err := c.ResolveOrganisation(viper.GetString("anypoint.group"))
		if err != nil {
			log.Fatalf("Failed to resolve organisation %+v", err)
		}
		// TODO: Do error checks
		environment, _ := c.ResolveEnvironment(org, viper.GetString("anypoint.environment"))
		fabric := anypointclient.Fabric{}
		if viper.GetString("fabric") != "" {
			fabric, err = c.ResolveFabricByName(org, viper.GetString("fabric"))
			if err != nil {
				log.Fatalf(err.Error())
			}
		}
		payload, _, err := c.ListDeployments(org, environment, fabric)
		if err != nil {
			log.Fatalf("Failed to list deployments %+v", err)
		}

		for _, d := range payload.Items {
			if pattern.Match([]byte(d.Name)) {
				if d.Application.Status == "NOT_RUNNING" {
					log.Printf("Skip deployment %s since it is NOT RUNNING",
						d.Name)
					continue
				}
				log.Printf("Stoping %s\n", d.Name)
				c.StopDeployment(org, environment, d)
				if err != nil {
					log.Printf("Failed to stop %s will skip this deployment: %+v", d.Name, err)
					continue
				}
				// Wait for the deployment to stop
				for {
					details, err := client.GetDeploymentDetails(org, environment, d.ID)
					if err != nil {
						// Do nothing
						continue
					}
					if details.Application.Status == "NOT_RUNNING" {
						break
					}
					time.Sleep(time.Duration(10) * time.Millisecond)

				}

				c.StartDeployment(org, environment, d)
				if err != nil {
					log.Fatalf("Failed to start %s : %+v", d.Name, err)
				}
				log.Printf("Starting %s\n", d.Name)
				if viper.GetBool("deploymentrestartwait") {
					for {
						details, err := client.GetDeploymentDetails(org, environment, d.ID)
						if err != nil {
							// Do nothing
							continue
						}
						if details.Application.Status == "RUNNING" {
							break
						}
						time.Sleep(time.Duration(10) * time.Millisecond)

					}
					log.Printf("Started %s\n", d.Name)
				}

			}
		}
	},
}

func init() {
	deploymentCmd.AddCommand(deploymentRestartCmd)

	deploymentRestartCmd.Flags().BoolP("wait", "w", false, "Wait for the deployments to fully start before continuing")
	viper.BindPFlag("deploymentrestartwait", deploymentRestartCmd.Flags().Lookup("wait"))

}
