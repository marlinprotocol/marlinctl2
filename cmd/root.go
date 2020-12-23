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

	"github.com/marlinprotocol/ctl2/modules/initialise"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	ver "github.com/marlinprotocol/ctl2/version"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var logLevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ctl2.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "loglevel", logLevel, "marlinctl loglevel (default is INFO)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.SetConfigFile(home + "/.marlinctl/marlinctl_config.yaml")
	}

	// viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Error("Cannot read config file: ", viper.ConfigFileUsed())
		log.Info("Seems like you do not have a config file for marlinctl on your machine. Let's create one for you now.")

		initCfg := initialise.InitConfig{}

		err := initCfg.Initialise()
		if err != nil {
			log.Error("Failed to initialise configuration. ", err.Error())
			os.Exit(1)
		} else {
			log.Info("Successfully initialised configuration for marlinctl")
		}
	}
}
