package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/schollz/progressbar/v3"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/yaml.v2"
)

func GetUser() (*user.User, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	if os.Geteuid() == 0 {
		// Root, try to retrieve SUDO_USER if exists
		if u := os.Getenv("SUDO_USER"); u != "" {
			usr, err = user.Lookup(u)
			if err != nil {
				return nil, err
			}
		}
	}

	return usr, nil
}

func RemoveConfigEntry(key string) error {
	// viper yaml parser is bad -> it makes nil entries do weird things like
	// map[interface{}]interface{}. Need a better yaml parser.
	// Using Beats authors' workaround given here: https://github.com/go-yaml/yaml/issues/139
	// Make sure whenever you call this function, viper is in sync with disk
	// NVM - this yaml issue is also solved automatically with marshalling
	// IT IS WEIRD. ISSUES MAY OCCUR AT A LATER DAY

	configMap := viper.AllSettings()
	delete(configMap, key)

	encodedConfig, err := yaml.Marshal(configMap)
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

func IsCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
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
	if _, err := os.Stat("/usr/bin/supervisord"); err == nil {
		isSupervisorAvailable = true
	}
	if !IsCommandAvailable("supervisorctl") {
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
		currentUser, err := GetUser()
		if err != nil {
			return err
		}
		uid, err := strconv.Atoi(currentUser.Uid)
		if err != nil {
			return err
		}
		gid, err := strconv.Atoi(currentUser.Gid)
		if err != nil {
			return err
		}
		createErr := os.MkdirAll(dirPath, 0777)
		if createErr != nil {
			return createErr
		}
		return os.Chown(dirPath, uid, gid)
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

func ChownRmarlinctlDir() error {
	currentUser, err := GetUser()
	if err != nil {
		return err
	}
	uid, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(currentUser.Gid)
	if err != nil {
		return err
	}
	return filepath.Walk(currentUser.HomeDir+"/.marlin", func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chown(name, uid, gid)
		}
		return err
	})
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
		return errors.New("MD5 mismatch. Got " + calculatedMD5 + " while expecting " + md5hash + " @ " + filepath)
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

func IsValidUpdatePolicy(updatePolicy string) bool {
	var updatePolicies = map[string]bool{"major": true, "minor": true, "patch": true, "frozen": true}
	if found, ok := updatePolicies[updatePolicy]; found && ok {
		return true
	}
	return false
}

func IsValidSubscription(subscription string) bool {
	var subscriptions = map[string]bool{"public": true, "beta": true, "alpha": true, "dev": true}
	if found, ok := subscriptions[subscription]; found && ok {
		return true
	}
	return false
}

func DecodeVersionString(verString string) (int, int, int, string, int, error) {
	dashSplit := strings.Split(verString, "-")
	if len(dashSplit) < 1 {
		return 0, 0, 0, "", 0, errors.New("Invalid version string")
	} else if len(dashSplit) < 2 {
		var majVer, minVer, patchVer, build int
		var channel string
		dotSplit := strings.Split(dashSplit[0], ".")
		if len(dotSplit) != 3 {
			return 0, 0, 0, "", 0, errors.New("Invalid version string")
		}
		majVer, err1 := strconv.Atoi(dotSplit[0])
		minVer, err2 := strconv.Atoi(dotSplit[1])
		patchVer, err3 := strconv.Atoi(dotSplit[2])
		build = 0
		channel = "public"
		if err1 != nil || err2 != nil || err3 != nil {
			return 0, 0, 0, "", 0, errors.New("Invalid version string")
		}
		return majVer, minVer, patchVer, channel, build, nil
	} else if len(dashSplit) == 2 {
		var majVer, minVer, patchVer, build int
		var channel string
		dotSplit := strings.Split(dashSplit[0], ".")
		if len(dotSplit) != 3 {
			return 0, 0, 0, "", 0, errors.New("Invalid version string")
		}
		majVer, err1 := strconv.Atoi(dotSplit[0])
		minVer, err2 := strconv.Atoi(dotSplit[1])
		patchVer, err3 := strconv.Atoi(dotSplit[2])
		dotSplit2 := strings.Split(dashSplit[1], ".")
		if len(dotSplit2) != 2 {
			return 0, 0, 0, "", 0, errors.New("Invalid version string")
		}
		channel = dotSplit2[0]
		build, err4 := strconv.Atoi(dotSplit2[1])
		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return 0, 0, 0, "", 0, errors.New("Invalid version string")
		}
		return majVer, minVer, patchVer, channel, build, nil
	}
	return 0, 0, 0, "", 0, errors.New("Invalid version string")
}

func CanUseVersion(maj1 int, min1 int, patch1 int, sub1 string, build1 int,
	maj2 int, min2 int, patch2 int, sub2 string, build2 int,
	updatePolicy string) bool {
	switch updatePolicy {
	case "frozen":
		return (maj1 == maj2 && min1 == min2 && patch1 == patch2 && sub1 == sub2 && build1 == build2)
	case "patch":
		return (maj1 == maj2 && min1 == min2)
	case "minor":
		return (maj1 == maj2)
	case "major":
		return true
	}
	return false
}

func IsHigherVersion(maj1 int, min1 int, patch1 int,
	maj2 int, min2 int, patch2 int, sub2 string) bool {
	if maj1 > maj2 {
		return true
	} else if maj1 == maj2 {
		if min1 > min2 {
			return true
		} else if min1 == min2 {
			if patch1 > patch2 {
				return true
			} else if patch1 == patch2 && sub2 != "public" {
				return true
			}
		}
	}
	return false
}

func ExpandTilde(path string) string {
	usr, _ := GetUser()
	dir := usr.HomeDir
	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		path = filepath.Join(dir, path[2:])
	}
	return path
}

func PrintPrettyDiff(message string) {
	lines := strings.Split(message, "\n")
	for _, l := range lines {
		if len(l) < 1 {
			continue
		} else if l[0] == []byte("-")[0] {
			fmt.Println(text.FgRed.Sprintf(l))
		} else if l[0] == []byte("+")[0] {
			fmt.Println(text.FgGreen.Sprintf(l))
		} else {
			fmt.Println(l)
		}
	}
}

func GetFileSeekOffsetLastNLines(fname string, lines int) int64 {
	file, err := os.Open(fname)
	if err != nil {
		return 0
	}
	defer file.Close()

	buf := make([]byte, 1000*lines)
	stat, err := os.Stat(fname)
	start := stat.Size() - int64(1000*lines)
	if start < 0 {
		start = 0
	}

	_, err = file.ReadAt(buf, start)

	linesEncountered := 0
	offset := stat.Size() - start - 1

	for ; offset >= 0 && linesEncountered <= lines; offset-- {
		if buf[offset] == '\n' {
			linesEncountered += 1
		}
	}

	offset = offset + start + 1
	if offset < 0 {
		return 0
	}
	return offset
}

func ReadInputPasswordLine() (string, error) {
	textBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(textBytes), "\n"), nil
}

func ReadStringFromFile(filePath string) (string, error) {
	passBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(passBytes), "\n"), nil
}
