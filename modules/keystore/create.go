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
	"io/ioutil"
	"strings"

	ethKeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	log "github.com/sirupsen/logrus"
)

func Create(keystoreDirPath string, keystorePassPath string) error {

	// read the password file
	passBytes, err := ioutil.ReadFile(keystorePassPath)
	if err != nil {
		return errors.New("cannot read keystore password file at path " + keystorePassPath)
	}

	kstore := ethKeystore.NewKeyStore(keystoreDirPath, ethKeystore.StandardScryptN, ethKeystore.StandardScryptP)
	if len(kstore.Accounts()) != 0 {
		return errors.New("Keystore already exists")
	}
	passphrase := string(passBytes)
	passphrase = strings.TrimSuffix(passphrase, "\n")

	_, err = kstore.NewAccount(passphrase)
	if err != nil {
		return errors.New("error while creating new account")
	}

	log.Info("created new keysore with address ", kstore.Accounts()[0].Address)

	if err := ioutil.WriteFile(kstore.Accounts()[0].URL.Path+"-pass", passBytes, 0644); err != nil {
		log.Error("error in writing password file ", err)
		if err := kstore.Delete(kstore.Accounts()[0], string(passBytes)); err != nil {
			log.Error("error while deleting previous keystore", err)
		} else {
			log.Info("Deleted keystore. Please create again")
		}
		return err
	}
	return nil
}
