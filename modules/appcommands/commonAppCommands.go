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
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/marlinprotocol/ctl2/modules/keystore"
	"github.com/marlinprotocol/ctl2/modules/runner"
	"github.com/marlinprotocol/ctl2/modules/util"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CommandDetails struct {
	Use                  string
	DescShort            string
	DescLong             string
	Cmd                  *cobra.Command
	AdditionalPreRunTest func(cmd *cobra.Command, args []string) error
	ArgStore             map[string]interface{}
}

type app struct {
	ProjectID      string
	RunnerProvider func(runnerId string, version string, storage string, runnerData interface{}, skipRunnerData bool, skipChecksum bool, instanceId string) (runner.Runner, error)

	CreateCmd          CommandDetails
	DestroyCmd         CommandDetails
	LogsCmd            CommandDetails
	StatusCmd          CommandDetails
	RecreateCmd        CommandDetails
	RestartCmd         CommandDetails
	VersionsCmd        CommandDetails
	ConfigShowCmd      CommandDetails
	ConfigDiffCmd      CommandDetails
	ConfigModifyCmd    CommandDetails
	ConfigResetCmd     CommandDetails
	ConfigApplyCmd     CommandDetails
	KeystoreCreateCmd  CommandDetails
	KeystoreDestroyCmd CommandDetails
}

// Write Defaults logic

// Write initialiser logic
func GetNewApp(_projectID string,
	_runnerProvider func(runnerId string, version string, storage string, runnerData interface{}, skipRunnerData bool, skipChecksum bool, instanceId string) (runner.Runner, error),
	_createCmd CommandDetails,
	_destroyCmd CommandDetails,
	_logsCmd CommandDetails,
	_statusCmd CommandDetails,
	_recreateCmd CommandDetails,
	_restartCmd CommandDetails,
	_versionsCmd CommandDetails,
	_configShowCmd CommandDetails,
	_configDiffCmd CommandDetails,
	_configModifyCmd CommandDetails,
	_configResetCmd CommandDetails,
	_configApplyCmd CommandDetails,
	_keystoreCreateCmd CommandDetails,
	_keystoreDestroyCmd CommandDetails,
) (app, error) {
	createdApp := app{
		ProjectID:      _projectID,
		RunnerProvider: _runnerProvider,
	}

	createdApp.shallowCopyDescriptions(&createdApp.CreateCmd, _createCmd)
	createdApp.setupCreateCommand()

	createdApp.shallowCopyDescriptions(&createdApp.DestroyCmd, _destroyCmd)
	createdApp.setupDestroyCommand()

	createdApp.shallowCopyDescriptions(&createdApp.LogsCmd, _logsCmd)
	createdApp.setupLogsCommand()

	createdApp.shallowCopyDescriptions(&createdApp.StatusCmd, _statusCmd)
	createdApp.setupStatusCommand()

	createdApp.shallowCopyDescriptions(&createdApp.RecreateCmd, _recreateCmd)
	createdApp.setupRecreateCommand()

	createdApp.shallowCopyDescriptions(&createdApp.RestartCmd, _restartCmd)
	createdApp.setupRestartCommand()

	createdApp.shallowCopyDescriptions(&createdApp.VersionsCmd, _versionsCmd)
	createdApp.setupVersionsCommand()

	createdApp.shallowCopyDescriptions(&createdApp.ConfigShowCmd, _configShowCmd)
	createdApp.setupConfigShowCommand()

	createdApp.shallowCopyDescriptions(&createdApp.ConfigDiffCmd, _configDiffCmd)
	createdApp.setupConfigDiffCommand()

	createdApp.shallowCopyDescriptions(&createdApp.ConfigModifyCmd, _configModifyCmd)
	createdApp.setupConfigModifyCommand()

	createdApp.shallowCopyDescriptions(&createdApp.ConfigResetCmd, _configResetCmd)
	createdApp.setupConfigResetCommand()

	createdApp.shallowCopyDescriptions(&createdApp.ConfigApplyCmd, _configApplyCmd)
	createdApp.setupConfigApplyCommand()

	createdApp.shallowCopyDescriptions(&createdApp.KeystoreCreateCmd, _keystoreCreateCmd)
	createdApp.setupKeystoreCreateCommand()

	createdApp.shallowCopyDescriptions(&createdApp.KeystoreDestroyCmd, _keystoreDestroyCmd)
	createdApp.setupKeystoreDestroyCommand()

	return createdApp, nil
}

// Create command
func (a *app) setupCreateCommand() {
	a.CreateCmd.Cmd = &cobra.Command{
		Use:   a.CreateCmd.Use,
		Short: a.CreateCmd.DescShort,
		Long:  a.CreateCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.CreateCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			a.keystoreSanity()
			// Extract runtime variables
			version := a.CreateCmd.getStringFromArgStoreOrDie("version")
			instanceID := a.CreateCmd.getStringFromArgStoreOrDie("instance-id")
			skipChecksum := a.CreateCmd.getBoolFromArgStoreOrDie("skip-checksum")
			runtimeArgs := a.CreateCmd.getStringToStringFromArgStoreOrDie("runtime-args")

			// Run application
			projConfig := a.getProjectConfigOrDie()
			versionToRun := a.getVersionToRunOrDie(projConfig.UpdatePolicy, version)
			runner := a.getRunnerInstanceOrDie(versionToRun.RunnerId,
				versionToRun.Version,
				projConfig.Storage,
				versionToRun.RunnerData,
				false,
				skipChecksum,
				instanceID)

			// MESSY SUBSTITUTIONS
			a.beaconCreateSusbstitutions(versionToRun.RunnerId)
			a.relayEthCreateSubstitutions(versionToRun.RunnerId)
			a.gatewayDotCreateSubstitutions(versionToRun.RunnerId)
			a.gatewayNearCreateSubstitutions(versionToRun.RunnerId)
			a.gatewayIrisCreateSubstitutions(versionToRun.RunnerId)
			a.gatewayCosmosCreateSubstitutions(versionToRun.RunnerId)
			a.relayIrisCreateSubstitutions(versionToRun.RunnerId)
			a.relayCosmosCreateSubstitutions(versionToRun.RunnerId)

			a.doPreRunSanityOrDie(runner)
			a.doPrepareOrDie(runner)
			a.doCreateOrDie(runner, runtimeArgs)
			if version == "" {
				projConfig.CurrentVersion = versionToRun.Version
				a.doUpdateCurrentVersionOrDie(projConfig)
			}
		},
	}

	a.CreateCmd.ArgStore = make(map[string]interface{})

	a.CreateCmd.ArgStore["version"] = a.CreateCmd.Cmd.Flags().StringP("version", "x", "", "runtime version override")
	a.CreateCmd.ArgStore["instance-id"] = a.CreateCmd.Cmd.Flags().StringP("instance-id", "i", "001", "instance-id of spawned up resource")
	a.CreateCmd.ArgStore["skip-checksum"] = a.CreateCmd.Cmd.Flags().BoolP("skip-checksum", "s", false, "skip checksum verification while starting up binaries")
	a.CreateCmd.ArgStore["runtime-args"] = a.CreateCmd.Cmd.Flags().StringToStringP("runtime-args", "r", map[string]string{}, "runtime arguments while starting up")
}

// Destroy command
func (a *app) setupDestroyCommand() {
	a.DestroyCmd.Cmd = &cobra.Command{
		Use:   a.DestroyCmd.Use,
		Short: a.DestroyCmd.DescShort,
		Long:  a.DestroyCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.DestroyCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Extract runtime variables
			instanceID := a.DestroyCmd.getStringFromArgStoreOrDie("instance-id")

			// Run application
			projConfig := a.getProjectConfigOrDie()
			runnerID, version := a.getResourceMetadataOrDie(projConfig, instanceID)
			runner := a.getRunnerInstanceOrDie(runnerID,
				version,
				projConfig.Storage,
				struct{}{},
				true,
				true,
				instanceID)
			a.doPreRunSanityOrDie(runner)
			a.doDestroyOrDie(runner)
			a.doPostRunOrDie(runner)
		},
	}

	a.DestroyCmd.ArgStore = make(map[string]interface{})

	a.DestroyCmd.ArgStore["instance-id"] = a.DestroyCmd.Cmd.Flags().StringP("instance-id", "i", "001", "instance-id of resource to destroy")
}

// Logs command
func (a *app) setupLogsCommand() {
	a.LogsCmd.Cmd = &cobra.Command{
		Use:   a.LogsCmd.Use,
		Short: a.LogsCmd.DescShort,
		Long:  a.LogsCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.LogsCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Extract runtime variables
			instanceID := a.LogsCmd.getStringFromArgStoreOrDie("instance-id")

			// Run application
			projConfig := a.getProjectConfigOrDie()
			runnerID, version := a.getResourceMetadataOrDie(projConfig, instanceID)
			last := a.LogsCmd.getIntFromArgStoreOrDie("last")
			runner := a.getRunnerInstanceOrDie(runnerID,
				version,
				projConfig.Storage,
				struct{}{},
				true,
				true,
				instanceID)
			a.doPreRunSanityOrDie(runner)
			runner.Logs(last)
		},
	}

	a.LogsCmd.ArgStore = make(map[string]interface{})

	a.LogsCmd.ArgStore["instance-id"] = a.LogsCmd.Cmd.Flags().StringP("instance-id", "i", "001", "instance-id of resource to log")
	a.LogsCmd.ArgStore["last"] = a.LogsCmd.Cmd.Flags().IntP("last", "n", 100, "number of last lines to tail in logfile")
}

// Status command
func (a *app) setupStatusCommand() {
	a.StatusCmd.Cmd = &cobra.Command{
		Use:   a.StatusCmd.Use,
		Short: a.StatusCmd.DescShort,
		Long:  a.StatusCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.StatusCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Extract runtime variables
			instanceID := a.StatusCmd.getStringFromArgStoreOrDie("instance-id")

			// Run application
			projConfig := a.getProjectConfigOrDie()
			runnerID, version := a.getResourceMetadataOrDie(projConfig, instanceID)
			runner := a.getRunnerInstanceOrDie(runnerID,
				version,
				projConfig.Storage,
				struct{}{},
				true,
				true,
				instanceID)
			a.doPreRunSanityOrDie(runner)
			a.doStatusOrDie(runner)
		},
	}

	a.StatusCmd.ArgStore = make(map[string]interface{})

	a.StatusCmd.ArgStore["instance-id"] = a.StatusCmd.Cmd.Flags().StringP("instance-id", "i", "001", "instance-id of resource to find status of")
}

// Recreate command
func (a *app) setupRecreateCommand() {
	a.RecreateCmd.Cmd = &cobra.Command{
		Use:   a.RecreateCmd.Use,
		Short: a.RecreateCmd.DescShort,
		Long:  a.RecreateCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.RecreateCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Extract runtime variables
			instanceID := a.RecreateCmd.getStringFromArgStoreOrDie("instance-id")

			// Run application
			projConfig := a.getProjectConfigOrDie()
			runnerID, version := a.getResourceMetadataOrDie(projConfig, instanceID)
			runner := a.getRunnerInstanceOrDie(runnerID,
				version,
				projConfig.Storage,
				struct{}{},
				true,
				true,
				instanceID)
			a.doPreRunSanityOrDie(runner)
			a.doRecreateOrDie(runner)
		},
	}

	a.RecreateCmd.ArgStore = make(map[string]interface{})

	a.RecreateCmd.ArgStore["instance-id"] = a.RecreateCmd.Cmd.Flags().StringP("instance-id", "i", "001", "instance-id of resource to recreate")
}

// Restart command
func (a *app) setupRestartCommand() {
	a.RestartCmd.Cmd = &cobra.Command{
		Use:   a.RestartCmd.Use,
		Short: a.RestartCmd.DescShort,
		Long:  a.RestartCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.RestartCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Extract runtime variables
			instanceID := a.RestartCmd.getStringFromArgStoreOrDie("instance-id")

			// Run application
			projConfig := a.getProjectConfigOrDie()
			runnerID, version := a.getResourceMetadataOrDie(projConfig, instanceID)
			runner := a.getRunnerInstanceOrDie(runnerID,
				version,
				projConfig.Storage,
				struct{}{},
				true,
				true,
				instanceID)
			a.doPreRunSanityOrDie(runner)
			a.doRestartOrDie(runner)
		},
	}

	a.RestartCmd.ArgStore = make(map[string]interface{})

	a.RestartCmd.ArgStore["instance-id"] = a.RestartCmd.Cmd.Flags().StringP("instance-id", "i", "001", "instance-id of resource to restart")
}

// Versions command
func (a *app) setupVersionsCommand() {
	a.VersionsCmd.Cmd = &cobra.Command{
		Use:   a.VersionsCmd.Use,
		Short: a.VersionsCmd.DescShort,
		Long:  a.VersionsCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.VersionsCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Run application
			projConfig := a.getProjectConfigOrDie()
			a.doListVersionsOrDie(projConfig)
		},
	}

	a.VersionsCmd.ArgStore = make(map[string]interface{})
}

// Config Show command
func (a *app) setupConfigShowCommand() {
	a.ConfigShowCmd.Cmd = &cobra.Command{
		Use:   a.ConfigShowCmd.Use,
		Short: a.ConfigShowCmd.DescShort,
		Long:  a.ConfigShowCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.ConfigShowCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Run application
			projConfig := a.getProjectConfigOrDie()
			s, err := json.MarshalIndent(projConfig, "", "  ")
			if err != nil {
				log.Error("Error while decoding json: ", err.Error())
				os.Exit(1)
			}
			log.Info("Current config:")
			util.PrintPrettyDiff(string(s))
		},
	}

	a.ConfigShowCmd.ArgStore = make(map[string]interface{})
}

// Config Diff command
func (a *app) setupConfigDiffCommand() {
	a.ConfigDiffCmd.Cmd = &cobra.Command{
		Use:   a.ConfigDiffCmd.Use,
		Short: a.ConfigDiffCmd.DescShort,
		Long:  a.ConfigDiffCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.ConfigDiffCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Run application
			if !viper.IsSet(a.ProjectID + "_modified") {
				log.Error("No modifications on disk to show diff with")
				os.Exit(1)
			}
			projConfig := a.getProjectConfigOrDie()
			s, err := json.MarshalIndent(projConfig, "", "  ")
			if err != nil {
				log.Error("Error while decoding json: ", err.Error())
				os.Exit(1)
			}
			projConfigMod := a.getProjectConfigModOrDie()
			smod, err := json.MarshalIndent(projConfigMod, "", "  ")
			if err != nil {
				log.Error("Error while decoding json (mod): ", err.Error())
				os.Exit(1)
			}
			log.Info("Difference:")
			util.PrintPrettyDiff(cmp.Diff(string(s), string(smod)))
		},
	}

	a.ConfigDiffCmd.ArgStore = make(map[string]interface{})
}

// Config Modify command
func (a *app) setupConfigModifyCommand() {
	a.ConfigModifyCmd.Cmd = &cobra.Command{
		Use:   a.ConfigModifyCmd.Use,
		Short: a.ConfigModifyCmd.DescShort,
		Long:  a.ConfigModifyCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.ConfigModifyCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Extract runtime variables
			subscriptions := a.ConfigModifyCmd.getStringSliceFromArgStoreOrDie("subscriptions")
			updatePolicy := a.ConfigModifyCmd.getStringFromArgStoreOrDie("update-policy")
			currentVersion := a.ConfigModifyCmd.getStringFromArgStoreOrDie("current-version")
			storage := a.ConfigModifyCmd.getStringFromArgStoreOrDie("storage")
			runtime := a.ConfigModifyCmd.getStringFromArgStoreOrDie("runtime")
			forceRuntime := a.ConfigModifyCmd.getBoolFromArgStoreOrDie("force-runtime")

			// Run application
			projectConfigMod := a.getProjectConfigModOrProjectConfigBase()

			if len(subscriptions) > 0 {
				for _, v := range subscriptions {
					if !util.IsValidSubscription(v) {
						log.Error("Not a valid subscription line: ", v)
						os.Exit(1)
					}
				}
				projectConfigMod.Subscription = subscriptions
			}

			if updatePolicy != "" {
				if !util.IsValidUpdatePolicy(updatePolicy) {
					log.Error("Not a valid update policy: ", updatePolicy)
					os.Exit(1)
				}
				projectConfigMod.UpdatePolicy = updatePolicy
			}

			if currentVersion != "" {
				_, _, _, _, _, err := util.DecodeVersionString(currentVersion)
				if err != nil {
					log.Error("Error decoding version string: ", err.Error())
					os.Exit(1)
				}
				projectConfigMod.CurrentVersion = currentVersion
			}

			if storage != "" {
				projectConfigMod.Storage = storage
			}

			if runtime != "" {
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
				projectConfigMod.Runtime = runtime
				projectConfigMod.ForcedRuntime = forceRuntime
			}

			viper.Set(a.ProjectID+"_modified", projectConfigMod)

			err := viper.WriteConfig()
			if err != nil {
				log.Error("Error while writing staging configs to disk: ", err.Error())
				os.Exit(1)
			}
			log.Info("Modifications registered")
		},
	}

	a.ConfigModifyCmd.ArgStore = make(map[string]interface{})

	a.ConfigModifyCmd.ArgStore["subscriptions"] = a.ConfigModifyCmd.Cmd.Flags().StringSliceP("subscriptions", "s", []string{}, "Release channels to subscribe to")
	a.ConfigModifyCmd.ArgStore["update-policy"] = a.ConfigModifyCmd.Cmd.Flags().StringP("update-policy", "u", "", "Update policy to set")
	a.ConfigModifyCmd.ArgStore["current-version"] = a.ConfigModifyCmd.Cmd.Flags().StringP("current-version", "c", "", "Current version to set for executables")
	a.ConfigModifyCmd.ArgStore["storage"] = a.ConfigModifyCmd.Cmd.Flags().StringP("storage", "l", "", "Storage location")
	a.ConfigModifyCmd.ArgStore["runtime"] = a.ConfigModifyCmd.Cmd.Flags().StringP("runtime", "r", "", "Runtime to use")
	a.ConfigModifyCmd.ArgStore["force-runtime"] = a.ConfigModifyCmd.Cmd.Flags().BoolP("force-runtime", "f", false, "Forcefully set runtime")
}

// Config Reset command
func (a *app) setupConfigResetCommand() {
	a.ConfigResetCmd.Cmd = &cobra.Command{
		Use:   a.ConfigResetCmd.Use,
		Short: a.ConfigResetCmd.DescShort,
		Long:  a.ConfigResetCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.ConfigResetCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Run application
			err := util.RemoveConfigEntry(a.ProjectID)
			if err != nil {
				log.Error("Error while removing project config entry: ", a.ProjectID)
				os.Exit(1)
			}
			err = util.RemoveConfigEntry(a.ProjectID + "_modified")
			if err != nil {
				log.Error("Error while removing project config entry: ", a.ProjectID+"_modified")
				os.Exit(1)
			}
			err = a.setupDefaultConfigIfNotExists()
			if err != nil {
				log.Error("Error while setting up default config on disk: ", err)
				os.Exit(1)
			}
		},
	}

	a.ConfigResetCmd.ArgStore = make(map[string]interface{})
}

// Config Apply command
func (a *app) setupConfigApplyCommand() {
	a.ConfigApplyCmd.Cmd = &cobra.Command{
		Use:   a.ConfigApplyCmd.Use,
		Short: a.ConfigApplyCmd.DescShort,
		Long:  a.ConfigApplyCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.ConfigApplyCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Run application
			if !viper.IsSet(a.ProjectID + "_modified") {
				log.Error("No modifications on disk to apply")
				os.Exit(1)
			}
			var projectConfig types.Project
			err := viper.UnmarshalKey(a.ProjectID+"_modified", &projectConfig)
			if err != nil {
				log.Error("Error while reading project configs: ", err.Error())
				os.Exit(1)
			}

			viper.Set(a.ProjectID, projectConfig)

			err = viper.WriteConfig()
			if err != nil {
				log.Error("Error while writing configs to disk: ", err.Error())
				os.Exit(1)
			}
			err = util.RemoveConfigEntry(a.ProjectID + "_modified")
			if err != nil {
				log.Error("Error while removing project config entry: ", a.ProjectID+"_modified")
				os.Exit(1)
			}
		},
	}

	a.ConfigApplyCmd.ArgStore = make(map[string]interface{})
}

func (a *app) setupKeystoreCreateCommand() {
	a.KeystoreCreateCmd.Cmd = &cobra.Command{
		Use:   a.KeystoreCreateCmd.Use,
		Short: a.KeystoreCreateCmd.DescShort,
		Long:  a.KeystoreCreateCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.KeystoreCreateCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var passphrase string
			if !a.KeystoreCreateCmd.Cmd.Flags().Changed("pass-path") {
				// read from stdin
				log.Info("Enter passphrase to generate keystore.")
				var err error
				passphrase, err = util.ReadInputPasswordLine()
				if err != nil {
					log.Error("Error while reading passphrase", err)
					os.Exit(1)
				}
			} else {
				keystorePassPath := a.KeystoreCreateCmd.getStringFromArgStoreOrDie("pass-path")
				var err error
				passphrase, err = util.ReadStringFromFile(keystorePassPath)
				if err != nil {
					log.Error("Error while reading passphrase file", err)
					os.Exit(1)
				}
			}

			home, err := util.GetUser()
			if err == nil {
				keystoreDir := home.HomeDir + "/.marlin/ctl/storage/projects/" + a.ProjectID + "/common/keystore"
				err = keystore.Create(keystoreDir, passphrase)
			}
			if err != nil {
				log.Error("Error while creating keystore for project "+a.ProjectID+": ", err)
				os.Exit(1)
			}
		},
	}

	a.KeystoreCreateCmd.ArgStore = make(map[string]interface{})
	a.KeystoreCreateCmd.ArgStore["pass-path"] = a.KeystoreCreateCmd.Cmd.Flags().StringP("pass-path", "p", "", "path to the passphrase file")
}

func (a *app) setupKeystoreDestroyCommand() {
	a.KeystoreDestroyCmd.Cmd = &cobra.Command{
		Use:   a.KeystoreDestroyCmd.Use,
		Short: a.KeystoreDestroyCmd.DescShort,
		Long:  a.KeystoreDestroyCmd.DescLong,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			additionalTest := a.KeystoreDestroyCmd.AdditionalPreRunTest
			err := a.setupDefaultConfigIfNotExists()
			if err != nil {
				return err
			} else if err == nil && additionalTest != nil {
				return additionalTest(cmd, args)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			home, err := util.GetUser()
			if err == nil {
				keystoreDir := home.HomeDir + "/.marlin/ctl/storage/projects/" + a.ProjectID + "/common/keystore"
				err = keystore.Destroy(keystoreDir)
			}
			if err != nil {
				log.Error("Error while destroying keystore for project "+a.ProjectID+": ", err)
				os.Exit(1)
			}
		},
	}
}
