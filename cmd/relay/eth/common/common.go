package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	projectRunners "github.com/marlinprotocol/ctl2/modules/runner/relay_eth"
	"github.com/marlinprotocol/ctl2/types"
	"github.com/prometheus/common/log"
)

const (
	ProjectID string = types.ProjectID_relay_eth
)

func GetResourceMetaData(projectConfig types.Project, instanceId string) (string, string, error) {
	resFileLocation := projectRunners.GetResourceFileLocation(projectConfig.Storage, instanceId)
	if _, err := os.Stat(resFileLocation); os.IsNotExist(err) {
		return "", "", errors.New("Cannot locate resource: " + resFileLocation)
	}
	file, err := ioutil.ReadFile(resFileLocation)
	if err != nil {
		return "", "", err
	}
	var resourceMetaData = struct {
		Runner  string `json:"Runner"`
		Version string `json:"Version"`
	}{}
	err = json.Unmarshal([]byte(file), &resourceMetaData)
	if err != nil {
		return "", "", err
	}
	log.Debug("Resource metadata: ", resourceMetaData)
	return resourceMetaData.Runner, resourceMetaData.Version, nil
}
