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
package cp

import (
	"os"

	"github.com/marlinprotocol/ctl2/modules/appcommands"

	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/cp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var CpCmd = &cobra.Command{
	Use:   "cp",
	Short: "Marlin Control Plane",
	Long:  `Marlin Control Plane for Oyster(Enclaves Support)`,
}

func init() {
	app, err := appcommands.GetNewApp("cp", projectRunners.GetRunnerInstance,
		appcommands.CommandDetails{Use: "create", DescShort: "Create control plane", DescLong: "Create control plane"},
		appcommands.CommandDetails{Use: "destroy", DescShort: "Destroy control plane", DescLong: "Destroy control plane"},
		appcommands.CommandDetails{Use: "logs", DescShort: "Tail logs for running control plane instances", DescLong: "Tail logs for running control plane instances"},
		appcommands.CommandDetails{Use: "status", DescShort: "Show current status of currently running control plane instances", DescLong: "Show current status of currently running control plane instances"},
		appcommands.CommandDetails{Use: "recreate", DescShort: "Recreate end to end control plane instances", DescLong: "Recreate end to end control plane instances"},
		appcommands.CommandDetails{Use: "restart", DescShort: "Restart services for control plane instances", DescLong: "Restart services for control plane instances"},
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
		log.Error("Error while creating control plane application command tree")
		os.Exit(1)
	}

	CpCmd.AddCommand(app.CreateCmd.Cmd)
	CpCmd.AddCommand(app.DestroyCmd.Cmd)
	CpCmd.AddCommand(app.LogsCmd.Cmd)
	CpCmd.AddCommand(app.StatusCmd.Cmd)
	CpCmd.AddCommand(app.RecreateCmd.Cmd)
	CpCmd.AddCommand(app.RestartCmd.Cmd)
	CpCmd.AddCommand(app.VersionsCmd.Cmd)

	configCmd := &cobra.Command{Use: "config", Short: "Configurations of project set on disk", Long: "Configurations of project set on disk"}
	CpCmd.AddCommand(configCmd)
	configCmd.AddCommand(app.ConfigShowCmd.Cmd)
	configCmd.AddCommand(app.ConfigDiffCmd.Cmd)
	configCmd.AddCommand(app.ConfigModifyCmd.Cmd)
	configCmd.AddCommand(app.ConfigResetCmd.Cmd)
	configCmd.AddCommand(app.ConfigApplyCmd.Cmd)

	keystoreCmd := &cobra.Command{Use: "keystore", Short: "Create or Destroy keystore", Long: "Create or Destroy keystore"}
	CpCmd.AddCommand(keystoreCmd)
	keystoreCmd.AddCommand(app.KeystoreCreateCmd.Cmd)
	keystoreCmd.AddCommand(app.KeystoreDestroyCmd.Cmd)

	// Extra flag additions for cp -----------------------------------------------
	app.CreateCmd.ArgStore["profile"] = app.CreateCmd.Cmd.Flags().StringP("profile", "p", "default", "AWS profile")
	app.CreateCmd.ArgStore["key-name"]  = app.CreateCmd.Cmd.Flags().StringP("key-name", "k", "marlin", "AWS keypair name")
	app.CreateCmd.ArgStore["rpc"]  = app.CreateCmd.Cmd.Flags().StringP("rpc", "u", "", "RPC url")
	app.CreateCmd.ArgStore["regions"] = app.CreateCmd.Cmd.Flags().StringP("regions", "v", "ap-south-1", "Allowed AWS regions")

	// ----------------------------------------------------------------------------------
}
