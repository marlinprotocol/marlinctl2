package keys

import (
	"github.com/marlinprotocol/ctl2/cmd/keys/keystore"
	"github.com/spf13/cobra"
)

var KeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "manages keys",
	Long:  `manages keys`,
}

func init() {
	KeysCmd.AddCommand(keystore.KeystoreCmd)
}
