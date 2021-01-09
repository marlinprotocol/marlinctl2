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
	"github.com/marlinprotocol/ctl2/modules/runner"
	"github.com/spf13/cobra"
)

type CommandDetails struct {
	Use                  string
	DescShort            string
	DescLong             string
	Cmd                  *cobra.Command
	AdditionalPreRunTest func(cmd *cobra.Command, args []string) error
	ArgStore             map[string]interface{}
	ArgStoreSyncer       func()
}

type app struct {
	ProjectID      string
	RunnerProvider func(runnerId string, version string, storage string, runnerData interface{}, skipRunnerData bool, skipChecksum bool, instanceId string) (runner.Runner, error)

	CreateCmd       CommandDetails
	DestroyCmd      CommandDetails
	LogsCmd         CommandDetails
	StatusCmd       CommandDetails
	RecreateCmd     CommandDetails
	RestartCmd      CommandDetails
	VersionsCmd     CommandDetails
	ConfigShowCmd   CommandDetails
	ConfigDiffCmd   CommandDetails
	ConfigModifyCmd CommandDetails
	ConfigResetCmd  CommandDetails
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
			// Run additional argstore syncing procedures
			if a.CreateCmd.ArgStoreSyncer != nil {
				a.CreateCmd.ArgStoreSyncer()
			}

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
			a.doPreRunSanityOrDie(runner)
			a.doPrepareOrDie(runner)
			a.doCreateOrDie(runner, runtimeArgs)
			if version != "" {
				projConfig.CurrentVersion = version
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
			// Run additional argstore syncing procedures
			if a.DestroyCmd.ArgStoreSyncer != nil {
				a.DestroyCmd.ArgStoreSyncer()
			}

			// Extract runtime variables
			instanceID := a.CreateCmd.getStringFromArgStoreOrDie("instance-id")

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
			// Run additional argstore syncing procedures
			if a.LogsCmd.ArgStoreSyncer != nil {
				a.LogsCmd.ArgStoreSyncer()
			}

			// Extract runtime variables
			instanceID := a.LogsCmd.getStringFromArgStoreOrDie("instance-id")

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
			runner.Logs()
		},
	}

	a.LogsCmd.ArgStore = make(map[string]interface{})

	a.LogsCmd.ArgStore["instance-id"] = a.LogsCmd.Cmd.Flags().StringP("instance-id", "i", "001", "instance-id of resource to log")
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
			// Run additional argstore syncing procedures
			if a.StatusCmd.ArgStoreSyncer != nil {
				a.StatusCmd.ArgStoreSyncer()
			}

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
			// Run additional argstore syncing procedures
			if a.RecreateCmd.ArgStoreSyncer != nil {
				a.RecreateCmd.ArgStoreSyncer()
			}

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
			// Run additional argstore syncing procedures
			if a.RestartCmd.ArgStoreSyncer != nil {
				a.RestartCmd.ArgStoreSyncer()
			}

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
			// Run additional argstore syncing procedures
			if a.VersionsCmd.ArgStoreSyncer != nil {
				a.VersionsCmd.ArgStoreSyncer()
			}

			// Run application
			projConfig := a.getProjectConfigOrDie()
			a.doListVersionsOrDie(projConfig)
		},
	}

	a.VersionsCmd.ArgStore = make(map[string]interface{})
}
