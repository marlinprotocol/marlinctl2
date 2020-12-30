package iris_endnode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/hpcloud/tail"
	"github.com/marlinprotocol/ctl2/modules/util"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type linux_amd64_supervisor_runner01_runnerdata struct {
	Gateway         string
	GatewayChecksum string
	Bridge          string
	BridgeChecksum  string
}

type linux_amd64_supervisor_runner01 struct {
	Version      string
	Storage      string
	RunnerData   linux_amd64_supervisor_runner01_runnerdata
	SkipChecksum bool
}

func (r *linux_amd64_supervisor_runner01) PreRunSanity() error {
	if !util.IsSupervisorAvailable() {
		return errors.New("System does not support supervisor")
	}
	if !util.IsSupervisorInRunningState() {
		return errors.New("System does not have supervisor in running state")
	}
	return nil
}

func (r *linux_amd64_supervisor_runner01) Download() error {
	var dirPath = r.Storage + "/" + r.Version
	err := util.CreateDirPathIfNotExists(dirPath)
	if err != nil {
		return err
	}

	var gatewayLocation = dirPath + "/iris_gateway_linux-amd64"
	var bridgeLocation = dirPath + "/iris_bridge_linux-amd64"

	if _, err := os.Stat(gatewayLocation); os.IsNotExist(err) {
		log.Info("Fetching iris gateway from upstream")
		util.DownloadFile(gatewayLocation, r.RunnerData.Gateway)
	}
	if !r.SkipChecksum {
		err := util.VerifyChecksum(gatewayLocation, r.RunnerData.GatewayChecksum)
		if err != nil {
			return errors.New("Error while verifying gateway checksum: " + err.Error())
		} else {
			log.Info("Successully verified gateway's integrity")
		}
	}
	if _, err := os.Stat(bridgeLocation); os.IsNotExist(err) {
		log.Info("Fetching iris bridge from upstream")
		util.DownloadFile(bridgeLocation, r.RunnerData.Bridge)
	}
	if !r.SkipChecksum {
		err := util.VerifyChecksum(bridgeLocation, r.RunnerData.BridgeChecksum)
		if err != nil {
			return errors.New("Error while verifying gateway checksum: " + err.Error())
		} else {
			log.Info("Successully verified bridge's integrity")
		}
	}

	err = os.Chmod(gatewayLocation, 0755)
	if err != nil {
		return err
	}
	err = os.Chmod(bridgeLocation, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (r *linux_amd64_supervisor_runner01) Prepare() error {
	err := r.Download()
	if err != nil {
		return err
	}

	var keyfileDir = r.Storage + "/common"
	err = util.CreateDirPathIfNotExists(keyfileDir)
	if err != nil {
		return err
	}
	var keyfileLocation = r.Storage + "/common/iris_keyfile.json"
	if _, err := os.Stat(keyfileLocation); os.IsNotExist(err) {
		log.Info("Creating a new keyfile since none found at " + keyfileLocation)
		var gatewayLocation = r.Storage + "/" + r.Version + "/iris_gateway_linux-amd64"
		keyFileGenCommand := exec.Command(gatewayLocation, "keyfile", "--chain=irisnet", "--generate", "--filelocation="+keyfileLocation)
		_, err := keyFileGenCommand.Output()
		if err != nil {
			return errors.New("Keyfile generation error: " + err.Error())
		}
		log.Info("New Keyfile generated.")
	}
	keyfile, err := os.Open(keyfileLocation)
	if err != nil {
		return err
	}
	defer keyfile.Close()
	byteValue, _ := ioutil.ReadAll(keyfile)
	var keyFileData = struct {
		NodeId string `json:"IdString"`
	}{}
	json.Unmarshal(byteValue, &keyFileData)

	log.Info("Keyfile information")
	util.PrettyPrintKV(keyFileData)

	return nil
}

func (r *linux_amd64_supervisor_runner01) Create(runtimeArgs map[string]string) error {
	substitutions := struct {
		GatewayProgram, GatewayUser, GatewayRunDir, GatewayExecutablePath, GatewayKeyfile, GatewayListenPortPeer, GatewayMarlinIp, GatewayMarlinPort string
		BridgeProgram, BridgeUser, BridgeRunDir, BridgeExecutablePath, BridgeBootstrapAddr                                                           string
	}{
		"irisgateway", "root", "/", r.Storage + "/" + r.Version + "/iris_gateway_linux-amd64", r.Storage + "/common/iris_keyfile.json", "21900", "127.0.0.1", "21901",
		"irisbridge", "root", "/", r.Storage + "/" + r.Version + "/iris_bridge_linux-amd64", "127.0.0.1:8002",
	}

	for k, v := range runtimeArgs {
		if k != "GetewayProgram" && k != "BridgeProgram" &&
			reflect.ValueOf(&substitutions).Elem().FieldByName(k).CanSet() {
			reflect.ValueOf(&substitutions).Elem().FieldByName(k).SetString(v)
		}
	}

	log.Info("Running configuration")
	util.PrettyPrintKV(substitutions)

	gt := template.Must(template.New("gateway-template").Parse(util.TrimSpacesEveryLine(`
		[program:{{.GatewayProgram}}]
		process_name={{.GatewayProgram}}
		user={{.GatewayUser}}
		directory={{.GatewayRunDir}}
		command={{.GatewayExecutablePath}} dataconnect --keyfile {{.GatewayKeyfile}} --listenportpeer {{.GatewayListenPortPeer}} --marlinip {{.GatewayMarlinIp}} --marlinport {{.GatewayMarlinPort}}
		priority=100
		numprocs=1
		numprocs_start=1
		autostart=true
		autorestart=true
	`)))
	gFile, err := os.Create("/etc/supervisor/conf.d/irisgateway.conf")
	if err != nil {
		return err
	}
	defer gFile.Close()
	if err := gt.Execute(gFile, substitutions); err != nil {
		panic(err)
	}

	bt := template.Must(template.New("bridge-template").Parse(util.TrimSpacesEveryLine(`
		[program:{{.BridgeProgram}}]
		process_name={{.BridgeProgram}}
		user={{.BridgeUser}}
		directory={{.BridgeRunDir}}
		command={{.BridgeExecutablePath}} -b"{{.BridgeBootstrapAddr}}"
		priority=100
		numprocs=1
		numprocs_start=1
		autostart=true
		autorestart=true
	`)))
	bFile, err := os.Create("/etc/supervisor/conf.d/irisbridge.conf")
	if err != nil {
		return err
	}
	defer bFile.Close()
	if err := bt.Execute(bFile, substitutions); err != nil {
		panic(err)
	}

	_, err = exec.Command("supervisorctl", "reread").Output()
	if err != nil {
		return errors.New("Error while supervisorctl reread: " + err.Error())
	}

	_, err = exec.Command("supervisorctl", "update").Output()
	if err != nil {
		return errors.New("Error while supervisorctl update: " + err.Error())
	}

	_, err = exec.Command("supervisorctl", "start", substitutions.BridgeProgram).Output()
	if err != nil {
		return errors.New("Error while starting bridge: " + err.Error())
	}
	log.Info("Trigerred bridge run")

	_, err = exec.Command("supervisorctl", "start", substitutions.GatewayProgram).Output()
	if err != nil {
		return errors.New("Error while starting bridge: " + err.Error())
	}
	log.Info("Trigerred gateway run")

	log.Info("Waiting 10 seconds to poll for status")
	time.Sleep(10 * time.Second)

	status, err := exec.Command("supervisorctl", "status").Output()
	if err != nil {
		return errors.New("Error while reading supervisor status: " + err.Error())
	}
	statusLines := strings.Split(string(status), "\n")
	for _, v := range statusLines {
		if match, err := regexp.MatchString(substitutions.GatewayProgram+"|"+substitutions.BridgeProgram, v); err == nil && match {
			log.Info("{SUPERVISOR STATUS} " + v)
		}
	}

	return nil
}

func (r *linux_amd64_supervisor_runner01) Destroy() error {
	_, err := exec.Command("supervisorctl", "stop", "irisgateway").Output()
	if err != nil {
		return errors.New("Error while stopping gateway: " + err.Error())
	}
	log.Info("Trigerred gateway stop")

	_, err = exec.Command("supervisorctl", "stop", "irisbridge").Output()
	if err != nil {
		return errors.New("Error while stopping bridge: " + err.Error())
	}
	log.Info("Trigerred bridge stop")

	log.Info("Waiting 5 seconds for SIGTERM to take affect")
	time.Sleep(5 * time.Second)

	return nil
}

func (r *linux_amd64_supervisor_runner01) PostRun() error {
	var gatewayConfig = "/etc/supervisor/conf.d/irisgateway.conf"
	var bridgeConfig = "/etc/supervisor/conf.d/irisbridge.conf"

	if _, err := os.Stat(gatewayConfig); !os.IsNotExist(err) {
		if err := os.Remove(gatewayConfig); err != nil {
			return err
		}
	}
	if _, err := os.Stat(bridgeConfig); !os.IsNotExist(err) {
		if err := os.Remove(bridgeConfig); err != nil {
			return err
		}
	}

	var logRootDir = "/var/log/supervisor/"
	var oldLogsRootDir = "/var/log/old_logs/"
	err := util.CreateDirPathIfNotExists(oldLogsRootDir)
	if err != nil {
		return err
	}
	err = filepath.Walk(logRootDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString("irisgateway.*|irisbridge.*", f.Name())
			if err == nil && r {
				err2 := os.Rename(logRootDir+f.Name(), oldLogsRootDir+"previous_run_"+f.Name())
				if err2 != nil {
					return err2
				}
			} else if err != nil {
				return nil
			}
		}
		return nil
	})

	if err != nil {
		return errors.New("Error while marking logs as old: " + err.Error())
	}

	_, err = exec.Command("supervisorctl", "reread").Output()
	if err != nil {
		return errors.New("Error while supervisorctl reread: " + err.Error())
	}

	_, err = exec.Command("supervisorctl", "update").Output()
	if err != nil {
		return errors.New("Error while supervisorctl update: " + err.Error())
	}

	log.Info("All iris processes stopped, supervisor configs removed, logs marked as old")
	return nil
}

func (r *linux_amd64_supervisor_runner01) Status() error {
	var projectConfig types.Project
	err := viper.UnmarshalKey("iris_endnode", &projectConfig)
	if err != nil {
		return err
	}
	log.Info("Project configuration")
	util.PrettyPrintKV(projectConfig)

	var keyfileLocation = r.Storage + "/common/iris_keyfile.json"
	keyfile, err := os.Open(keyfileLocation)
	if err != nil {
		return err
	}
	defer keyfile.Close()
	byteValue, _ := ioutil.ReadAll(keyfile)
	var keyFileData = struct {
		NodeId string `json:"IdString"`
	}{}
	json.Unmarshal(byteValue, &keyFileData)

	log.Info("Keyfile information")
	util.PrettyPrintKV(keyFileData)

	status, err := exec.Command("supervisorctl", "status").Output()
	if err != nil {
		return errors.New("Error while reading supervisor status: " + err.Error())
	}
	statusLines := strings.Split(string(status), "\n")
	var anyStatusLine = false
	for _, v := range statusLines {
		if match, err := regexp.MatchString("irisgateway|irisbridge", v); err == nil && match {
			log.Info("{SUPERVISOR STATUS} " + v)
			anyStatusLine = true
		}
	}
	if !anyStatusLine {
		log.Info("No proceses seem to be running")
	}

	return nil
}

func (r *linux_amd64_supervisor_runner01) Logs() error {
	fileSubscriptions := make(map[string]string)
	var logRootDir = "/var/log/supervisor/"
	err := filepath.Walk(logRootDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			for _, v := range []string{"irisgateway-stdout.*", "irisgateway-stderr.*", "irisbridge-stdout.*", "irisbridge-stdout.*"} {
				r, err := regexp.MatchString(v, f.Name())
				if err == nil && r {
					fileSubscriptions[v[:len(v)-2]] = logRootDir + f.Name()
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil
	}

	var wg sync.WaitGroup
	for k, v := range fileSubscriptions {
		wg.Add(1)
		go func(filename string, filelocation string) {
			t, err := tail.TailFile(filelocation, tail.Config{Follow: true})
			if err != nil {
				fmt.Println(err)
			}
			for line := range t.Lines {
				log.Info(fmt.Sprintf("[%20s] ", filename) + line.Text)
			}
			wg.Done()
		}(k, v)
	}
	wg.Wait()
	return nil
}
