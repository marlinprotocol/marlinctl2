/*
Copyright © 2020 MARLIN TEAM <info@marlin.pro>

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

	cmn "github.com/marlinprotocol/ctl2/cmd/relay/eth/common"
	cfg "github.com/marlinprotocol/ctl2/cmd/relay/eth/config"
	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/relay_eth"
	"github.com/marlinprotocol/ctl2/types"
)

// AppCmd represents the registry command
var RestartCmd = &cobra.Command{
	Use:     "restart",
	Short:   "Trigger restart for the service",
	Long:    `Trigger restart fot the service`,
	PreRunE: cfg.ConfigTest,
	Run: func(cmd *cobra.Command, args []string) {
		var projectConfig types.Project
		err := viper.UnmarshalKey(cmn.ProjectID, &projectConfig)
		if err != nil {
			log.Error("Error while reading project config: ", err)
			os.Exit(1)
		}

		runnerId, version, err := cmn.GetResourceMetaData(projectConfig, instanceId)
		if err != nil {
			log.Error("Error while fetching resource information: ", err)
			os.Exit(1)
		}

		runner, err := projectRunners.GetRunnerInstance(runnerId, version, projectConfig.Storage, struct{}{}, true, true, instanceId)
		if err != nil {
			log.Error("Cannot get runner: ", err.Error())
			os.Exit(1)
		}

		err = runner.PreRunSanity()
		if err != nil {
			log.Error("Failure during pre run sanity: ", err.Error())
			return
		}

		err = runner.Restart()
		if err != nil {
			log.Error("Failure during restart: ", err.Error())
			return
		}
	},
}

func init() {
	RestartCmd.Flags().StringVarP(&instanceId, "instance-id", "i", "001", "instance-id of the resource")
}
