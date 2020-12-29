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

	"github.com/marlinprotocol/ctl2/modules/util"
	"github.com/marlinprotocol/ctl2/types"
)

var enableBeta, clearCache, forceRuntime bool
var version, runtime string

var projectId string = "iris_endnode"

// AppCmd represents the registry command
var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var projectConfig types.Project

		if viper.IsSet(projectId) {
			err := viper.UnmarshalKey(projectId, &projectConfig)
			if err != nil {
				log.Error("Error encountered while retrieving project configuration: ", err)
				return
			}

			if clearCache {
				err := clearCacheFunc(projectConfig, projectId)
				if err != nil {
					log.Error("Error encountered while removing project cache: ", err)
					return
				}
			}
		} else {

		}

		if !clearCache {
			suitableRuntimes := util.GetRuntimes()

			if !forceRuntime {
				if suitable, ok := suitableRuntimes[runtime]; !ok || !suitable {
					log.Error("Runtime provided for configuration: " + runtime +
						" may not be supported by marlinctl or is not supported by your system." +
						" If you think this is incorrect, override this check using --force-runtime.")
					os.Exit(1)
				} else {
					log.Info("Runtime provided for configuration: " + runtime +
						" seems to be supported. Going ahead with configuring this.")
				}
			}

			var releaseSubscriptions = []string{"rtw"}

			if enableBeta {
				releaseSubscriptions = append(releaseSubscriptions, "beta")
			}

			if version != "latest" {
				// TODO
				log.Warn("Yet to implement non latest versions")
			}

			viper.Set(projectId, types.Project{
				Subscription:  releaseSubscriptions,
				Version:       "latest",
				Storage:       viper.GetString("homedir") + "/project/" + projectId,
				Runtime:       runtime,
				ForcedRuntime: false,
			})
			viper.WriteConfig()
		}
	},
}

func init() {
	ConfigureCmd.Flags().BoolVarP(&clearCache, "clear-cache", "c", false, "clear any locally available binaries and configurations")
	ConfigureCmd.Flags().BoolVarP(&enableBeta, "enable-beta", "b", false, "enable beta releases")
	ConfigureCmd.Flags().StringVarP(&version, "version", "v", "latest", "Version to run")
	ConfigureCmd.Flags().StringVarP(&runtime, "runtime", "r", "", "Application runtime")
	ConfigureCmd.Flags().BoolVarP(&forceRuntime, "force-runtime", "f", false, "Forcefully set application runtime")
}

func clearCacheFunc(projectConfig types.Project, projectId string) error {
	err := util.RemoveDirPathIfExists(projectConfig.Storage)
	if err != nil {
		return err
	}

	err = util.RemoveConfigEntry(projectId)
	if err != nil {
		return err
	}
	return nil
}
