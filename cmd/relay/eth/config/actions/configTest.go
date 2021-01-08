package actions

import (
	cmn "github.com/marlinprotocol/ctl2/cmd/relay/eth/common"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ConfigTest = func(cmd *cobra.Command, args []string) error {
	var marlinConfig types.Project
	err := viper.UnmarshalKey(types.ProjectID_marlinctl, &marlinConfig)
	if err != nil {
		return err
	}
	if !viper.IsSet(cmn.ProjectID) {
		log.Debug("Setting up default config for running relay_eth.")
		updPol, ok1 := marlinConfig.AdditionalInfo["defaultprojectupdatepolicy"]
		defRun, ok2 := marlinConfig.AdditionalInfo["defaultprojectruntime"]
		if ok1 && ok2 {
			SetupConfiguration(false,
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
