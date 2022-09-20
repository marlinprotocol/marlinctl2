package proxy

import (
	"github.com/marlinprotocol/ctl2/cmd/proxy/tcptovsock"
	"github.com/marlinprotocol/ctl2/cmd/proxy/vsocktotcp"
	"github.com/spf13/cobra"
)

var ProxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Proxy servers",
	Long:  `Run and control proxy servers`,
}

func init() {
	ProxyCmd.AddCommand(tcptovsock.TcpToVsockCmd)
	ProxyCmd.AddCommand(vsocktotcp.VsockToTcpCmd)
}


