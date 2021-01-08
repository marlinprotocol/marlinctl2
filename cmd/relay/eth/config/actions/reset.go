package actions

import (
	"errors"
	"os"

	cmn "github.com/marlinprotocol/ctl2/cmd/relay/eth/common"
	"github.com/marlinprotocol/ctl2/modules/registry"
	"github.com/marlinprotocol/ctl2/modules/util"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppCmd represents the registry command
var ConfigResetCmd = &cobra.Command{
	Use:     "reset",
	Short:   "Reset to default configs",
	Long:    `Reset to default configs`,
	PreRunE: ConfigTest,
	Run: func(cmd *cobra.Command, args []string) {
		var marlinConfig types.Project
		err := viper.UnmarshalKey(types.ProjectID_marlinctl, &marlinConfig)
		if err != nil {
			log.Error("Error while reading marlinctl configs: ", err.Error())
			os.Exit(1)
		}
		log.Debug("Setting up default config for running relay_eth.")
		updPol, ok1 := marlinConfig.AdditionalInfo["defaultprojectupdatepolicy"]
		defRun, ok2 := marlinConfig.AdditionalInfo["defaultprojectruntime"]
		if ok1 && ok2 {
			err = SetupConfiguration(false,
				false,
				updPol.(string),
				defRun.(string),
				"latest")
			if err != nil {
				log.Error("Error while resetting project to default config values ", err.Error())
				os.Exit(1)
			} else {
				log.Info("Successfully reset project to default config values")
			}
		}
		if viper.IsSet(cmn.ProjectID + "_modified") {
			err := util.RemoveConfigEntry(cmn.ProjectID + "_modified")
			if err != nil {
				log.Error("Error while removing modifications relating to the project from config file: " + err.Error())
				os.Exit(1)
			}
		}
	},
}

func SetupConfiguration(enableBeta bool, forceRuntime bool, updatePolicy string, runtime string, version string) error {
	if !util.IsValidUpdatePolicy(updatePolicy) {
		return errors.New("Unknown update policy: " + updatePolicy)
	}

	var projectConfig types.Project

	err := viper.UnmarshalKey(cmn.ProjectID, &projectConfig)
	if err != nil {
		return err
	}

	suitableRuntimes := util.GetRuntimes()

	if !forceRuntime {
		if suitable, ok := suitableRuntimes[runtime]; !ok || !suitable {
			log.Error("Runtime provided for configuration: " + runtime +
				" may not be supported by marlinctl or is not supported by your system." +
				" If you think this is incorrect, override this check using --force-runtime.")
			os.Exit(1)
		} else {
			log.Debug("Runtime provided for configuration: " + runtime +
				" seems to be supported. Going ahead with configuring this.")
		}
	} else {
		log.Warning("Skipped runtime suitability check due to forced runtime")
	}

	var releaseSubscriptions = []string{"public"}

	if enableBeta {
		releaseSubscriptions = append(releaseSubscriptions, "beta")
	}

	var currentVersion = "0.0.0"

	if version != "latest" {
		versions, err := registry.GlobalRegistry.GetVersions(cmn.ProjectID, releaseSubscriptions, version, updatePolicy, runtime)
		if err != nil {
			log.Error("Error while fetching from global registry: ", err)
			os.Exit(1)
		}
		var foundVersion = false
		for _, v := range versions {
			if v.Version == version {
				foundVersion = true
				currentVersion = version
				break
			}
		}
		if !foundVersion {
			log.Error("Version was not found in global registry: ", version)
			os.Exit(1)
		}
	}

	viper.Set(cmn.ProjectID, types.Project{
		Subscription:   releaseSubscriptions,
		UpdatePolicy:   updatePolicy,
		CurrentVersion: currentVersion,
		Storage:        viper.GetString("homedir") + "/projects/" + cmn.ProjectID,
		Runtime:        runtime,
		ForcedRuntime:  false,
		AdditionalInfo: nil,
	})

	return viper.WriteConfig()
}
