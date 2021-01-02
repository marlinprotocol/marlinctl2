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
package actions

import (
	"os"

	"github.com/marlinprotocol/ctl2/modules/registry"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppCmd represents the registry command
var ListVersionsCmd = &cobra.Command{
	Use:     "listversions",
	Short:   "List versions for beacon",
	Long:    `List versions for beacon`,
	PreRunE: ConfigTest,
	Run: func(cmd *cobra.Command, args []string) {
		var projectConfig types.Project
		err := viper.UnmarshalKey(projectId, &projectConfig)
		if err != nil {
			log.Error("Error while reading project config: ", err)
			os.Exit(1)
		}
		versions, err := registry.GlobalRegistry.GetVersions(projectId, projectConfig.Subscription, "0.0.0", "major", projectConfig.Runtime)

		if err != nil {
			log.Error("Error encountered while listing versions: ", err)
			return
		}

		registry.GlobalRegistry.PrettyPrintProjectVersions(versions)
	},
}

func init() {
}
