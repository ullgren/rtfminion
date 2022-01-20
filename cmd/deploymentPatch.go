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
	"strings"
	"time"

	"github.com/Redpill-Linpro/rtfminion/pkg/anypointclient"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// deploymentCmd represents the deployment command
var deploymentPatchCmd = &cobra.Command{
	Use:   "patch [deployment name]",
	Short: "Patches the runtime for one or more deployments that matches the given pattern",
	Long: `Patches the runtime version for one or more deployment, that mathcer the given pattern, in a specific organization and environment. 
By default the version is updated to the latest patch version of the running version.

The deployment pattern is given using go regexp.

Optionaly other criteria can be added, see flags.

This command requires you to specifiy an organisation and environment.`,
	Args: cobra.ExactArgs(1),
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

		// TODO: Validate that the toversion is a valid runtime version
		for _, d := range payload.Items {
			if pattern.Match([]byte(d.Name)) {
				/*if d.Application.Status == "NOT_RUNNING" {
					log.Printf("Skip deployment %s since it is NOT RUNNING",
						d.Name)
					continue
				}*/
				if viper.GetString("fromversion") != "" {
					matched, err := regexp.MatchString(viper.GetString("fromversion"), d.Target.DeploymentSettings.RuntimeVersion)
					if !matched || err != nil {
						log.Printf("Skip deployment %s since version %s does not match %s",
							d.Name,
							d.Target.DeploymentSettings.RuntimeVersion,
							viper.GetString("fromversion"))
						continue
					}
				}
				if strings.HasPrefix(d.Target.DeploymentSettings.RuntimeVersion, viper.GetString("toversion")) {
					log.Printf("Skip deployment %s since version is already %s",
						d.Name,
						d.Target.DeploymentSettings.RuntimeVersion)
					continue
				}
				d.Target.DeploymentSettings.RuntimeVersion = viper.GetString("toversion")

				err = c.UpdateDeployment(org, environment, d)
				if err != nil {
					log.Fatalf("Failed to update %s : %+v", d.Name, err)
				}
				log.Printf("Updated %s to %s\n",
					d.Name, d.Target.DeploymentSettings.RuntimeVersion)
				if viper.GetInt("delay") != 0 {
					time.Sleep(time.Duration(viper.GetInt("delay")) * time.Millisecond)
				}
			}
		}
	},
}

func init() {
	deploymentCmd.AddCommand(deploymentPatchCmd)

	deploymentPatchCmd.Flags().String("fromversion", "", "Version to update from. If specifified only update if current version matches this pattern")
	viper.BindPFlag("fromversion", deploymentPatchCmd.Flags().Lookup("fromversion"))
	deploymentPatchCmd.Flags().String("toversion", "", "Version to update to")
	viper.BindPFlag("toversion", deploymentPatchCmd.Flags().Lookup("toversion"))
	deploymentPatchCmd.Flags().Uint("delay", 0, "Delay in milliseconds between updates. Default is 0 (no delay)")
	viper.BindPFlag("delay", deploymentPatchCmd.Flags().Lookup("delay"))

}
