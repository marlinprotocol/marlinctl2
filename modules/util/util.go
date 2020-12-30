package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/schollz/progressbar/v3"
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

func IsSystemdAvailable() bool {
	var isSystemdAvailable bool = false
	if _, err := os.Stat("/run/systemd/system"); err == nil {
		isSystemdAvailable = true
	}
	return isSystemdAvailable
}

func IsSupervisorAvailable() bool {
	var isSupervisorAvailable bool = false
	if _, err := os.Stat("/bin/supervisord"); err == nil {
		isSupervisorAvailable = true
	}
	return isSupervisorAvailable
}

func IsSupervisorInRunningState() bool {
	var isSupervisorInRunningState bool = false
	if _, err := os.Stat("/run/supervisor.sock"); err == nil {
		isSupervisorInRunningState = true
	}
	return isSupervisorInRunningState
}

func GetRuntimes() map[string]bool {
	availableRuntimes := []string{"linux-amd64.supervisor", "linux-amd64.systemd"}

	systemPlatform := runtime.GOOS + "-" + runtime.GOARCH

	var isSystemdAvailable bool = IsSystemdAvailable()
	var isSupervisorAvailable bool = IsSupervisorAvailable()

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

func DownloadFile(filepath string, url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading File",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
	return nil
}

func VerifyChecksum(filepath string, md5hash string) error {
	var calculatedMD5 string

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	hashInBytes := hash.Sum(nil)[:16]

	calculatedMD5 = hex.EncodeToString(hashInBytes)

	if calculatedMD5 != md5hash {
		return errors.New("MD5 mismatch. Got " + calculatedMD5 + " while expecting " + md5hash)
	}
	return nil
}

func TrimSpacesEveryLine(s string) string {
	s = strings.Trim(s, " \t\n")
	sArray := strings.Split(s, "\n")
	retString := ""
	ls := len(sArray)

	for i := 0; i < ls; i++ {
		retString = retString + strings.Trim(sArray[i], " \t")
		if i != ls-1 {
			retString = retString + "\n"
		}
	}
	return retString
}

func PrettyPrintKVStruct(s interface{}) {
	v := reflect.ValueOf(s)

	t := GetTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key", "Value"})

	for i := 0; i < v.NumField(); i++ {
		t.AppendRow(table.Row{v.Type().Field(i).Name, v.Field(i).Interface()})
	}
	t.Render()
}

func PrettyPrintKVMap(s map[string]interface{}) {
	t := GetTable()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key", "Value"})

	for k, v := range s {
		t.AppendRow(table.Row{k, v})
	}
	t.Render()
}

func GetTable() table.Writer {
	t := table.NewWriter()
	t.SetStyle(table.Style{Box: table.BoxStyle{
		BottomLeft:       " ",
		BottomRight:      " ",
		BottomSeparator:  " ",
		Left:             " ",
		LeftSeparator:    " ",
		MiddleHorizontal: " ",
		MiddleSeparator:  " ",
		MiddleVertical:   " ",
		PaddingLeft:      " ",
		PaddingRight:     " ",
		PageSeparator:    "\n",
		Right:            " ",
		RightSeparator:   " ",
		TopLeft:          " ",
		TopRight:         " ",
		TopSeparator:     " ",
		UnfinishedRow:    " ",
	}, Color: table.ColorOptions{
		Header: text.Colors{text.FgBlue},
	}})
	return t
}
