package gateway_iris

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
	InstanceId   string
	RunnerData   linux_amd64_supervisor_runner01_runnerdata
	SkipChecksum bool
}

const (
	gatewayName               = "gateway_iris_linux-amd64"
	bridgeName                = "bridge_iris_linux-amd64"
	gatewayProgramName        = "gatewayiris"
	bridgeProgramName         = "bridgeiris"
	defaultUser               = "root"
	supervisorConfFiles       = "/etc/supervisor/conf.d"
	gatewaySupervisorConfFile = "gatewayiris"
	bridgeSupervisorConfFile  = "bridgeiris"
	logRootDir                = "/var/log/supervisor"
	oldLogRootDir             = "/var/log/old_logs"
	projectName               = "gateway_iris"
)

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

	var gatewayLocation = dirPath + "/" + gatewayName
	var bridgeLocation = dirPath + "/" + bridgeName

	if _, err := os.Stat(gatewayLocation); os.IsNotExist(err) {
		log.Info("Fetching gateway from upstream for version ", r.Version)
		util.DownloadFile(gatewayLocation, r.RunnerData.Gateway)
	}
	if !r.SkipChecksum {
		err := util.VerifyChecksum(gatewayLocation, r.RunnerData.GatewayChecksum)
		if err != nil {
			return errors.New("Error while verifying gateway checksum: " + err.Error())
		} else {
			log.Debug("Successully verified gateway's integrity")
		}
	}
	if _, err := os.Stat(bridgeLocation); os.IsNotExist(err) {
		log.Info("Fetching bridge from upstream for version ", r.Version)
		util.DownloadFile(bridgeLocation, r.RunnerData.Bridge)
	}
	if !r.SkipChecksum {
		err := util.VerifyChecksum(bridgeLocation, r.RunnerData.BridgeChecksum)
		if err != nil {
			return errors.New("Error while verifying bridge checksum: " + err.Error())
		} else {
			log.Debug("Successully verified bridge's integrity")
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
	var keyfileLocation = r.Storage + "/common/keyfile.json"
	if _, err := os.Stat(keyfileLocation); os.IsNotExist(err) {
		log.Debug("Creating a new keyfile since none found at " + keyfileLocation)
		var gatewayLocation = r.Storage + "/" + r.Version + "/" + gatewayName
		keyFileGenCommand := exec.Command(gatewayLocation, "keyfile", "--chain=irisnet", "--generate", "--filelocation="+keyfileLocation)
		_, err := keyFileGenCommand.Output()
		if err != nil {
			return errors.New("Keyfile generation error: " + err.Error())
		}
		log.Debug("New Keyfile generated.")
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
	util.PrettyPrintKVStruct(keyFileData)

	return nil
}

func (r *linux_amd64_supervisor_runner01) Create(runtimeArgs map[string]string) error {
	if _, err := os.Stat(GetResourceFileLocation(r.Storage, r.InstanceId)); err == nil {
		return errors.New("Resource file already exisits, cannot create a new instance: " + GetResourceFileLocation(r.Storage, r.InstanceId))
	}

	substitutions := resource{
		"linux-amd64.supervisor.runner01", r.Version, time.Now().Format(time.RFC822Z),
		gatewayProgramName + r.InstanceId, defaultUser, "/", r.Storage + "/" + r.Version + "/" + gatewayName, r.Storage + "/common/keyfile.json", "21900", "127.0.0.1", "21901",
		bridgeProgramName + r.InstanceId, defaultUser, "/", r.Storage + "/" + r.Version + "/" + bridgeName, "127.0.0.1:8002",
	}

	for k, v := range runtimeArgs {
		if k != "GatewayProgram" && k != "BridgeProgram" &&
			reflect.ValueOf(&substitutions).Elem().FieldByName(k).CanSet() {
			reflect.ValueOf(&substitutions).Elem().FieldByName(k).SetString(v)
		}
	}

	log.Info("Running configuration")
	util.PrettyPrintKVStruct(substitutions)

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
	gFile, err := os.Create(supervisorConfFiles + "/" + gatewaySupervisorConfFile + r.InstanceId + ".conf")
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
	bFile, err := os.Create(supervisorConfFiles + "/" + bridgeSupervisorConfFile + r.InstanceId + ".conf")
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
	log.Debug("Trigerred bridge run")

	_, err = exec.Command("supervisorctl", "start", substitutions.GatewayProgram).Output()
	if err != nil {
		return errors.New("Error while starting bridge: " + err.Error())
	}
	log.Debug("Trigerred gateway run")

	log.Info("Waiting 10 seconds to poll for status")
	time.Sleep(10 * time.Second)

	status, err := exec.Command("supervisorctl", "status").Output()
	if err != nil {
		log.Warning("Error while reading supervisor status: " + err.Error())
	} else {
		var supervisorStatus = make(map[string]interface{})

		statusLines := strings.Split(string(status), "\n")
		var anyStatusLine = false
		for _, v := range statusLines {
			if match, err := regexp.MatchString(gatewayProgramName+r.InstanceId+"|"+bridgeProgramName+r.InstanceId, v); err == nil && match {
				vSplit := strings.Split(v, " ")
				supervisorStatus[vSplit[0]] = strings.Trim(strings.Join(vSplit[1:], " "), " ")
				anyStatusLine = true
			}
		}
		if !anyStatusLine {
			log.Info("No proceses seem to be running")
		} else {
			log.Info("Process status")
			util.PrettyPrintKVMap(supervisorStatus)
		}
	}
	r.writeResourceToFile(substitutions, GetResourceFileLocation(r.Storage, r.InstanceId))

	return nil
}

func (r *linux_amd64_supervisor_runner01) Destroy() error {
	available, _, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exists. Can't destroy")
	}

	_, err = exec.Command("supervisorctl", "stop", gatewaySupervisorConfFile+r.InstanceId).Output()
	if err != nil {
		return errors.New("Error while stopping gateway: " + err.Error())
	}
	log.Debug("Trigerred gateway stop")

	_, err = exec.Command("supervisorctl", "stop", bridgeSupervisorConfFile+r.InstanceId).Output()
	if err != nil {
		return errors.New("Error while stopping bridge: " + err.Error())
	}
	log.Debug("Trigerred bridge stop")

	log.Info("Waiting 5 seconds for SIGTERM to take affect")
	time.Sleep(5 * time.Second)

	return nil
}

func (r *linux_amd64_supervisor_runner01) PostRun() error {
	var gatewayConfig = supervisorConfFiles + "/" + gatewaySupervisorConfFile + r.InstanceId + ".conf"
	var bridgeConfig = supervisorConfFiles + "/" + bridgeSupervisorConfFile + r.InstanceId + ".conf"

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

	err := util.CreateDirPathIfNotExists(oldLogRootDir)
	if err != nil {
		return err
	}
	err = filepath.Walk(logRootDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(gatewayProgramName+r.InstanceId+".*|"+bridgeProgramName+r.InstanceId+".*", f.Name())
			if err == nil && r {
				err2 := os.Rename(logRootDir+"/"+f.Name(), oldLogRootDir+"/previous_run_"+f.Name())
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

	err = os.Remove(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return errors.New("Error while removing resource file: " + err.Error())
	}

	log.Info("All relevant processes stopped, resources deleted, supervisor configs removed, logs marked as old")
	return nil
}

func (r *linux_amd64_supervisor_runner01) Status() error {
	available, resData, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exist. Can't return status.")
	}

	var projectConfig types.Project
	err = viper.UnmarshalKey(projectName, &projectConfig)
	if err != nil {
		return err
	}
	log.Info("Project configuration")
	util.PrettyPrintKVStruct(projectConfig)

	var keyfileLocation = r.Storage + "/common/keyfile.json"
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
	util.PrettyPrintKVStruct(keyFileData)

	log.Info("Resource information")
	util.PrettyPrintKVStruct(resData)

	status, err := exec.Command("supervisorctl", "status").Output()
	if err != nil {
		return errors.New("Error while reading supervisor status: " + err.Error())
	}

	var supervisorStatus = make(map[string]interface{})

	statusLines := strings.Split(string(status), "\n")
	var anyStatusLine = false
	for _, v := range statusLines {
		if match, err := regexp.MatchString(gatewayProgramName+r.InstanceId+"|"+bridgeProgramName+r.InstanceId, v); err == nil && match {
			vSplit := strings.Split(v, " ")
			supervisorStatus[vSplit[0]] = strings.Trim(strings.Join(vSplit[1:], " "), " ")
			anyStatusLine = true
		}
	}
	if !anyStatusLine {
		log.Info("No proceses seem to be running")
	} else {
		log.Info("Process status")
		util.PrettyPrintKVMap(supervisorStatus)
	}

	return nil
}

func (r *linux_amd64_supervisor_runner01) Logs() error {
	available, _, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exists. Can't tail logs")
	}
	// Check for resource
	fileSubscriptions := make(map[string]string)
	var logRootDir = "/var/log/supervisor/"
	err = filepath.Walk(logRootDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			for _, v := range []string{gatewaySupervisorConfFile + r.InstanceId + "-stdout.*",
				gatewaySupervisorConfFile + r.InstanceId + "-stderr.*",
				bridgeSupervisorConfFile + r.InstanceId + "-stdout.*",
				bridgeSupervisorConfFile + r.InstanceId + "-stderr.*"} {
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

type resource struct {
	Runner, Version, StartTime                                                                                                                   string
	GatewayProgram, GatewayUser, GatewayRunDir, GatewayExecutablePath, GatewayKeyfile, GatewayListenPortPeer, GatewayMarlinIp, GatewayMarlinPort string
	BridgeProgram, BridgeUser, BridgeRunDir, BridgeExecutablePath, BridgeBootstrapAddr                                                           string
}

func (r *linux_amd64_supervisor_runner01) fetchResourceInformation(fileLocation string) (bool, resource, error) {
	if _, err := os.Stat(fileLocation); os.IsNotExist(err) {
		return false, resource{}, err
	}

	file, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return false, resource{}, err
	}

	var resData = resource{}
	err = json.Unmarshal([]byte(file), &resData)

	return true, resData, err
}

func (r *linux_amd64_supervisor_runner01) writeResourceToFile(resData resource, fileLocation string) error {
	lSplice := strings.Split(fileLocation, "/")
	var dirPath string

	for i := 0; i < len(lSplice)-1; i++ {
		dirPath = dirPath + "/" + lSplice[i]
	}
	err := util.CreateDirPathIfNotExists(dirPath)
	if err != nil {
		log.Error("Error while creating directory ", dirPath, " ", err.Error())
	}

	fileData, err := json.MarshalIndent(resData, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileLocation, fileData, 0644)
}
