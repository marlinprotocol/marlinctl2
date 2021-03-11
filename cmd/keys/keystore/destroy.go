/*
Copyright © 2020 MARLIN TEAM <info@marlin.pro>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package keystore

import (
	"io/ioutil"
	"os"
	"strings"

	ethKeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/marlinprotocol/ctl2/modules/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Delete keystore",
	Long:  `Delete the existing keystore`,
	Run: func(cmd *cobra.Command, args []string) {

		home, err := util.GetUser()
		if err != nil {
			log.Error("error while getting user ", err)
			return
		}

		kstore := ethKeystore.NewKeyStore(home.HomeDir+"/.marlin/ctl/keys/keystore", ethKeystore.StandardScryptN, ethKeystore.StandardScryptP)
		if len(kstore.Accounts()) == 0 {
			log.Error("No keystore found")
			return
		}
		keystorePassPath := kstore.Accounts()[0].URL.Path + "-pass"
		passBytes, err := ioutil.ReadFile(keystorePassPath)
		if err != nil {
			log.Error("cannot read keystore password file at path ", keystorePassPath)
			return
		}
		passphrase := string(passBytes)
		passphrase = strings.TrimSuffix(passphrase, "\n")
		if err := kstore.Delete(kstore.Accounts()[0], passphrase); err != nil {
			log.Error("error while deleting keystore ", err)
			return
		}
		if err := os.Remove(keystorePassPath); err != nil {
			log.Error("error in deleting password file ", err)
			return
		}
		log.Info("successfully deleted keystore")
	},
}

func init() {

}
