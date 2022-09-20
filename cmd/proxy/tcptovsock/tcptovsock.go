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

package tcptovsock

import (
	"os"

	"github.com/marlinprotocol/ctl2/modules/appcommands"
	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/proxy_tcptovsock"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var TcpToVsockCmd = &cobra.Command{
	Use:   "tcptovsock",
	Short: "TCP to Vsock Proxy",
	Long:  `TCP to Vsock Proxy`,
}

func init() {
	app, err := appcommands.GetNewApp("tcptovsock", projectRunners.GetRunnerInstance,
		appcommands.CommandDetails{Use: "create", DescShort: "Create TCP to Vsock proxy", DescLong: "Create proxy from local TCP address to local Vsock address"},
		appcommands.CommandDetails{Use: "destroy", DescShort: "Destroy running TCP to Vsock proxy", DescLong: "Create running proxy from local TCP address to local Vsock address"},
		appcommands.CommandDetails{Use: "logs", DescShort: "Tail logs for running TCP to Vsock proxy", DescLong: "Tail logs for running proxy from local TCP address to local Vsock address"},
		appcommands.CommandDetails{Use: "status", DescShort: "Show status of currently running TCP to Vsock proxy", DescLong: "Show status of currently running proxy from local TCP address to local Vsock address"},
		appcommands.CommandDetails{Use: "recreate", DescShort: "Recreate end to end TCP to Vsock proxy", DescLong: "Recreate end to end proxy from local TCP address to local Vsock address"},
		appcommands.CommandDetails{Use: "restart", DescShort: "Restart services for TCP to Vsock proxy", DescLong: "Restart services for proxy from local TCP address to local Vsock address"},
		appcommands.CommandDetails{Use: "versions", DescShort: "Show available versions for use", DescLong: "Show available versions for use"},

		appcommands.CommandDetails{Use: "show", DescShort: "Show current configuration residing on disk", DescLong: "Show current configuration residing on disk"},
		appcommands.CommandDetails{Use: "diff", DescShort: "Show soft modifications to config staged for apply", DescLong: "Show soft modifications to config staged for apply"},
		appcommands.CommandDetails{Use: "modify", DescShort: "Modify configs on disk", DescLong: "Modify configs on disk"},
		appcommands.CommandDetails{Use: "reset", DescShort: "Reset Configurations on disk", DescLong: "Reset Configurations on disk"},
		appcommands.CommandDetails{Use: "apply", DescShort: "Apply modifications to config", DescLong: "Apply modifications to config"},
	)
	if err != nil {
		log.Error("Error while creating TCP to vsock proxy application command tree")
		os.Exit(1)
	}

	TcpToVsockCmd.AddCommand(app.CreateCmd.Cmd)
	TcpToVsockCmd.AddCommand(app.DestroyCmd.Cmd)
	TcpToVsockCmd.AddCommand(app.LogsCmd.Cmd)
	TcpToVsockCmd.AddCommand(app.StatusCmd.Cmd)
	TcpToVsockCmd.AddCommand(app.RecreateCmd.Cmd)
	TcpToVsockCmd.AddCommand(app.RestartCmd.Cmd)
	TcpToVsockCmd.AddCommand(app.VersionsCmd.Cmd)

	configCmd := &cobra.Command{Use: "config", Short: "Configurations of project set on disk", Long: "Configurations of project set on disk"}
	TcpToVsockCmd.AddCommand(configCmd)
	configCmd.AddCommand(app.ConfigShowCmd.Cmd)
	configCmd.AddCommand(app.ConfigDiffCmd.Cmd)
	configCmd.AddCommand(app.ConfigModifyCmd.Cmd)
	configCmd.AddCommand(app.ConfigResetCmd.Cmd)
	configCmd.AddCommand(app.ConfigApplyCmd.Cmd)
	// Extra flag additions for gateway_cosmos -----------------------------------------------

	app.CreateCmd.ArgStore["ip-addr"] = app.CreateCmd.Cmd.Flags().StringP("listner-addr", "l", "127.0.0.1:8000", "Listner TCP port")
	app.CreateCmd.ArgStore["vsock-addr"]  = app.CreateCmd.Cmd.Flags().StringP("server-addr", "s", "0:9001", "Server Vsock address")

	// ----------------------------------------------------------------------------------
}
