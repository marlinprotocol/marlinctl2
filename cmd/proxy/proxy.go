package proxy

import (
	"github.com/marlinprotocol/ctl2/cmd/proxy/iptovsock"
	"github.com/marlinprotocol/ctl2/cmd/proxy/vsocktoip"
	"github.com/spf13/cobra"
)

var ProxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Proxy servers",
	Long:  `Run and control proxy servers`,
}

func init() {
	ProxyCmd.AddCommand(iptovsock.IpToVsockCmd)
	ProxyCmd.AddCommand(vsocktoip.VsockToIpCmd)
}


