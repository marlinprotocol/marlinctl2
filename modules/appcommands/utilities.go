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

package appcommands

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/marlinprotocol/ctl2/modules/keystore"
	"github.com/marlinprotocol/ctl2/modules/registry"
	"github.com/marlinprotocol/ctl2/modules/runner"
	"github.com/marlinprotocol/ctl2/modules/util"

	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (a *app) shallowCopyDescriptions(dst *CommandDetails, src CommandDetails) {
	dst.Use = src.Use
	dst.DescShort = src.DescShort
	dst.DescLong = src.DescLong
	dst.AdditionalPreRunTest = src.AdditionalPreRunTest
}

func (a *app) setupDefaultConfigIfNotExists() error {
	var marlinConfig types.Project
	err := viper.UnmarshalKey(types.ProjectID_marlinctl, &marlinConfig)
	if err != nil {
		return err
	}
	if !viper.IsSet(a.ProjectID) {
		log.Debug("Setting up default config for running relay_eth.")
		updPol, ok1 := marlinConfig.AdditionalInfo["defaultprojectupdatepolicy"]
		defRun, ok2 := marlinConfig.AdditionalInfo["defaultprojectruntime"]
		if ok1 && ok2 {
			err = a.setupConfiguration(false,
				false,
				updPol.(string),
				defRun.(string),
				"latest")
			if err != nil {
				log.Error("Error while seting up default config for project "+a.ProjectID+": ", err)
				os.Exit(1)
			}
		}
	} else {
		log.Debug("Project config found. Not creating defaults.")
	}
	return nil
}

func (a *app) getProjectConfigOrDie() types.Project {
	var projectConfig types.Project
	err := viper.UnmarshalKey(a.ProjectID, &projectConfig)
	if err != nil {
		log.Error("Error while reading project config: ", err)
		os.Exit(1)
	}
	return projectConfig
}

func (a *app) getProjectConfigModOrDie() types.Project {
	var projectConfig types.Project
	err := viper.UnmarshalKey(a.ProjectID+"_modified", &projectConfig)
	if err != nil {
		log.Error("Error while reading project config: ", err)
		os.Exit(1)
	}
	return projectConfig
}

func (a *app) getProjectConfigModOrProjectConfigBase() types.Project {
	var projectConfig types.Project
	if viper.IsSet(a.ProjectID + "_modified") {
		err := viper.UnmarshalKey(a.ProjectID+"_modified", &projectConfig)
		if err != nil {
			log.Error("Error while reading project config (mod): ", err)
			os.Exit(1)
		}
	} else {
		err := viper.UnmarshalKey(a.ProjectID, &projectConfig)
		if err != nil {
			log.Error("Error while reading project config: ", err)
			os.Exit(1)
		}
	}

	return projectConfig
}

func (a *app) getVersionToRunOrDie(updatePolicy string, version string) registry.ProjectVersion {
	versionToRun, err := registry.GlobalRegistry.GetVersionToRun(a.ProjectID, updatePolicy, version)
	if err != nil {
		log.Error("Error while getting version to run for project "+a.ProjectID+": ", err)
		os.Exit(1)
	}
	return versionToRun
}

func (a *app) getResourceMetadata(projectConfig types.Project, instanceId string) (string, string, error) {
	resFileLocation := projectConfig.Storage + "/common/project_" + a.ProjectID + "_instance" + instanceId + ".resource"
	if _, err := os.Stat(resFileLocation); os.IsNotExist(err) {
		return "", "", errors.New("Cannot locate resource: " + resFileLocation)
	}
	file, err := ioutil.ReadFile(resFileLocation)
	if err != nil {
		return "", "", err
	}
	var resourceMetaData = struct {
		Runner  string `json:"Runner"`
		Version string `json:"Version"`
	}{}
	err = json.Unmarshal([]byte(file), &resourceMetaData)
	if err != nil {
		return "", "", err
	}
	log.Debug("Resource metadata: ", resourceMetaData)
	return resourceMetaData.Runner, resourceMetaData.Version, nil
}

func (a *app) getResourceMetadataOrDie(projConfig types.Project, instanceID string) (string, string) {
	runnerId, version, err := a.getResourceMetadata(projConfig, instanceID)
	if err != nil {
		log.Error("Error while getting resource file information for project "+a.ProjectID+" instance "+instanceID+": ", err)
		os.Exit(1)
	}
	return runnerId, version
}

func (a *app) getRunnerInstanceOrDie(runnerId string, version string, storage string, runnerData interface{}, skipRunnerData bool, skipChecksum bool, instanceId string) runner.Runner {
	runner, err := a.RunnerProvider(runnerId, version, storage, runnerData, skipRunnerData, skipChecksum, instanceId)
	if err != nil {
		log.Error("Error while getting project runner to run for project "+a.ProjectID+": ", err)
		os.Exit(1)
	}
	return runner
}

func (a *app) doPreRunSanityOrDie(r runner.Runner) {
	err := r.PreRunSanity()
	if err != nil {
		log.Error("Error while doing prerun sanity for project "+a.ProjectID+": ", err)
		os.Exit(1)
	}
}

func (a *app) doPrepareOrDie(r runner.Runner) {
	err := r.Prepare()
	if err != nil {
		log.Error("Error while doing preparation for project "+a.ProjectID+": ", err)
		os.Exit(1)
	}
}

func (a *app) doCreateOrDie(r runner.Runner, runtimeArgs map[string]string) {
	err := r.Create(runtimeArgs)
	if err != nil {
		log.Error("Error while creating application for project "+a.ProjectID+": ", err)
		os.Exit(1)
	}
}

func (a *app) doDestroyOrDie(r runner.Runner) {
	err := r.Destroy()
	if err != nil {
		log.Error("Error while destroying for project "+a.ProjectID+": ", err)
		os.Exit(1)
	}
}

func (a *app) doStatusOrDie(r runner.Runner) {
	err := r.Status()
	if err != nil {
		log.Error("Error while fetching status for project "+a.ProjectID+": ", err)
		os.Exit(1)
	}
}

func (a *app) doRecreateOrDie(r runner.Runner) {
	err := r.Recreate()
	if err != nil {
		log.Error("Error while recreating for project "+a.ProjectID+": ", err)
		os.Exit(1)
	}
}

func (a *app) doRestartOrDie(r runner.Runner) {
	err := r.Restart()
	if err != nil {
		log.Error("Error while restarting for project "+a.ProjectID+": ", err)
		os.Exit(1)
	}
}

func (a *app) doPostRunOrDie(r runner.Runner) {
	err := r.PostRun()
	if err != nil {
		log.Error("Error while running post run for project "+a.ProjectID+": ", err)
		os.Exit(1)
	}
}

func (a *app) doListVersionsOrDie(projConfig types.Project) {
	versions, err := registry.GlobalRegistry.GetVersions(a.ProjectID, projConfig.Subscription, "0.0.0", "major", projConfig.Runtime)

	if err != nil {
		log.Error("Error encountered while listing versions: ", err)
		os.Exit(1)
	}

	registry.GlobalRegistry.PrettyPrintProjectVersions(versions)
}

func (a *app) doUpdateCurrentVersionOrDie(cfg types.Project) {
	viper.Set(a.ProjectID, cfg)
	err := viper.WriteConfig()
	if err != nil {
		log.Error("Error while updating configuration on disk for project "+a.ProjectID+": ", err)
		os.Exit(1)
	}
}

func (a *app) keystoreSanity() {
	if a.ProjectID == "beacon" || a.ProjectID == "gateway_dot" || a.ProjectID == "gateway_near" || a.ProjectID == "gateway_maticbor" {
		if err := keystore.KeystoreCheck(a.CreateCmd.Cmd, a.ProjectID); err != nil {
			log.Error("keystore error: ", err)
			os.Exit(1)
		}
	}
}

func (c *CommandDetails) getStringFromArgStoreOrDie(key string) string {
	if v, ok := c.ArgStore[key]; ok {
		return *(v.(*string))
	} else {
		log.Error("Cannot find key " + key + " in argstore. Aborting")
		os.Exit(1)
	}
	return ""
}

func (c *CommandDetails) getIntFromArgStoreOrDie(key string) int {
	if v, ok := c.ArgStore[key]; ok {
		return *(v.(*int))
	} else {
		log.Error("Cannot find key " + key + " in argstore. Aborting")
		os.Exit(1)
	}
	return 0
}

func (c *CommandDetails) getStringSliceFromArgStoreOrDie(key string) []string {
	if v, ok := c.ArgStore[key]; ok {
		return *(v.(*[]string))
	} else {
		log.Error("Cannot find key " + key + " in argstore. Aborting")
		os.Exit(1)
	}
	return []string{}
}

func (c *CommandDetails) getBoolFromArgStoreOrDie(key string) bool {
	if v, ok := c.ArgStore[key]; ok {
		return *(v.(*bool))
	} else {
		log.Error("Cannot find key " + key + " in argstore. Aborting")
		os.Exit(1)
	}
	return false
}

func (c *CommandDetails) getStringToStringFromArgStoreOrDie(key string) map[string]string {
	if v, ok := c.ArgStore[key]; ok {
		return *(v.(*map[string]string))
	} else {
		log.Error("Cannot find key " + key + " in argstore. Aborting")
		os.Exit(1)
	}
	return map[string]string{}
}

func (a *app) setupConfiguration(enableBeta bool, forceRuntime bool, updatePolicy string, runtime string, version string) error {
	if !util.IsValidUpdatePolicy(updatePolicy) {
		return errors.New("Unknown update policy: " + updatePolicy)
	}

	var projectConfig types.Project

	err := viper.UnmarshalKey(a.ProjectID, &projectConfig)
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
		versions, err := registry.GlobalRegistry.GetVersions(a.ProjectID, releaseSubscriptions, version, updatePolicy, runtime)
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

	viper.Set(a.ProjectID, types.Project{
		Subscription:   releaseSubscriptions,
		UpdatePolicy:   updatePolicy,
		CurrentVersion: currentVersion,
		Storage:        viper.GetString("homedir") + "/projects/" + a.ProjectID,
		Runtime:        runtime,
		ForcedRuntime:  false,
		AdditionalInfo: nil,
	})

	return viper.WriteConfig()
}
