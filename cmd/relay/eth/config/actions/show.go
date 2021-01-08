package actions

import (
	"encoding/json"
	"os"

	cmn "github.com/marlinprotocol/ctl2/cmd/relay/eth/common"
	"github.com/marlinprotocol/ctl2/modules/util"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppCmd represents the registry command
var ConfigShowCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show current project configuration",
	Long:    `Show current project configuration`,
	PreRunE: ConfigTest,
	Run: func(cmd *cobra.Command, args []string) {
		var projectConfig types.Project
		err := viper.UnmarshalKey(cmn.ProjectID, &projectConfig)
		if err != nil {
			log.Error("Error while reading project configs: ", err.Error())
			os.Exit(1)
		}
		s, err := json.MarshalIndent(projectConfig, "", "  ")
		if err != nil {
			log.Error("Error while decoding json: ", err.Error())
			os.Exit(1)
		}
		log.Info("Current config:")
		util.PrintPrettyDiff(string(s))
	},
}
