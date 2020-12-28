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
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	// "github.com/marlinprotocol/ctl2/modules/config"
	ver "github.com/marlinprotocol/ctl2/version"

	"github.com/marlinprotocol/ctl2/cmd/app"
	"github.com/marlinprotocol/ctl2/cmd/config"
	"github.com/marlinprotocol/ctl2/cmd/registry"
)

var cfgFile string
var logLevel string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "marlinctl",
	Short: "Marlinctl provides a command line interface for setting up the different components of the Marlin network.",
	Long: `Marlinctl provides a command line interface for setting up the different components of the Marlin network.
It can spawn up beacons, gateways, relays on various platforms and runtimes. Check out!`,
	Version: ver.RootCmdVersion,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(registry.RegistryCmd)
	RootCmd.AddCommand(config.ConfigCmd)
	RootCmd.AddCommand(app.AppCmd)

	RootCmd.PersistentFlags().StringVar(&logLevel, "loglevel", logLevel, "marlinctl loglevel (default is INFO)")
}
