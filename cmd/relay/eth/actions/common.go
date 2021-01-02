package actions

import (
	"errors"
	"os"

	"github.com/marlinprotocol/ctl2/modules/registry"
	"github.com/marlinprotocol/ctl2/modules/util"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	enableBeta   bool
	runtime      string
	version      string
	instanceId   string
	projectId    string = "relay_eth"
	updatePolicy string
	skipChecksum bool
	forceRuntime bool
	runtimeArgs  map[string]string
)

var ConfigTest = func(cmd *cobra.Command, args []string) error {
	var marlinConfig types.Project
	err := viper.UnmarshalKey("marlinctl", &marlinConfig)
	if err != nil {
		return err
	}
	if !viper.IsSet("relay_eth") {
		log.Debug("Setting up default config for running relay_eth.")
		updPol, ok1 := marlinConfig.AdditionalInfo["defaultprojectupdatepolicy"]
		defRun, ok2 := marlinConfig.AdditionalInfo["defaultprojectruntime"]
		if ok1 && ok2 {
			setupConfiguration(false,
				false,
				updPol.(string),
				defRun.(string),
				"latest")
		}
	} else {
		log.Debug("Project config found. Not creating defaults.")
	}
	return nil
}

func setupConfiguration(enableBeta bool, forceRuntime bool, updatePolicy string, runtime string, version string) error {
	if !util.IsValidUpdatePolicy(updatePolicy) {
		return errors.New("Unknown update policy: " + updatePolicy)
	}

	var projectConfig types.Project

	err := viper.UnmarshalKey(projectId, &projectConfig)
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
		versions, err := registry.GlobalRegistry.GetVersions(projectId, releaseSubscriptions, version, updatePolicy, runtime)
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

	viper.Set(projectId, types.Project{
		Subscription:   releaseSubscriptions,
		UpdatePolicy:   updatePolicy,
		CurrentVersion: currentVersion,
		Storage:        viper.GetString("homedir") + "/projects/" + projectId,
		Runtime:        runtime,
		ForcedRuntime:  false,
	})

	return viper.WriteConfig()
}
