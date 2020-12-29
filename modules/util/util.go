package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
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

	var isSupervisorAvailable bool = false
	if _, err := os.Stat("/bin/supervisord"); err == nil {
		isSupervisorAvailable = true
	}

	var returnMap = make(map[string]bool)

	for _, runtime := range availableRuntimes {
		rInfo := strings.Split(runtime, ".")

		platform := rInfo[0]
		runner := rInfo[1]

		if platform == systemPlatform {
			switch runner {
			case "supervisor":
				if isSupervisorAvailable {
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
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0777)
	}
	return nil
}

func RemoveDirContents(dirPath string) error {
	files, err := filepath.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		log.Info("Removing ", file)
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func GitPullHead(upstreamUrl string, branch string, dirPath string) error {
	_, err := git.PlainClone(dirPath, false, &git.CloneOptions{
		URL:           upstreamUrl,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		SingleBranch:  true,
		Depth:         1,
	})
	return err
}

func MoveDir(srcDir string, dstDir string) error {
	err := RemoveDirPathIfExists(dstDir)
	if err != nil {
		return nil
	}
	return os.Rename(srcDir, dstDir)
}

func RemoveDirPathIfExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil
	} else {
		return os.RemoveAll(dirPath)
	}
}
