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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marlinprotocol/ctl2/modules/registry"
	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/iris_endnode"
	"github.com/marlinprotocol/ctl2/types"
)

// AppCmd represents the registry command
var LogsCmd = &cobra.Command{
	Use:   "logs",
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

		var versionToRun registry.ProjectVersion
		if projectConfig.Version == "latest" {
			if len(versions) > 0 {
				versionToRun = versions[0]
				log.Info("Latest version is being picked: ", versionToRun.Version)
			} else {
				log.Error("No version available to run for latest for this project. Aborting")
				os.Exit(1)
			}
		} else {
			var isVersionAvailable bool = false
			for _, v := range versions {
				if projectConfig.Version == v.Version {
					isVersionAvailable = true
					versionToRun = v
					break
				}
			}
			if !isVersionAvailable {
				log.Error("Explicitly configured version " + projectConfig.Version + " is not available in registries. Aborting")
				os.Exit(1)
			}
		}

		runner, err := projectRunners.GetRunnerInstance(versionToRun.RunnerId, versionToRun.Version, projectConfig.Storage, versionToRun.RunnerData, skipChecksum)
		if err != nil {
			log.Error("Cannot get runner: ", err.Error())
			os.Exit(1)
		}

		err = runner.PreRunSanity()
		if err != nil {
			log.Error("Failure during pre run sanity: ", err.Error())
			return
		}

		err = runner.Logs()
		if err != nil {
			log.Error("Failure during logging: ", err.Error())
			return
		}
	},
}

func init() {
	// NIL
}
