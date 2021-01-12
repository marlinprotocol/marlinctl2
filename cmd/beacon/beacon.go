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
package beacon

import (
	"os"

	"github.com/marlinprotocol/ctl2/modules/appcommands"
	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/beacon"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var BeaconCmd = &cobra.Command{
	Use:   "beacon",
	Short: "Marlin Beacon",
	Long:  `Marlin Beacon`,
}

func init() {
	// BeaconCmd.AddCommand(actions.CreateCmd)
	// BeaconCmd.AddCommand(actions.StatusCmd)
	// BeaconCmd.AddCommand(actions.DestroyCmd)
	// BeaconCmd.AddCommand(actions.VersionsCmd)
	// BeaconCmd.AddCommand(actions.LogsCmd)
	app, err := appcommands.GetNewApp("beacon", projectRunners.GetRunnerInstance,
		appcommands.CommandDetails{Use: "create", DescShort: "Create marlin beacon", DescLong: "Create marlin beacon"},
		appcommands.CommandDetails{Use: "destroy", DescShort: "Destroy marlin beacon", DescLong: "Destroy marlin beacon"},
		appcommands.CommandDetails{Use: "logs", DescShort: "Tail logs for running beacon instances", DescLong: "Tail logs for running beacon instances"},
		appcommands.CommandDetails{Use: "status", DescShort: "Show current status of currently running marlin beacon instances", DescLong: "Show current status of currently running marlin beacon instances"},
		appcommands.CommandDetails{Use: "recreate", DescShort: "Recreate end to end marlin beacon instances", DescLong: "Recreate end to end marlin beacon instances"},
		appcommands.CommandDetails{Use: "restart", DescShort: "Restart services for marlin beacon instances", DescLong: "Restart services for marlin beacon instances"},
		appcommands.CommandDetails{Use: "versions", DescShort: "Show available versions for use", DescLong: "Show available versions for use"},

		appcommands.CommandDetails{Use: "show", DescShort: "Show current configuration residing on disk", DescLong: "Show current configuration residing on disk"},
		appcommands.CommandDetails{Use: "diff", DescShort: "Show soft modifications to config staged for apply", DescLong: "Show soft modifications to config staged for apply"},
		appcommands.CommandDetails{Use: "modify", DescShort: "Modify configs on disk", DescLong: "Modify configs on disk"},
		appcommands.CommandDetails{Use: "reset", DescShort: "Reset Configurations on disk", DescLong: "Reset Configurations on disk"},
		appcommands.CommandDetails{Use: "apply", DescShort: "Apply modifications to config", DescLong: "Apply modifications to config"})
	if err != nil {
		log.Error("Error while creating beacon application command tree")
		os.Exit(1)
	}

	BeaconCmd.AddCommand(app.CreateCmd.Cmd)
	BeaconCmd.AddCommand(app.DestroyCmd.Cmd)
	BeaconCmd.AddCommand(app.LogsCmd.Cmd)
	BeaconCmd.AddCommand(app.StatusCmd.Cmd)
	BeaconCmd.AddCommand(app.RecreateCmd.Cmd)
	BeaconCmd.AddCommand(app.RestartCmd.Cmd)
	BeaconCmd.AddCommand(app.VersionsCmd.Cmd)

	configCmd := &cobra.Command{Use: "config", Short: "Configurations of project set on disk", Long: "Configurations of project set on disk"}
	BeaconCmd.AddCommand(configCmd)
	configCmd.AddCommand(app.ConfigShowCmd.Cmd)
	configCmd.AddCommand(app.ConfigDiffCmd.Cmd)
	configCmd.AddCommand(app.ConfigModifyCmd.Cmd)
	configCmd.AddCommand(app.ConfigResetCmd.Cmd)
	configCmd.AddCommand(app.ConfigApplyCmd.Cmd)

	// Extra flag additions for beacon -----------------------------------------------

	app.CreateCmd.ArgStore["discovery-addr"] = app.CreateCmd.Cmd.Flags().StringP("discovery-addr", "a", "127.0.0.1:8002", "Discovery address of beacon")
	app.CreateCmd.ArgStore["heartbeat-addr"] = app.CreateCmd.Cmd.Flags().StringP("heartbeat-addr", "g", "127.0.0.1:8003", "Heartbeat address of beacon")
	app.CreateCmd.ArgStore["bootstrap-addr"] = app.CreateCmd.Cmd.Flags().StringP("bootstrap-addr", "b", "", "Bootstrap address of beacon")
	app.CreateCmd.ArgStore["keystore-path"] = app.CreateCmd.Cmd.Flags().StringP("keystore-path", "k", "", "Keystore path")
	app.CreateCmd.ArgStore["keystore-pass-path"] = app.CreateCmd.Cmd.Flags().StringP("keystore-pass-path", "p", "", "Keystore pass path")

	// ----------------------------------------------------------------------------------
}
