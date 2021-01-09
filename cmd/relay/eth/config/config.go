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
package config

import (
	"github.com/marlinprotocol/ctl2/cmd/relay/eth/config/actions"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure state information of the project",
	Long:  `Configure state information of the project`,
}

func init() {
	ConfigCmd.AddCommand(actions.ConfigShowCmd)
	ConfigCmd.AddCommand(actions.ConfigModifyCmd)
	ConfigCmd.AddCommand(actions.ConfigDiffCmd)
	ConfigCmd.AddCommand(actions.ConfigApplyCmd)
	ConfigCmd.AddCommand(actions.ConfigResetCmd)
}