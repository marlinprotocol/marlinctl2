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

package vsocktoip

import (
	"os"

	"github.com/marlinprotocol/ctl2/modules/appcommands"
	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/vsock_to_ip"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var VsockToIpCmd = &cobra.Command{
	Use:   "vsocktoip",
	Short: "Vsock to IP Proxy",
	Long:  `Vsock to IP Proxy`,
}

func init() {
	app, err := appcommands.GetNewApp("vsocktoip", projectRunners.GetRunnerInstance,
		appcommands.CommandDetails{Use: "create", DescShort: "Create Vsock to IP proxy", DescLong: "Create proxy from local Vsock address to local IP address"},
		appcommands.CommandDetails{Use: "destroy", DescShort: "Destroy running Vsock to IP proxy", DescLong: "Create running proxy from local Vsock address to local IP address"},
		appcommands.CommandDetails{Use: "logs", DescShort: "Tail logs for running Vsock to IP proxy", DescLong: "Tail logs for running proxy from local Vsock address to local IP address"},
		appcommands.CommandDetails{Use: "status", DescShort: "Show status of currently running Vsock to IP proxy", DescLong: "Show status of currently running proxy from local Vsock address to local IP address"},
		appcommands.CommandDetails{Use: "recreate", DescShort: "Recreate end to end Vsock to IP proxy", DescLong: "Recreate end to end proxy from local Vsock address to local IP address"},
		appcommands.CommandDetails{Use: "restart", DescShort: "Restart services for Vsock to IP proxy", DescLong: "Restart services for proxy from local Vsock address to local IP address"},
		appcommands.CommandDetails{Use: "versions", DescShort: "Show available versions for use", DescLong: "Show available versions for use"},

		appcommands.CommandDetails{Use: "show", DescShort: "Show current configuration residing on disk", DescLong: "Show current configuration residing on disk"},
		appcommands.CommandDetails{Use: "diff", DescShort: "Show soft modifications to config staged for apply", DescLong: "Show soft modifications to config staged for apply"},
		appcommands.CommandDetails{Use: "modify", DescShort: "Modify configs on disk", DescLong: "Modify configs on disk"},
		appcommands.CommandDetails{Use: "reset", DescShort: "Reset Configurations on disk", DescLong: "Reset Configurations on disk"},
		appcommands.CommandDetails{Use: "apply", DescShort: "Apply modifications to config", DescLong: "Apply modifications to config"},
	)
	if err != nil {
		log.Error("Error while creating vsock to ip proxy application command tree")
		os.Exit(1)
	}

	VsockToIpCmd.AddCommand(app.CreateCmd.Cmd)
	VsockToIpCmd.AddCommand(app.DestroyCmd.Cmd)
	VsockToIpCmd.AddCommand(app.LogsCmd.Cmd)
	VsockToIpCmd.AddCommand(app.StatusCmd.Cmd)
	VsockToIpCmd.AddCommand(app.RecreateCmd.Cmd)
	VsockToIpCmd.AddCommand(app.RestartCmd.Cmd)
	VsockToIpCmd.AddCommand(app.VersionsCmd.Cmd)

	configCmd := &cobra.Command{Use: "config", Short: "Configurations of project set on disk", Long: "Configurations of project set on disk"}
	VsockToIpCmd.AddCommand(configCmd)
	configCmd.AddCommand(app.ConfigShowCmd.Cmd)
	configCmd.AddCommand(app.ConfigDiffCmd.Cmd)
	configCmd.AddCommand(app.ConfigModifyCmd.Cmd)
	configCmd.AddCommand(app.ConfigResetCmd.Cmd)
	configCmd.AddCommand(app.ConfigApplyCmd.Cmd)
	// Extra flag additions for gateway_cosmos -----------------------------------------------

	app.CreateCmd.ArgStore["listner-addr"] = app.CreateCmd.Cmd.Flags().StringP("listner-addr", "l", "0:8000", "Listner Vsock port")
	app.CreateCmd.ArgStore["server-addr"]  = app.CreateCmd.Cmd.Flags().StringP("server-addr", "s", "9001", "Server IP address")

	// ----------------------------------------------------------------------------------
}
