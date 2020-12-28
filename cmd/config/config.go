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
package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marlinprotocol/ctl2/types"
	"github.com/marlinprotocol/ctl2/version"
	log "github.com/sirupsen/logrus"
)

var defaultCreate bool
var location string

// VIPER defaults ------------------------

var defaultReleaseUpstreams = []types.Registry{
	types.Registry{
		Link:    "https://github.com/marlinprotocol/releases.git",
		Branch:  "rtw",
		Enabled: true,
	},
	types.Registry{
		Link:    "https://github.com/marlinprotocol/releases.git",
		Branch:  "beta",
		Enabled: true,
	},
	types.Registry{
		Link:    "https://github.com/marlinprotocol/releases.git",
		Branch:  "alpha",
		Enabled: false,
	},
	types.Registry{
		Link:    "https://github.com/marlinprotocol/releases.git",
		Branch:  "dev",
		Enabled: false,
	},
}

// --------------------------------------

// registryCmd represents the registry command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if defaultCreate {
			log.Warning("Any already existing marlinctl config will be overwritten. Do you really want this? (y/n)")

			if askForConfirmation() == false {
				log.Info("Aborting creation of new marlinctl config")
				return
			}

			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if location == "" {
				location = home + "/.marlinctl/marlinctl_config.yaml"
			}

			viper.SetConfigFile(location)

			for i := 0; i < len(defaultReleaseUpstreams); i++ {
				defaultReleaseUpstreams[i].Local = home + "/.marlinctl/registries/" + defaultReleaseUpstreams[i].Branch
			}
			viper.Set("config_version", version.CfgVersion)
			viper.Set("registries", defaultReleaseUpstreams)
			viper.WriteConfig()

			log.Info("Default marlinctl config written to disk successfully")
		}
	},
}

func init() {
	ConfigCmd.Flags().BoolVarP(&defaultCreate, "default-create", "d", false, "Create default marlinctl config file")
	ConfigCmd.Flags().StringVarP(&location, "location", "l", "", "config file (default is $HOME/.marlinctl/marlinctl_config.yaml)")
}

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
