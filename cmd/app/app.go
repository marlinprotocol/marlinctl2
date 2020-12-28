/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package app

import (
	"fmt"
	"os"
	"strconv"

	"github.com/marlinprotocol/ctl2/cmd/app/projects/iris"
	"github.com/marlinprotocol/ctl2/types"
	"github.com/marlinprotocol/ctl2/version"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marlinprotocol/ctl2/modules/registry"
)

var cfgFile string

// AppCmd represents the registry command
var AppCmd = &cobra.Command{
	Use:   "app",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		readConfig()

		var configuredRegistries []types.Registry

		err := viper.UnmarshalKey("registries", &configuredRegistries)
		if err != nil {
			log.Error("Error reading registries from cfg file: ", err)
			os.Exit(1)
		}

		registry.SetupGlobalRegistry(configuredRegistries)

		registry.GlobalRegistry.Sync()
	},
}

func init() {
	AppCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.marlinctl/marlinctl_config.yaml)")
	AppCmd.AddCommand(iris.IrisCmd)
}

// initConfig reads in config file and ENV variables if set.
func readConfig() {
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

		viper.SetConfigFile(home + "/.marlinctl/marlinctl_config.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		var cfgVersionOnDisk = viper.GetInt("config_version")
		if cfgVersionOnDisk != version.CfgVersion {
			log.Error("Cannot use the given config file as it does not match marlinctl's cfgversion. Wanted "+strconv.Itoa(version.CfgVersion), " but found ", cfgVersionOnDisk)
			os.Exit(1)
		}
		log.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Error("No config file available on local machine. Please create one first.")
		os.Exit(1)
	}
}
