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
	"errors"

	ethKeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/marlinprotocol/ctl2/modules/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AppCmd represents the registry command
var KeystoreCmd = &cobra.Command{
	Use:   "keystore",
	Short: "Create or destroy keystore",
	Long:  `Create or destroy Keystore`,
}

func init() {
	KeystoreCmd.AddCommand(createCmd)
	KeystoreCmd.AddCommand(destroyCmd)
}

// returns keystore and keystorePass file path if exists at default location, else return error
func GetKeystoreDetails() (string, string, error) {
	home, err := util.GetUser()
	if err != nil {
		return "", "", err
	}
	kstore := ethKeystore.NewKeyStore(home.HomeDir+"/.marlin/ctl/keys/keystore", ethKeystore.StandardScryptN, ethKeystore.StandardScryptP)
	if len(kstore.Accounts()) == 0 {
		return "", "", errors.New("no existing keystore found")
	}
	return kstore.Accounts()[0].URL.Path, kstore.Accounts()[0].URL.Path + "-pass", nil
}

func KeystoreCheck(cmd *cobra.Command, args []string) error {
	if !cmd.Flags().Changed("keystore-path") {
		_, _, err := GetKeystoreDetails()
		if err != nil {
			log.Error(err)
			log.Error("Please either create a new keystore or provide existing")
			return err
		}
	}
	return nil
}
