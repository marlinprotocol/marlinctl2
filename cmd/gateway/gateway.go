package gateway

import (
	"github.com/marlinprotocol/ctl2/cmd/gateway/iris"
	"github.com/spf13/cobra"
)

var GatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Run gateways of various blockchaing",
	Long:  `Allows controlling gateways (+bridges) for multiple blockchains`,
}

func init() {
	GatewayCmd.AddCommand(iris.IrisCmd)
}
