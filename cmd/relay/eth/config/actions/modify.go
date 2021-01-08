package actions

import (
	"os"

	"github.com/marlinprotocol/ctl2/modules/util"

	"github.com/getlantern/deepcopy"
	cmn "github.com/marlinprotocol/ctl2/cmd/relay/eth/common"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppCmd represents the registry command
var ConfigModifyCmd = &cobra.Command{
	Use:     "modify",
	Short:   "Modify state information of the project",
	Long:    `Modify state information of the project`,
	PreRunE: ConfigTest,
	Run: func(cmd *cobra.Command, args []string) {
		var projectConfig types.Project
		err := viper.UnmarshalKey(cmn.ProjectID, &projectConfig)
		if err != nil {
			log.Error("Error while reading project configs: ", err.Error())
			os.Exit(1)
		}

		modifiedProjectID := cmn.ProjectID + "_modified"
		var projectConfigMod types.Project
		if viper.IsSet(modifiedProjectID) {
			err = viper.UnmarshalKey(modifiedProjectID, &projectConfigMod)
			if err != nil {
				log.Error("No existing modifications found.")
				os.Exit(1)
			}
		} else {
			log.Info("No existing modifications found. Sensing from set project config.")
			deepcopy.Copy(&projectConfigMod, &projectConfig)
		}

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
			_, _, _, _, _, err = util.DecodeVersionString(currentVersion)
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
		viper.Set(modifiedProjectID, projectConfigMod)

		err = viper.WriteConfig()
		if err != nil {
			log.Error("Error while writing staging configs to disk: ", err.Error())
			os.Exit(1)
		}
		log.Info("Modifications registered")
	},
}

func init() {
	ConfigModifyCmd.Flags().StringSliceVar(&subscriptions, "subscriptions", []string{}, "Subscriptions - public, beta")
	ConfigModifyCmd.Flags().StringVar(&updatePolicy, "update-policy", "", "Update policy - major, minor, patch, frozen")
	ConfigModifyCmd.Flags().StringVar(&currentVersion, "current-version", "", "version to use")
	ConfigModifyCmd.Flags().StringVar(&storage, "storage", "", "storage location to use")
	ConfigModifyCmd.Flags().StringVar(&runtime, "runtime", "", "runtime to use")
	ConfigModifyCmd.Flags().BoolVar(&forceRuntime, "forced-runtime", false, "forcefully set runtime to use")
}
