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
package dot

import (
	"os"

	"github.com/marlinprotocol/ctl2/modules/appcommands"
	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/gateway_dot"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var DotCmd = &cobra.Command{
	Use:   "dot",
	Short: "Polkadot Gateway",
	Long:  `Polkadot Gateway`,
}

func init() {
	app, err := appcommands.GetNewApp("gateway_dot", projectRunners.GetRunnerInstance,
		appcommands.CommandDetails{Use: "create", DescShort: "Create gateway for polkadot blockchain", DescLong: "Create gateway for polkadot blockchain"},
		appcommands.CommandDetails{Use: "destroy", DescShort: "Destroy gateway for polkadot blockchain", DescLong: "Destroy gateway for polkadot blockchain"},
		appcommands.CommandDetails{Use: "logs", DescShort: "Tail logs for running gateway (polkadot) instances", DescLong: "Tail logs for running gateway (polkadot) instances"},
		appcommands.CommandDetails{Use: "status", DescShort: "Show status of currently running gateway (polkadot) instances", DescLong: "Show status of currently running gateway (polkadot) instances"},
		appcommands.CommandDetails{Use: "recreate", DescShort: "Recreate end to end gateway (polkadot) instances", DescLong: "Recreate end to end gateway (polkadot) instances"},
		appcommands.CommandDetails{Use: "restart", DescShort: "Restart services for gateway (polkadot) instances", DescLong: "Restart services for gateway (polkadot) instances"},
		appcommands.CommandDetails{Use: "versions", DescShort: "Show available versions for use", DescLong: "Show available versions for use"},

		appcommands.CommandDetails{Use: "show", DescShort: "Show current configuration residing on disk", DescLong: "Show current configuration residing on disk"},
		appcommands.CommandDetails{Use: "diff", DescShort: "Show soft modifications to config staged for apply", DescLong: "Show soft modifications to config staged for apply"},
		appcommands.CommandDetails{Use: "modify", DescShort: "Modify configs on disk", DescLong: "Modify configs on disk"},
		appcommands.CommandDetails{Use: "reset", DescShort: "Reset Configurations on disk", DescLong: "Reset Configurations on disk"},
		appcommands.CommandDetails{Use: "apply", DescShort: "Apply modifications to config", DescLong: "Apply modifications to config"})
	if err != nil {
		log.Error("Error while creating gateway_dot application command tree")
		os.Exit(1)
	}

	DotCmd.AddCommand(app.CreateCmd.Cmd)
	DotCmd.AddCommand(app.DestroyCmd.Cmd)
	DotCmd.AddCommand(app.LogsCmd.Cmd)
	DotCmd.AddCommand(app.StatusCmd.Cmd)
	DotCmd.AddCommand(app.RecreateCmd.Cmd)
	DotCmd.AddCommand(app.RestartCmd.Cmd)
	DotCmd.AddCommand(app.VersionsCmd.Cmd)

	configCmd := &cobra.Command{Use: "config", Short: "Configurations of project set on disk", Long: "Configurations of project set on disk"}
	DotCmd.AddCommand(configCmd)
	configCmd.AddCommand(app.ConfigShowCmd.Cmd)
	configCmd.AddCommand(app.ConfigDiffCmd.Cmd)
	configCmd.AddCommand(app.ConfigModifyCmd.Cmd)
	configCmd.AddCommand(app.ConfigResetCmd.Cmd)
	configCmd.AddCommand(app.ConfigApplyCmd.Cmd)

	// Extra flag additions for gateway_dot -----------------------------------------------

	app.CreateCmd.ArgStore["gateway-keystore-path"] = app.CreateCmd.Cmd.Flags().StringP("gateway-keystore-path", "a", "gateway_dot.key", "Gateway's keystore path")
	app.CreateCmd.ArgStore["gateway-listen-port"] = app.CreateCmd.Cmd.Flags().StringP("gateway-listen-port", "g", "20900", "Port on which gateway listens for connections from peer")
	app.CreateCmd.ArgStore["bridge-discovery-addr"] = app.CreateCmd.Cmd.Flags().StringP("bridge-discovery-addr", "d", "0.0.0.0:20702", "Bridge discovery address")
	app.CreateCmd.ArgStore["bridge-pubsub-addr"] = app.CreateCmd.Cmd.Flags().StringP("bridge-pubsub-addr", "p", "0.0.0.0:20700", "Bridge pubsub address")
	app.CreateCmd.ArgStore["bridge-bootstrap-addr"] = app.CreateCmd.Cmd.Flags().StringP("bridge-bootstrap-addr", "b", "127.0.0.1:8002", "Bridge bootstrap address")
	app.CreateCmd.ArgStore["bridge-listen-addr"] = app.CreateCmd.Cmd.Flags().StringP("bridge-listen-address", "l", "127.0.0.1:20901", "Bridge listen address")
	app.CreateCmd.ArgStore["bridge-keystore-path"] = app.CreateCmd.Cmd.Flags().StringP("bridge-keystore-path", "k", "/etc/dot-keystore-path", "Keystore Path")
	app.CreateCmd.ArgStore["bridge-keystore-pass-path"] = app.CreateCmd.Cmd.Flags().StringP("bridge-keystore-pass-path", "v", "/etc/dot-keystore-pass-path", "Keystore pass path")
	app.CreateCmd.ArgStore["bridge-contracts"] = app.CreateCmd.Cmd.Flags().StringP("bridge-contracts", "c", "mainnet", "mainnet/kovan")

	// ----------------------------------------------------------------------------------
}