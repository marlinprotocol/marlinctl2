package actions

import (
	"os"

	"github.com/marlinprotocol/ctl2/modules/util"

	cmn "github.com/marlinprotocol/ctl2/cmd/relay/eth/common"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppCmd represents the registry command
var ConfigApplyCmd = &cobra.Command{
	Use:     "apply",
	Short:   "Apply project config modifications",
	Long:    `Apply project config modifications`,
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
				log.Error("Error while unmarshalling modifications: ", err.Error())
				os.Exit(1)
			}
		} else {
			log.Info("No existing modifications found.")
			os.Exit(1)
		}

		viper.Set(cmn.ProjectID, projectConfigMod)
		err = viper.WriteConfig()
		if err != nil {
			log.Error("Error while writing configs to disk: ", err.Error())
			os.Exit(1)
		}
		err = util.RemoveConfigEntry(modifiedProjectID)
		if err != nil {
			log.Error("Error while removing staging configs from disk: ", err.Error())
			os.Exit(1)
		}
		log.Info("Configs have been updated. These will take effect next time you create an application instance.")
	},
}

func init() {
	ConfigApplyCmd.Flags().BoolVar(&skipRecreate, "skip-recreate", false, "Skip recreating project's running resources while updating configs")
}
