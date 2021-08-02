package gateway

import (
	"github.com/marlinprotocol/ctl2/cmd/gateway/cosmos"
	"github.com/marlinprotocol/ctl2/cmd/gateway/dot"
	"github.com/marlinprotocol/ctl2/cmd/gateway/iris"
	"github.com/marlinprotocol/ctl2/cmd/gateway/near"
	"github.com/marlinprotocol/ctl2/cmd/gateway/polygon"
	"github.com/spf13/cobra"
)

var GatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Run gateways of various blockchaing",
	Long:  `Allows controlling gateways (+bridges) for multiple blockchains`,
}

func init() {
	GatewayCmd.AddCommand(iris.IrisCmd)
	GatewayCmd.AddCommand(cosmos.CosmosCmd)
	GatewayCmd.AddCommand(near.NearCmd)
	GatewayCmd.AddCommand(dot.DotCmd)
	GatewayCmd.AddCommand(polygon.PolygonCmd)
}
