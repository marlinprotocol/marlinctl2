package gateway_cosmos

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

type linux_amd64_supervisor_runner02_runnerdata struct {
	Gateway         string
	GatewayChecksum string
	Bridge          string
	BridgeChecksum  string
}

type linux_amd64_supervisor_runner02 struct {
	Version      string
	Storage      string
	InstanceId   string
	RunnerData   linux_amd64_supervisor_runner02_runnerdata
	SkipChecksum bool
}

const (
	runner02gatewayName               = "gateway_cosmos_linux-amd64"
	runner02bridgeName                = "bridge_cosmos_linux-amd64"
	runner02gatewayProgramName        = "gateway_cosmos"
	runner02bridgeProgramName         = "bridge_cosmos"
	runner02defaultUser               = "root"
	runner02supervisorConfFiles       = "/etc/supervisor/conf.d"
	runner02gatewaySupervisorConfFile = "gateway_cosmos"
	runner02bridgeSupervisorConfFile  = "bridge_cosmos"
	runner02logRootDir                = "/var/log/supervisor"
	runner02oldLogRootDir             = "/var/log/old_logs"
	runner02projectName               = "gateway_cosmos"
)

func (r *linux_amd64_supervisor_runner02) PreRunSanity() error {
	if !util.IsSupervisorAvailable() {
		return errors.New("System does not support supervisor")
	}
	if !util.IsSupervisorInRunningState() {
		return errors.New("System does not have supervisor in running state")
	}
	return nil
}

func (r *linux_amd64_supervisor_runner02) Download() error {
	var dirPath = r.Storage + "/" + r.Version
	err := util.CreateDirPathIfNotExists(dirPath)
	if err != nil {
		return err
	}

	var gatewayLocation = dirPath + "/" + runner02gatewayName
	var bridgeLocation = dirPath + "/" + runner02bridgeName

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

func (r *linux_amd64_supervisor_runner02) Prepare() error {
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
		var gatewayLocation = r.Storage + "/" + r.Version + "/" + runner02gatewayName
		keyFileGenCommand := exec.Command(gatewayLocation, "keyfile", "--chain=cosmos", "--generate", "--filelocation="+keyfileLocation)
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

	err = util.ChownRmarlinctlDir()
	if err != nil {
		return err
	}

	return nil
}

func (r *linux_amd64_supervisor_runner02) Create(runtimeArgs map[string]string) error {
	if _, err := os.Stat(GetResourceFileLocation(r.Storage, r.InstanceId)); err == nil {
		return errors.New("Resource file already exisits, cannot create a new instance: " + GetResourceFileLocation(r.Storage, r.InstanceId))
	}

	substitutions := runner02resource{
		"linux-amd64.supervisor.runner02", r.Version, time.Now().Format(time.RFC822Z),
		runner02gatewayProgramName + "_" + r.InstanceId, runner02defaultUser, "/", r.Storage + "/" + r.Version + "/" + runner02gatewayName, r.Storage + "/common/keyfile.json", "22400", "127.0.0.1", "22401",
		runner02bridgeProgramName + "_" + r.InstanceId, runner02defaultUser, "/", r.Storage + "/" + r.Version + "/" + runner02bridgeName, "127.0.0.1:8002",
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
		stdout_logfile=/var/log/supervisor/{{.GatewayProgram}}-stdout.log
		stderr_logfile=/var/log/supervisor/{{.GatewayProgram}}-stderr.log
	`)))
	gFile, err := os.Create(runner02supervisorConfFiles + "/" + runner02gatewaySupervisorConfFile + "_" + r.InstanceId + ".conf")
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
		stdout_logfile=/var/log/supervisor/{{.BridgeProgram}}-stdout.log
		stderr_logfile=/var/log/supervisor/{{.BridgeProgram}}-stderr.log
	`)))
	bFile, err := os.Create(runner02supervisorConfFiles + "/" + runner02bridgeSupervisorConfFile + "_" + r.InstanceId + ".conf")
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
			if match, err := regexp.MatchString(runner02gatewayProgramName+"_"+r.InstanceId+"|"+runner02bridgeProgramName+"_"+r.InstanceId, v); err == nil && match {
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

func (r *linux_amd64_supervisor_runner02) Restart() error {
	available, resData, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exist. Can't return status.")
	}

	_, err1 := exec.Command("supervisorctl", "restart", resData.GatewayProgram).Output()
	_, err2 := exec.Command("supervisorctl", "restart", resData.BridgeProgram).Output()

	if err1 == nil && err2 == nil {
		log.Info("Triggered restart")
	} else {
		log.Warning("Triggered restart, however supervisor did return some errors. ", err1.Error(), " ", err2.Error())
	}

	return nil
}

func (r *linux_amd64_supervisor_runner02) Recreate() error {
	available, resData, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exist. Can't return status.")
	}
	err = r.Destroy()
	if err != nil {
		return err
	}

	err = r.PostRun()
	if err != nil {
		return err
	}

	err = r.Prepare()
	if err != nil {
		return err
	}

	ref := reflect.ValueOf(resData)
	typeOfref := ref.Type()
	runtimeArgs := make(map[string]string)

	for i := 0; i < ref.NumField(); i++ {
		var name = typeOfref.Field(i).Name
		if name != "StartTime" {
			runtimeArgs[name] = ref.Field(i).String()
		}
	}

	err = r.Create(runtimeArgs)
	if err != nil {
		return err
	}
	return nil
}

func (r *linux_amd64_supervisor_runner02) Destroy() error {
	available, resData, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exists. Can't destroy")
	}

	returned, err := exec.Command("supervisorctl", "stop", resData.GatewayProgram).Output()
	if err != nil {
		alreadyDead, err2 := regexp.MatchString("not running", string(returned))
		if !alreadyDead || err2 != nil {
			return errors.New("Error while stopping gateway: " + err.Error())
		}
	}
	log.Debug("Trigerred gateway stop")

	returned, err = exec.Command("supervisorctl", "stop", resData.BridgeProgram).Output()
	if err != nil {
		alreadyDead, err2 := regexp.MatchString("not running", string(returned))
		if !alreadyDead || err2 != nil {
			return errors.New("Error while stopping bridge: " + err.Error())
		}
	}
	log.Debug("Trigerred bridge stop")

	log.Info("Waiting 5 seconds for SIGTERM to take effect")
	time.Sleep(5 * time.Second)

	return nil
}

func (r *linux_amd64_supervisor_runner02) PostRun() error {
	var gatewayConfig = runner02supervisorConfFiles + "/" + runner02gatewaySupervisorConfFile + "_" + r.InstanceId + ".conf"
	var bridgeConfig = runner02supervisorConfFiles + "/" + runner02bridgeSupervisorConfFile + "_" + r.InstanceId + ".conf"

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

	_, err := exec.Command("supervisorctl", "reread").Output()
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

	log.Info("All relevant processes stopped, resources deleted, supervisor configs removed")
	return nil
}

func (r *linux_amd64_supervisor_runner02) Status() error {
	available, resData, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exist. Can't return status.")
	}

	var projectConfig types.Project
	err = viper.UnmarshalKey(runner02projectName, &projectConfig)
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

	status, _ := exec.Command("supervisorctl", "status").Output()

	var supervisorStatus = make(map[string]interface{})

	statusLines := strings.Split(string(status), "\n")
	var anyStatusLine = false
	for _, v := range statusLines {
		if match, err := regexp.MatchString(runner02gatewayProgramName+"_"+r.InstanceId+"|"+runner02bridgeProgramName+"_"+r.InstanceId, v); err == nil && match {
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

func (r *linux_amd64_supervisor_runner02) Logs(lines int) error {
	available, _, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exists. Can't tail logs")
	}
	// Check for resource
	fileSubscriptions := make(map[string]string)
	var runner02logRootDir = "/var/log/supervisor/"
	err = filepath.Walk(runner02logRootDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			for _, v := range []string{runner02gatewaySupervisorConfFile + "_" + r.InstanceId + "-stdout.*",
				runner02gatewaySupervisorConfFile + "_" + r.InstanceId + "-stderr.*",
				runner02bridgeSupervisorConfFile + "_" + r.InstanceId + "-stdout.*",
				runner02bridgeSupervisorConfFile + "_" + r.InstanceId + "-stderr.*"} {
				r, err := regexp.MatchString(v, f.Name())
				if err == nil && r {
					fileSubscriptions[v[:len(v)-2]] = runner02logRootDir + f.Name()
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
			seeklocation := util.GetFileSeekOffsetLastNLines(filelocation, lines)
			t, err := tail.TailFile(filelocation, tail.Config{Location: &tail.SeekInfo{Offset: seeklocation}, Follow: true, Logger: tail.DiscardingLogger})
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

type runner02resource struct {
	Runner, Version, StartTime                                                                                                                   string
	GatewayProgram, GatewayUser, GatewayRunDir, GatewayExecutablePath, GatewayKeyfile, GatewayListenPortPeer, GatewayMarlinIp, GatewayMarlinPort string
	BridgeProgram, BridgeUser, BridgeRunDir, BridgeExecutablePath, BridgeBootstrapAddr                                                           string
}

func (r *linux_amd64_supervisor_runner02) fetchResourceInformation(fileLocation string) (bool, runner02resource, error) {
	if _, err := os.Stat(fileLocation); os.IsNotExist(err) {
		return false, runner02resource{}, err
	}

	file, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return false, runner02resource{}, err
	}

	var resData = runner02resource{}
	err = json.Unmarshal([]byte(file), &resData)

	return true, resData, err
}

func (r *linux_amd64_supervisor_runner02) writeResourceToFile(resData runner02resource, fileLocation string) error {
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
