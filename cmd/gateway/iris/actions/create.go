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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marlinprotocol/ctl2/modules/registry"
	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/gateway_iris"
	"github.com/marlinprotocol/ctl2/types"
)

var runtimeArgs map[string]string
var instanceId string
var skipChecksum bool

// AppCmd represents the registry command
var CreateCmd = &cobra.Command{
	Use:     "create",
	Short:   `Create a gateway on local system`,
	PreRunE: ConfigTest,
	Run: func(cmd *cobra.Command, args []string) {
		var projectConfig types.Project
		err := viper.UnmarshalKey(projectId, &projectConfig)
		if err != nil {
			log.Error("Error while reading project config: ", err)
			return
		}

		versionToRun, err := registry.GlobalRegistry.GetVersionToRun(projectId)
		if err != nil {
			log.Error("Error while getting version to run: ", err)
			return
		}

		runner, err := projectRunners.GetRunnerInstance(versionToRun.RunnerId, versionToRun.Version, projectConfig.Storage, versionToRun.RunnerData, false, skipChecksum, instanceId)
		if err != nil {
			log.Error("Cannot get runner: ", err.Error())
			return
		}

		err = runner.PreRunSanity()
		if err != nil {
			log.Error("Failure during pre run sanity: ", err.Error())
			return
		}

		err = runner.Prepare()
		if err != nil {
			log.Error("Failure during preparation: ", err.Error())
			return
		}

		err = runner.Create(runtimeArgs)
		if err != nil {
			log.Error("Failure during start: ", err.Error())
			return
		}

		projectConfig.CurrentVersion = versionToRun.Version

		viper.Set(projectId, projectConfig)
		err = viper.WriteConfig()
		if err != nil {
			log.Error("Failure while updating config for current version: ", err.Error())
			return
		}
	},
}

func init() {
	runtimeArgs = make(map[string]string)
	CreateCmd.Flags().StringVarP(&instanceId, "instance-id", "i", "001", "instance-id of the resource")
	CreateCmd.Flags().BoolVarP(&skipChecksum, "skip-checksum", "s", false, "skips checking file integrity during run")
	CreateCmd.Flags().StringToStringVarP(&runtimeArgs, "runtime-arguments", "r", map[string]string{}, "runtime arguments for iris endnode")
}
