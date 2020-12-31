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
	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/gateway_iris"
	"github.com/marlinprotocol/ctl2/types"
)

// AppCmd represents the registry command
var DestroyCmd = &cobra.Command{
	Use:     "destroy",
	Short:   "Destroy any running iris gateway",
	Long:    `Destroy any running iris gateway`,
	PreRunE: ConfigTest,
	Run: func(cmd *cobra.Command, args []string) {
		var projectConfig types.Project
		err := viper.UnmarshalKey(projectId, &projectConfig)
		if err != nil {
			log.Error("Error while reading project config: ", err)
			os.Exit(1)
		}

		versionToRun, err := registry.GlobalRegistry.GetVersionToRun(projectId)
		if err != nil {
			log.Error("Error while getting version to run: ", err)
			os.Exit(1)
		}

		runner, err := projectRunners.GetRunnerInstance(versionToRun.RunnerId, versionToRun.Version, projectConfig.Storage, versionToRun.RunnerData, skipChecksum, instanceId)
		if err != nil {
			log.Error("Cannot get runner: ", err.Error())
			os.Exit(1)
		}

		err = runner.PreRunSanity()
		if err != nil {
			log.Error("Failure during pre run sanity: ", err.Error())
			return
		}

		err = runner.Destroy()
		if err != nil {
			log.Error("Failure during destroy: ", err.Error())
			log.Warning("Destroy failure can occur when creation and destruction of processes is done manually and not all through marlinctl." +
				" Failure may not reflect current process state.")
			return
		}

		err = runner.PostRun()
		if err != nil {
			log.Error("Failure during post run: ", err.Error())
			return
		}
	},
}

func init() {
	DestroyCmd.Flags().StringVarP(&instanceId, "instance-id", "i", "001", "instance-id of the resource")
}
