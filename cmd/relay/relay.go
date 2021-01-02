package relay

import (
	"github.com/marlinprotocol/ctl2/cmd/relay/eth"
	"github.com/spf13/cobra"
)

var RelayCmd = &cobra.Command{
	Use:   "relay",
	Short: "Run relays of various blockchaing",
	Long:  `Allows controlling relays (+abci) for multiple blockchains`,
}

func init() {
	RelayCmd.AddCommand(eth.EthCmd)
}
