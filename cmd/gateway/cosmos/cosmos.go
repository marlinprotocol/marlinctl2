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
package cosmos

import (
	"os"

	"github.com/marlinprotocol/ctl2/modules/appcommands"
	"github.com/marlinprotocol/ctl2/modules/keystore"
	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/gateway_cosmos"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var CosmosCmd = &cobra.Command{
	Use:   "cosmos",
	Short: "Cosmos Gateway",
	Long:  `Cosmos Gateway`,
}

func init() {
	app, err := appcommands.GetNewApp("gateway_cosmos", projectRunners.GetRunnerInstance,
		appcommands.CommandDetails{Use: "create", DescShort: "Create gateway for cosmos blockchain", DescLong: "Create gateway for cosmos blockchain"},
		appcommands.CommandDetails{Use: "destroy", DescShort: "Destroy gateway for cosmos blockchain", DescLong: "Destroy gateway for cosmos blockchain"},
		appcommands.CommandDetails{Use: "logs", DescShort: "Tail logs for running gateway (cosmos) instances", DescLong: "Tail logs for running gateway (cosmos) instances"},
		appcommands.CommandDetails{Use: "status", DescShort: "Show status of currently running gateway (cosmos) instances", DescLong: "Show status of currently running gateway (cosmos) instances"},
		appcommands.CommandDetails{Use: "recreate", DescShort: "Recreate end to end gateway (cosmos) instances", DescLong: "Recreate end to end gateway (cosmos) instances"},
		appcommands.CommandDetails{Use: "restart", DescShort: "Restart services for gateway (cosmos) instances", DescLong: "Restart services for gateway (cosmos) instances"},
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
		log.Error("Error while creating gateway_cosmos application command tree")
		os.Exit(1)
	}

	CosmosCmd.AddCommand(app.CreateCmd.Cmd)
	CosmosCmd.AddCommand(app.DestroyCmd.Cmd)
	CosmosCmd.AddCommand(app.LogsCmd.Cmd)
	CosmosCmd.AddCommand(app.StatusCmd.Cmd)
	CosmosCmd.AddCommand(app.RecreateCmd.Cmd)
	CosmosCmd.AddCommand(app.RestartCmd.Cmd)
	CosmosCmd.AddCommand(app.VersionsCmd.Cmd)

	configCmd := &cobra.Command{Use: "config", Short: "Configurations of project set on disk", Long: "Configurations of project set on disk"}
	CosmosCmd.AddCommand(configCmd)
	configCmd.AddCommand(app.ConfigShowCmd.Cmd)
	configCmd.AddCommand(app.ConfigDiffCmd.Cmd)
	configCmd.AddCommand(app.ConfigModifyCmd.Cmd)
	configCmd.AddCommand(app.ConfigResetCmd.Cmd)
	configCmd.AddCommand(app.ConfigApplyCmd.Cmd)

	keystoreCmd := &cobra.Command{Use: "keystore", Short: "Create or Destroy keystore", Long: "Create or Destroy keystore"}
	CosmosCmd.AddCommand(keystoreCmd)
	keystoreCmd.AddCommand(app.KeystoreCreateCmd.Cmd)
	keystoreCmd.AddCommand(app.KeystoreDestroyCmd.Cmd)

	// Extra flag additions for gateway_cosmos -----------------------------------------------
	keystorePath, keystorePassPath, err := keystore.GetKeystoreDetails("gateway_cosmos")
	if err != nil {
		log.Warning("No Keystore for gateway_cosmos")
	}

	app.CreateCmd.ArgStore["discovery-addr"] = app.CreateCmd.Cmd.Flags().StringP("discovery-addr", "d", "0.0.0.0:22002", "Bridge discovery address")
	app.CreateCmd.ArgStore["pubsub-addr"] = app.CreateCmd.Cmd.Flags().StringP("pubsub-addr", "p", "0.0.0.0:22000", "Bridge pubsub address")
	app.CreateCmd.ArgStore["bootstrap-addr"] = app.CreateCmd.Cmd.Flags().StringP("bootstrap-addr", "b", "", "Bridge bootstrap address")
	app.CreateCmd.ArgStore["internal-listen-addr"] = app.CreateCmd.Cmd.Flags().StringP("internal-listen-address", "l", "127.0.0.1:22401", "Bridge listen address")
	app.CreateCmd.ArgStore["keystore-path"] = app.CreateCmd.Cmd.Flags().StringP("keystore-path", "k", keystorePath, "Keystore Path")
	app.CreateCmd.ArgStore["keystore-pass-path"] = app.CreateCmd.Cmd.Flags().StringP("keystore-pass-path", "y", keystorePassPath, "Keystore pass path")
	app.CreateCmd.ArgStore["contracts"] = app.CreateCmd.Cmd.Flags().StringP("contracts", "c", "mainnet", "mainnet/kovan")

	// ----------------------------------------------------------------------------------
}
