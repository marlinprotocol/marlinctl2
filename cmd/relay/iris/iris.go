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
package iris

import (
	// "github.com/marlinprotocol/ctl2/cmd/relay/eth/actions"
	// "github.com/marlinprotocol/ctl2/cmd/relay/eth/config"

	"os"

	"github.com/marlinprotocol/ctl2/modules/appcommands"
	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/relay_iris"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var IrisCmd = &cobra.Command{
	Use:   "iris",
	Short: "iris relay",
	Long:  `iris relay`,
}

func init() {
	app, err := appcommands.GetNewApp("relay_iris", projectRunners.GetRunnerInstance,
		appcommands.CommandDetails{Use: "create", DescShort: "Create relay for iris blockchain", DescLong: "Create relay for iris blockchain"},
		appcommands.CommandDetails{Use: "destroy", DescShort: "Destroy relay for iris blockchain", DescLong: "Destroy relay for iris blockchain"},
		appcommands.CommandDetails{Use: "logs", DescShort: "Tail logs for running relay (iris) instances", DescLong: "Tail logs for running relay (iris) instances"},
		appcommands.CommandDetails{Use: "status", DescShort: "Show current status of currently running relay instances", DescLong: "Show current status of currently running relay instances"},
		appcommands.CommandDetails{Use: "recreate", DescShort: "Recreate end to end relay (iris) instances", DescLong: "Recreate end to end relay (iris) instances"},
		appcommands.CommandDetails{Use: "restart", DescShort: "Restart services for relay (iris) instances", DescLong: "Restart services for relay (iris) instances"},
		appcommands.CommandDetails{Use: "versions", DescShort: "Show available versions for use", DescLong: "Show available versions for use"},

		appcommands.CommandDetails{Use: "show", DescShort: "Show current configuration residing on disk", DescLong: "Show current configuration residing on disk"},
		appcommands.CommandDetails{Use: "diff", DescShort: "Show soft modifications to config staged for apply", DescLong: "Show soft modifications to config staged for apply"},
		appcommands.CommandDetails{Use: "modify", DescShort: "Modify configs on disk", DescLong: "Modify configs on disk"},
		appcommands.CommandDetails{Use: "reset", DescShort: "Reset Configurations on disk", DescLong: "Reset Configurations on disk"},
		appcommands.CommandDetails{Use: "apply", DescShort: "Apply modifications to config", DescLong: "Apply modifications to config"},

		appcommands.CommandDetails{Use: "create", DescShort: "Create keystore", DescLong: "Create keystore"},
		appcommands.CommandDetails{Use: "destroy", DescShort: "Destroy keystore", DescLong: "Destroy keystore"},
	)
	if err != nil {
		log.Error("Error while creating relay_iris application command tree")
		os.Exit(1)
	}

	IrisCmd.AddCommand(app.CreateCmd.Cmd)
	IrisCmd.AddCommand(app.DestroyCmd.Cmd)
	IrisCmd.AddCommand(app.LogsCmd.Cmd)
	IrisCmd.AddCommand(app.StatusCmd.Cmd)
	IrisCmd.AddCommand(app.RecreateCmd.Cmd)
	IrisCmd.AddCommand(app.RestartCmd.Cmd)
	IrisCmd.AddCommand(app.VersionsCmd.Cmd)

	configCmd := &cobra.Command{Use: "config", Short: "Configurations of project set on disk", Long: "Configurations of project set on disk"}
	IrisCmd.AddCommand(configCmd)
	configCmd.AddCommand(app.ConfigShowCmd.Cmd)
	configCmd.AddCommand(app.ConfigDiffCmd.Cmd)
	configCmd.AddCommand(app.ConfigModifyCmd.Cmd)
	configCmd.AddCommand(app.ConfigResetCmd.Cmd)
	configCmd.AddCommand(app.ConfigApplyCmd.Cmd)

	// Extra flag additions for relay_iris -----------------------------------------------

	app.CreateCmd.ArgStore["discovery-addrs"] = app.CreateCmd.Cmd.Flags().StringP("discovery-addrs", "a", "127.0.0.1:8002", "Discovery address of relay")
	app.CreateCmd.ArgStore["heartbeat-addrs"] = app.CreateCmd.Cmd.Flags().StringP("heartbeat-addrs", "g", "127.0.0.1:8003", "Heartbeat address of relay")
	app.CreateCmd.ArgStore["discovery-bind-addr"] = app.CreateCmd.Cmd.Flags().StringP("discovery-bind-addr", "f", "0.0.0.0:22000", "Discovery bind addr")
	app.CreateCmd.ArgStore["pubsub-bind-addr"] = app.CreateCmd.Cmd.Flags().StringP("pubsub-bind-addr", "p", "0.0.0.0:22002", "PubSub bind addr")

	// ----------------------------------------------------------------------------------
}
