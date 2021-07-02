package matic

import (
	"github.com/marlinprotocol/ctl2/cmd/gateway/matic/bor"
	"github.com/spf13/cobra"
)

var MaticCmd = &cobra.Command{
	Use:   "matic",
	Short: "Matic Gateway",
	Long:  `Matic Gateway`,
}

func init() {
	MaticCmd.AddCommand(bor.BorCmd)
}
