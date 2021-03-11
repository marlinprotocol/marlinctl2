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
	"strings"

	ethKeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/marlinprotocol/ctl2/modules/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var keystorePassPath string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create keystore",
	Long:  `Create a new Keystore`,
	Run: func(cmd *cobra.Command, args []string) {
		// read the password file
		passBytes, err := ioutil.ReadFile(keystorePassPath)
		if err != nil {
			log.Error("cannot read keystore password file at path ", keystorePassPath)
			return
		}

		home, err := util.GetUser()
		if err != nil {
			log.Error("error while getting user ", err)
			return
		}

		kstore := ethKeystore.NewKeyStore(home.HomeDir+"/.marlin/ctl/keys/keystore", ethKeystore.StandardScryptN, ethKeystore.StandardScryptP)
		if len(kstore.Accounts()) != 0 {
			log.Error("Keystore already exists.")
			return
		}
		passphrase := string(passBytes)
		passphrase = strings.TrimSuffix(passphrase, "\n")
		_, err = kstore.NewAccount(passphrase)
		if err != nil {
			log.Error("error while creating new account", err)
			return
		}
		log.Info("created new keysore with address ", kstore.Accounts()[0].Address)

		if err := ioutil.WriteFile(kstore.Accounts()[0].URL.Path+"-pass", passBytes, 0644); err != nil {
			log.Error("error in writing password file ", err)
			if err := kstore.Delete(kstore.Accounts()[0], string(passBytes)); err != nil {
				log.Error("error while deleting previous keystore", err)
			} else {
				log.Info("Deleted keystore. Please create again")
			}
		}
	},
}

func init() {
	createCmd.Flags().StringVarP(&keystorePassPath, "pass-path", "p", "", "path to the passphrase file")
	createCmd.MarkFlagRequired("pass-path")
}
