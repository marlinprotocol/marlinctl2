package util

import (
	"bytes"
	"encoding/json"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func RemoveConfigEntry(key string) error {
	configMap := viper.AllSettings()
	delete(configMap, key)
	encodedConfig, err := json.MarshalIndent(configMap, "", " ")
	if err != nil {
		return err
	}
	err = viper.ReadConfig(bytes.NewReader(encodedConfig))
	if err != nil {
		return err
	}
	viper.WriteConfig()
	return nil
}

func GetRuntimes() map[string]bool {
	availableRuntimes := []string{"linux-amd64.supervisor", "linux-amd64.systemd"}

	systemPlatform := runtime.GOOS + "-" + runtime.GOARCH

	var isSystemdAvailable bool = false
	if _, err := os.Stat("/run/systemd/system"); err == nil {
		isSystemdAvailable = true
	}

	var isSupervisordAvailable bool = false
	if _, err := os.Stat("/bin/supervisord"); err == nil {
		isSupervisordAvailable = true
	}

	var returnMap map[string]bool

	for _, runtime := range availableRuntimes {
		rInfo := strings.Split(runtime, ".")

		platform := rInfo[0]
		runner := rInfo[1]

		if platform == systemPlatform {
			switch runner {
			case "supervisor":
				if isSupervisordAvailable {
					returnMap[runtime] = true
				} else {
					returnMap[runtime] = false
				}
			case "systemd":
				if isSystemdAvailable {
					returnMap[runtime] = true
				} else {
					returnMap[runtime] = false
				}
			default:
				returnMap[runtime] = false
			}
		} else {
			returnMap[runtime] = false
		}
	}

	return returnMap
}

func CreateDirPathIfNotExists(dirPath string) error {
	log.Warning("Yet to implement")
	return nil
}

func RemoveDirContents(dirPath string) error {
	log.Warning("Yet to implement")
	return nil
}

func GitPullHead(upstreamUrl string, branch string, dirPath string) error {
	log.Warning("Yet to implement")
	return nil
}

func MoveDir(srcDir string, dstDir string) error {
	log.Warning("Yet to implement")
	return nil
}

func RemoveDirPathIfExists(dirPath string) error {
	log.Warning("Yet to implement")
	return nil
}
