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
	"github.com/marlinprotocol/ctl2/modules/registry"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppCmd represents the registry command
var ListVersionsCmd = &cobra.Command{
	Use:   "listversions",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var projectConfig types.Project
		err := viper.UnmarshalKey(projectId, &projectConfig)
		versions, err := registry.GlobalRegistry.GetVersions("iris_endnode", projectConfig.Subscription, projectConfig.Runtime)

		if err != nil {
			log.Error("Error encountered while listing versions: ", err)
			return
		}

		registry.GlobalRegistry.PrettyPrintProjectVersions(versions)
	},
}

func init() {
	ListVersionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
