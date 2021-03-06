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
package near

import (
	"os"

	"github.com/marlinprotocol/ctl2/modules/appcommands"
	"github.com/marlinprotocol/ctl2/modules/keystore"
	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/gateway_near"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var NearCmd = &cobra.Command{
	Use:   "near",
	Short: "Near Gateway",
	Long:  `Near Gateway`,
}

func init() {
	app, err := appcommands.GetNewApp("gateway_near", projectRunners.GetRunnerInstance,
		appcommands.CommandDetails{Use: "create", DescShort: "Create gateway for near blockchain", DescLong: "Create gateway for near blockchain"},
		appcommands.CommandDetails{Use: "destroy", DescShort: "Destroy gateway for near blockchain", DescLong: "Destroy gateway for near blockchain"},
		appcommands.CommandDetails{Use: "logs", DescShort: "Tail logs for running gateway (near) instances", DescLong: "Tail logs for running gateway (near) instances"},
		appcommands.CommandDetails{Use: "status", DescShort: "Show status of currently running gateway (near) instances", DescLong: "Show status of currently running gateway (near) instances"},
		appcommands.CommandDetails{Use: "recreate", DescShort: "Recreate end to end gateway (near) instances", DescLong: "Recreate end to end gateway (near) instances"},
		appcommands.CommandDetails{Use: "restart", DescShort: "Restart services for gateway (near) instances", DescLong: "Restart services for gateway (near) instances"},
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
		log.Error("Error while creating gateway_near application command tree")
		os.Exit(1)
	}

	NearCmd.AddCommand(app.CreateCmd.Cmd)
	NearCmd.AddCommand(app.DestroyCmd.Cmd)
	NearCmd.AddCommand(app.LogsCmd.Cmd)
	NearCmd.AddCommand(app.StatusCmd.Cmd)
	NearCmd.AddCommand(app.RecreateCmd.Cmd)
	NearCmd.AddCommand(app.RestartCmd.Cmd)
	NearCmd.AddCommand(app.VersionsCmd.Cmd)

	configCmd := &cobra.Command{Use: "config", Short: "Configurations of project set on disk", Long: "Configurations of project set on disk"}
	NearCmd.AddCommand(configCmd)
	configCmd.AddCommand(app.ConfigShowCmd.Cmd)
	configCmd.AddCommand(app.ConfigDiffCmd.Cmd)
	configCmd.AddCommand(app.ConfigModifyCmd.Cmd)
	configCmd.AddCommand(app.ConfigResetCmd.Cmd)
	configCmd.AddCommand(app.ConfigApplyCmd.Cmd)

	keystoreCmd := &cobra.Command{Use: "keystore", Short: "Create or Destroy keystore", Long: "Create or Destroy keystore"}
	NearCmd.AddCommand(keystoreCmd)
	keystoreCmd.AddCommand(app.KeystoreCreateCmd.Cmd)
	keystoreCmd.AddCommand(app.KeystoreDestroyCmd.Cmd)

	// Extra flag additions for gateway_near -----------------------------------------------
	keystorePath, keystorePassPath, _ := keystore.GetKeystoreDetails("gateway_near")

	app.CreateCmd.ArgStore["chain-identity"] = app.CreateCmd.Cmd.Flags().StringP("chain-identity", "a", "gateway_near.key", "Gateway's keystore path")
	app.CreateCmd.ArgStore["listen-addr"] = app.CreateCmd.Cmd.Flags().StringP("listen-addr", "g", "0.0.0.0:21400", "Address on which gateway listens for connections from peer")
	app.CreateCmd.ArgStore["discovery-addr"] = app.CreateCmd.Cmd.Flags().StringP("discovery-addr", "d", "0.0.0.0:21202", "Discovery address")
	app.CreateCmd.ArgStore["pubsub-addr"] = app.CreateCmd.Cmd.Flags().StringP("pubsub-addr", "p", "0.0.0.0:21200", "Pubsub address")
	app.CreateCmd.ArgStore["bootstrap-addr"] = app.CreateCmd.Cmd.Flags().StringP("bootstrap-addr", "b", "", "Bootstrap address")
	app.CreateCmd.ArgStore["keystore-path"] = app.CreateCmd.Cmd.Flags().StringP("keystore-path", "k", keystorePath, "Keystore Path")
	app.CreateCmd.ArgStore["keystore-pass-path"] = app.CreateCmd.Cmd.Flags().StringP("keystore-pass-path", "y", keystorePassPath, "Keystore pass path")
	app.CreateCmd.ArgStore["contracts"] = app.CreateCmd.Cmd.Flags().StringP("contracts", "c", "mainnet", "mainnet/kovan")

	// ----------------------------------------------------------------------------------
}
