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
	// "github.com/marlinprotocol/ctl2/cmd/relay/eth/actions"
	// "github.com/marlinprotocol/ctl2/cmd/relay/eth/config"

	"os"

	"github.com/marlinprotocol/ctl2/modules/appcommands"
	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/relay_eth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var EthCmd = &cobra.Command{
	Use:   "eth",
	Short: "Eth relay",
	Long:  `Eth relay`,
}

func init() {
	app, err := appcommands.GetNewApp("relay_eth", projectRunners.GetRunnerInstance,
		appcommands.CommandDetails{Use: "create", DescShort: "Create relay for ethereum blockchain", DescLong: "Create relay for ethereum blockchain"},
		appcommands.CommandDetails{Use: "destroy", DescShort: "Destroy relay for ethereum blockchain", DescLong: "Destroy relay for ethereum blockchain"},
		appcommands.CommandDetails{Use: "logs", DescShort: "Tail logs for running relay (eth) instances", DescLong: "Tail logs for running relay (eth) instances"},
		appcommands.CommandDetails{Use: "status", DescShort: "Show current status of currently running relay instances", DescLong: "Show current status of currently running relay instances"},
		appcommands.CommandDetails{Use: "recreate", DescShort: "Recreate end to end relay (eth) instances", DescLong: "Recreate end to end relay (eth) instances"},
		appcommands.CommandDetails{Use: "restart", DescShort: "Restart services for relay (eth) instances", DescLong: "Restart services for relay (eth) instances"},
		appcommands.CommandDetails{Use: "versions", DescShort: "Show available versions for use", DescLong: "Show available versions for use"},

		appcommands.CommandDetails{Use: "show", DescShort: "Show current configuration residing on disk", DescLong: "Show current configuration residing on disk"},
		appcommands.CommandDetails{Use: "diff", DescShort: "Show soft modifications to config staged for apply", DescLong: "Show soft modifications to config staged for apply"},
		appcommands.CommandDetails{Use: "modify", DescShort: "Modify configs on disk", DescLong: "Modify configs on disk"},
		appcommands.CommandDetails{Use: "reset", DescShort: "Reset Configurations on disk", DescLong: "Reset Configurations on disk"},
		appcommands.CommandDetails{Use: "apply", DescShort: "Apply modifications to config", DescLong: "Apply modifications to config"})
	if err != nil {
		log.Error("Error while creating relay_eth application command tree")
		os.Exit(1)
	}

	EthCmd.AddCommand(app.CreateCmd.Cmd)
	EthCmd.AddCommand(app.DestroyCmd.Cmd)
	EthCmd.AddCommand(app.LogsCmd.Cmd)
	EthCmd.AddCommand(app.StatusCmd.Cmd)
	EthCmd.AddCommand(app.RecreateCmd.Cmd)
	EthCmd.AddCommand(app.RestartCmd.Cmd)
	EthCmd.AddCommand(app.VersionsCmd.Cmd)

	configCmd := &cobra.Command{Use: "config", Short: "Configurations of project set on disk", Long: "Configurations of project set on disk"}
	EthCmd.AddCommand(configCmd)
	configCmd.AddCommand(app.ConfigShowCmd.Cmd)
	configCmd.AddCommand(app.ConfigDiffCmd.Cmd)
	configCmd.AddCommand(app.ConfigModifyCmd.Cmd)
	configCmd.AddCommand(app.ConfigResetCmd.Cmd)
	configCmd.AddCommand(app.ConfigApplyCmd.Cmd)

	// Extra flag additions for relay_eth -----------------------------------------------

	app.CreateCmd.ArgStore["discovery-addrs"] = app.CreateCmd.Cmd.Flags().StringP("discovery-addrs", "a", "127.0.0.1:8002", "Discovery address of relay")
	app.CreateCmd.ArgStore["heartbeat-addrs"] = app.CreateCmd.Cmd.Flags().StringP("heartbeat-addrs", "g", "127.0.0.1:8003", "Heartbeat address of relay")
	app.CreateCmd.ArgStore["datadir"] = app.CreateCmd.Cmd.Flags().StringP("datadir", "d", "~/.ethereum/", "Data directory")
	app.CreateCmd.ArgStore["discovery-port"] = app.CreateCmd.Cmd.Flags().StringP("discovery-port", "f", "", "Discovery port")
	app.CreateCmd.ArgStore["pubsub-port"] = app.CreateCmd.Cmd.Flags().StringP("pubsub-port", "p", "", "PubSub port")
	app.CreateCmd.ArgStore["address"] = app.CreateCmd.Cmd.Flags().StringP("address", "b", "", "Address")
	app.CreateCmd.ArgStore["name"] = app.CreateCmd.Cmd.Flags().StringP("name", "n", "", "Name of relay")
	app.CreateCmd.ArgStore["abci-version"] = app.CreateCmd.Cmd.Flags().StringP("abci-version", "c", "", "ABCI version")
	app.CreateCmd.ArgStore["sync-mode"] = app.CreateCmd.Cmd.Flags().StringP("sync-mode", "m", "light", "Sync mode of GETH")

	// ----------------------------------------------------------------------------------
}
