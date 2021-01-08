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
package eth

import (
	"github.com/marlinprotocol/ctl2/cmd/relay/eth/actions"
	"github.com/marlinprotocol/ctl2/cmd/relay/eth/config"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var EthCmd = &cobra.Command{
	Use:   "eth",
	Short: "Eth relay",
	Long:  `Eth relay`,
}

func init() {
	EthCmd.AddCommand(actions.CreateCmd)
	EthCmd.AddCommand(actions.RestartCmd)
	EthCmd.AddCommand(actions.RecreateCmd)
	EthCmd.AddCommand(actions.StatusCmd)
	EthCmd.AddCommand(actions.DestroyCmd)
	EthCmd.AddCommand(actions.VersionsCmd)
	EthCmd.AddCommand(actions.LogsCmd)
	EthCmd.AddCommand(config.ConfigCmd)
}
