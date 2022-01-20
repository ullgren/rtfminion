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
	"log"
	"os"
	"strings"

	"github.com/Redpill-Linpro/rtfminion/pkg/anypointclient"
	"github.com/Redpill-Linpro/rtfminion/pkg/rtfprinter"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var printer rtfprinter.RTFPrinter

var client anypointclient.AnypointClient

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rtfminion",
	Short: "Minion to help automate operations towards Anypoint Runtime Fabric",
	Long: `This minion helps you automate tasks in Anypoint Runtime Fabric by 
leveraring the Anypoint platform API. Primarly the API related to Runtime Manager 
and Exchange. In this way it differs from rtfctl which mainly interact with the 
Kubernetes cluster directly and is meant to be run on one of the Kubernetes 
cluster nodes.

Using RTF Minion you can currently execute tasks that otherwise normaly is 
done through the Anypoint Runtime Manager UI or in somecases using the 
Mule Maven Plugin.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rtfminion.yaml)")
	rootCmd.PersistentFlags().StringP("user", "u", "", "Anypoint username")
	viper.BindPFlag("anypoint.user", rootCmd.PersistentFlags().Lookup("user"))
	rootCmd.PersistentFlags().StringP("password", "p", "", "Anypoint password")
	viper.BindPFlag("anypoint.password", rootCmd.PersistentFlags().Lookup("password"))
	rootCmd.PersistentFlags().StringP("bearer", "b", "", "authentication bearer token used to authenticate with Anypoint")
	viper.BindPFlag("anypoint.bearer", rootCmd.PersistentFlags().Lookup("bearer"))
	rootCmd.PersistentFlags().StringP("region", "r", anypointclient.REGION_US, "Anypoint region. Allowed values US or EU.")
	viper.BindPFlag("anypoint.region", rootCmd.PersistentFlags().Lookup("region"))
	rootCmd.PersistentFlags().StringP("group", "g", "", "organization or buissness group within Anypoint Platform. Provided as a '/' separated path.")
	viper.BindPFlag("anypoint.group", rootCmd.PersistentFlags().Lookup("group"))
	rootCmd.PersistentFlags().StringP("output", "O", "PRETTY", "Output format. Allowed valued PRETTY and CSV")
	viper.BindPFlag("output.format", rootCmd.PersistentFlags().Lookup("output"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".rtfminion" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rtfminion")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	p, err := rtfprinter.New(os.Stdout, viper.GetString("output.format"))
	if err != nil {
		log.Fatalf("Failed to create printer")
	}
	printer = p

	if viper.GetString("anypoint.bearer") == "" {
		client = anypointclient.NewAnypointClientWithCredentials(viper.GetString("anypoint.region"), viper.GetString("anypoint.user"), viper.GetString("anypoint.password"))
	} else {
		client = anypointclient.NewAnypointClientWithToken(viper.GetString("anypoint.region"), viper.GetString("anypoint.bearer"))
	}
}
