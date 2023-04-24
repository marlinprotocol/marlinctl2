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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/marlinprotocol/ctl2/modules/util"
	"github.com/marlinprotocol/ctl2/types"
	"github.com/marlinprotocol/ctl2/version"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/marlinprotocol/ctl2/modules/registry"
	log "github.com/sirupsen/logrus"

	"github.com/inconshreveable/go-update"
	"github.com/marlinprotocol/ctl2/cmd/beacon"
	"github.com/marlinprotocol/ctl2/cmd/cp"
	"github.com/marlinprotocol/ctl2/cmd/gateway"
	"github.com/marlinprotocol/ctl2/cmd/relay"
)

var cfgFile string
var logLevel string
var skipRegistrySync, forcefulRegistrySync bool
var skipMarlinctlUpdateCheck bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "marlinctl",
	Short: "Marlinctl provides a command line interface for setting up the different components of the Marlin network.",
	Long: `Marlinctl provides a command line interface for setting up the different components of the Marlin network.
It can spawn up beacons, gateways, relays on various platforms and runtimes.`,
	Version: version.RootCmdVersion,
	Run: func(cmd *cobra.Command, args []string) {
		// NOP for upgrade
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		lvl, err := log.ParseLevel(logLevel)
		if err != nil {
			log.Error("Invalid loglevel: ", logLevel)
			os.Exit(1)
		}
		log.SetLevel(lvl)

		// Uncomment following line to show CMD tree
		// showCmdTree(cmd.Root(), "")

		err = readConfig()
		var configuredRegistries []types.Registry
		err = viper.UnmarshalKey("registries", &configuredRegistries)
		if err != nil {
			log.Error("Error reading registries from cfg file: ", err)
			os.Exit(1)
		}
		registry.SetupGlobalRegistry(configuredRegistries)

		currentTime := time.Now().Unix()
		lastSyncTime := viper.GetTime("last_registry_sync").Unix()

		if skipRegistrySync == false && (currentTime-lastSyncTime) > 15*60 || forcefulRegistrySync {
			err = registry.GlobalRegistry.Sync()
			if err != nil {
				log.Error("Error while syncing registry: " + err.Error())
				os.Exit(1)
			}
			viper.Set("last_registry_sync", time.Now())
			err := viper.WriteConfig()
			if err != nil {
				log.Error("Error while writing sync time to state: " + err.Error())
				os.Exit(1)
			}
		} else {
			log.Debug("Skipping registry sync procedure. Metrics: curr: ", currentTime, " lst: ", lastSyncTime, " skip,force: ", skipRegistrySync, " ", forcefulRegistrySync)
		}

		if !skipMarlinctlUpdateCheck {
			hasUpgraded, err := checkMarlinctlUpdates()
			if err != nil {
				log.Error("Error while upgrading marlinctl: " + err.Error())
				os.Exit(1)
			}
			if hasUpgraded {
				log.Info("Starting upgraded marlinctl")
				cmd := exec.Command(os.Args[0], os.Args[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				err = cmd.Start()
				if err != nil {
					fmt.Println("Unable to start updated cli. Retry running same command ", err)
				}
				err = cmd.Wait()
				if err != nil {
					fmt.Println("error occured: ", err)
				}
				os.Exit(2)
			}
		} else {
			log.Debug("Skipping marlinctl update check")
		}

		err = util.ChownRmarlinctlDir()
		if err != nil {
			log.Error("Error while chowning .marlin " + err.Error())
			os.Exit(1)
		}
	},
}

func showCmdTree(cmd *cobra.Command, lvl string) {
	fmt.Println(lvl + "-" + cmd.Use)
	for _, c := range cmd.Commands() {
		showCmdTree(c, lvl+"  |")
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(gateway.GatewayCmd)
	RootCmd.AddCommand(beacon.BeaconCmd)
	RootCmd.AddCommand(relay.RelayCmd)
	RootCmd.AddCommand(cp.CpCmd)

	RootCmd.PersistentFlags().BoolVar(&skipRegistrySync, "skip-sync", false, "skip registry sync during run")
	RootCmd.PersistentFlags().BoolVar(&forcefulRegistrySync, "registry-sync", false, "forceful registry sync from remote. May be used to check for upgrades.")
	RootCmd.PersistentFlags().BoolVar(&skipMarlinctlUpdateCheck, "skip-update-check", false, "skip update check during run")
	RootCmd.PersistentFlags().StringVar(&logLevel, "loglevel", "info", "marlinctl loglevel (default is INFO)")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.marlin/ctl/state.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func readConfig() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := util.GetUser()
		if err != nil {
			return err
		}

		viper.SetConfigFile(home.HomeDir + "/.marlin/ctl/state.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		var cfgVersionOnDisk = viper.GetInt("config_version")
		if cfgVersionOnDisk != version.CfgVersion {
			return errors.New("Cannot use the given config file as it does not match marlinctl's cfgversion. Wanted " + strconv.Itoa(version.CfgVersion) + " but found " + strconv.Itoa(cfgVersionOnDisk))
		}
	} else {
		log.Warning("No config file available on local machine. Creating default for you.")
		err = setupDefaultConfig()
		if err != nil {
			return err
		}
		home, err := util.GetUser()
		if err != nil {
			return err
		}
		viper.SetConfigFile(home.HomeDir + "/.marlin/ctl/state.yaml")
		return viper.ReadInConfig()
	}
	return nil
}

// VIPER defaults ------------------------
var defaultReleaseUpstreams = []types.Registry{
	types.Registry{
		Name:    "public",
		Link:    "https://github.com/marlinprotocol/releases.git",
		Branch:  "public",
		Enabled: true,
	},
	types.Registry{
		Name:    "beta",
		Link:    "https://github.com/marlinprotocol/releases.git",
		Branch:  "beta",
		Enabled: true,
	},
	types.Registry{
		Name:    "alpha",
		Link:    "https://github.com/marlinprotocol/releases.git",
		Branch:  "alpha",
		Enabled: false,
	},
	types.Registry{
		Name:    "dev",
		Link:    "https://github.com/marlinprotocol/releases.git",
		Branch:  "dev",
		Enabled: false,
	},
}

// --------------------------------------

func setupDefaultConfig() error {
	home, err := util.GetUser()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	location := home.HomeDir + "/.marlin/ctl/state.yaml"

	lSplice := strings.Split(location, "/")
	var dirPath string

	for i := 0; i < len(lSplice)-1; i++ {
		dirPath = dirPath + "/" + lSplice[i]
	}

	err = util.CreateDirPathIfNotExists(dirPath)
	if err != nil {
		log.Error("Error while creating directory ", dirPath, " ", err.Error())
	}

	viper.SetConfigFile(location)

	for i := 0; i < len(defaultReleaseUpstreams); i++ {
		defaultReleaseUpstreams[i].Local = home.HomeDir + "/.marlin/ctl/registries/" + defaultReleaseUpstreams[i].Branch
	}

	var defaultProjectRuntime = "linux-amd64.supervisor"
	if runtime.GOOS+"-"+runtime.GOARCH != "linux-amd64" {
		return errors.New("don't know how to service non linux-amd64 system as of now.")
	}

	viper.Set("config_version", version.CfgVersion)
	viper.Set("homedir", home.HomeDir+"/.marlin/ctl/storage")
	viper.Set("registries", defaultReleaseUpstreams)
	viper.Set("marlinctl", types.Project{
		Subscription:   []string{"public"},
		UpdatePolicy:   "minor",
		CurrentVersion: version.ApplicationVersion,
		Storage:        home.HomeDir + "/.marlin/ctl/storage/projects/marlinctl",
		Runtime:        runtime.GOOS + "-" + runtime.GOARCH,
		ForcedRuntime:  false,
		AdditionalInfo: map[string]interface{}{
			"defaultprojectruntime":      defaultProjectRuntime,
			"defaultprojectupdatepolicy": "minor",
		},
	})
	err = viper.WriteConfig()

	if err != nil {
		log.Error("Error while writing config file to ", location, " ", err.Error())
	}

	log.Info("Default marlinctl config written to disk successfully to ", location)
	return nil
}

func checkMarlinctlUpdates() (bool, error) {
	ver, err := registry.GlobalRegistry.GetVersionToRun("marlinctl", "", "")
	if err != nil {
		return false, err
	}
	if version.ApplicationVersion == ver.Version {
		log.Debug("Latest marlinctl described upstream is current marlinctl's version. No updates to do.")
		return false, nil
	}
	log.Info("MarlinCTL needs to upgrade, going from ", version.ApplicationVersion, " to ", ver.Version)

	executableURL := ver.RunnerData.(map[string]interface{})["executable"].(string)
	executableChecksum := ver.RunnerData.(map[string]interface{})["checksum"].(string)
	tempDownloadLoc := "/tmp/marlinctl.tempdownload." + strconv.FormatInt(time.Now().Unix(), 10)

	log.Debug("Downloading marlinctl to ", tempDownloadLoc)

	err = util.DownloadFile(tempDownloadLoc, executableURL)
	if err != nil {
		return false, err
	}

	err = util.VerifyChecksum(tempDownloadLoc, executableChecksum)
	if err != nil {
		return false, err
	}

	log.Debug("Patching start")

	updateFile, err := os.Open(tempDownloadLoc)
	if err != nil {
		return false, nil
	}
	defer updateFile.Close()

	err = update.Apply(updateFile, update.Options{})
	if err != nil {
		return false, err
	}

	log.Debug("Patching complete")

	return true, nil
}
