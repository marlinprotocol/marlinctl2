package actions

import (
	"os"

	"github.com/google/go-cmp/cmp"
	cmn "github.com/marlinprotocol/ctl2/cmd/relay/eth/common"
	"github.com/marlinprotocol/ctl2/modules/util"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppCmd represents the registry command
var ConfigDiffCmd = &cobra.Command{
	Use:     "diff",
	Short:   "Show difference between current project configs and modifications",
	Long:    `Show difference between current project configs and modifications`,
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

		log.Info("Difference:")
		util.PrintPrettyDiff(cmp.Diff(projectConfig, projectConfigMod))
	},
}
