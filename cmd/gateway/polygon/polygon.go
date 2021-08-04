package polygon

import (
	"github.com/marlinprotocol/ctl2/cmd/gateway/polygon/bor"
	"github.com/spf13/cobra"
)

var PolygonCmd = &cobra.Command{
	Use:   "polygon",
	Short: "Polygon Gateway",
	Long:  `Polygon Gateway`,
}

func init() {
	PolygonCmd.AddCommand(bor.BorCmd)
}
