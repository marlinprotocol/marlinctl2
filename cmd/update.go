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
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/marlinprotocol/ctl2/modules/upstream"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update marlinctl registry with upstream",
	Long:  `marlinctl updates binaries and related configurations from remote using this command`,
	Run: func(cmd *cobra.Command, args []string) {
		upstreamCfg := upstream.UpstreamConfig{
			RemoteVCS: viper.GetString("upstream"),
			HomeClone: viper.GetString("marlindir") + "/releases",
		}

		returnMessage, err := upstreamCfg.FetchUpstreamRegistry()

		if err != nil {
			log.Error("Error while fetching registry from upstream: ", err.Error())
		} else {
			log.Info("Update completed: " + returnMessage)
		}

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
