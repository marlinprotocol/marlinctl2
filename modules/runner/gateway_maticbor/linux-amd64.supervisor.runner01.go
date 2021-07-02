package gateway_maticbor

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
}

type linux_amd64_supervisor_runner01 struct {
	Version      string
	Storage      string
	InstanceId   string
	RunnerData   linux_amd64_supervisor_runner01_runnerdata
	SkipChecksum bool
}

const (
	runner01gatewayName               = "gateway_maticbor_linux-amd64"
	runner01gatewayProgramName        = "gateway_maticbor"
	runner01defaultUser               = "root"
	runner01supervisorConfFiles       = "/etc/supervisor/conf.d"
	runner01gatewaySupervisorConfFile = "gateway_maticbor"
	runner01logRootDir                = "/var/log/supervisor"
	runner01oldLogRootDir             = "/var/log/old_logs"
	runner01projectName               = "gateway_maticbor"
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

	var gatewayLocation = dirPath + "/" + runner01gatewayName

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

	err = os.Chmod(gatewayLocation, 0755)
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

	err = util.ChownRmarlinctlDir()
	if err != nil {
		return err
	}

	return nil
}

func (r *linux_amd64_supervisor_runner01) Create(runtimeArgs map[string]string) error {
	if _, err := os.Stat(GetResourceFileLocation(r.Storage, r.InstanceId)); err == nil {
		return errors.New("Resource file already exisits, cannot create a new instance: " + GetResourceFileLocation(r.Storage, r.InstanceId))
	}

	currentUser, err := util.GetUser()
	if err != nil {
		return err
	}

	substitutions := runner01resource{
		"linux-amd64.supervisor.runner01", r.Version, time.Now().Format(time.RFC822Z),
		runner01gatewayProgramName + "_" + r.InstanceId, currentUser.Username, currentUser.HomeDir, r.Storage + "/" + r.Version + "/" + runner01gatewayName,
		"", "", "", "", "",
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
		command={{.GatewayExecutablePath}} --discovery-addr {{.DiscoveryAddr}} --pubsub-addr {{.PubsubAddr}} {{if .BootstrapAddr}} --beacon-addr {{.BootstrapAddr}}{{end}} {{if .KeystorePath}} --keystore-path {{.KeystorePath}}{{end}} {{if .KeystorePassPath}} --keystore-pass-path {{.KeystorePassPath}} {{end}}
		priority=100
		numprocs=1
		numprocs_start=1
		autostart=true
		autorestart=true
		stdout_logfile=/var/log/supervisor/{{.GatewayProgram}}-stdout.log
		stderr_logfile=/var/log/supervisor/{{.GatewayProgram}}-stderr.log
	`)))
	gFile, err := os.Create(runner01supervisorConfFiles + "/" + runner01gatewaySupervisorConfFile + "_" + r.InstanceId + ".conf")
	if err != nil {
		return err
	}
	defer gFile.Close()
	if err := gt.Execute(gFile, substitutions); err != nil {
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
			if match, err := regexp.MatchString(runner01gatewayProgramName+"_"+r.InstanceId, v); err == nil && match {
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

func (r *linux_amd64_supervisor_runner01) Restart() error {
	available, resData, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exist. Can't return status.")
	}

	_, err1 := exec.Command("supervisorctl", "restart", resData.GatewayProgram).Output()

	if err1 == nil {
		log.Info("Triggered restart")
	} else {
		log.Warning("Triggered restart, however supervisor did return some errors. ", err1.Error())
	}

	return nil
}

func (r *linux_amd64_supervisor_runner01) Recreate() error {
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

func (r *linux_amd64_supervisor_runner01) Destroy() error {
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

	log.Info("Waiting 5 seconds for SIGTERM to take effect")
	time.Sleep(5 * time.Second)

	return nil
}

func (r *linux_amd64_supervisor_runner01) PostRun() error {
	var gatewayConfig = runner01supervisorConfFiles + "/" + runner01gatewaySupervisorConfFile + "_" + r.InstanceId + ".conf"

	if _, err := os.Stat(gatewayConfig); !os.IsNotExist(err) {
		if err := os.Remove(gatewayConfig); err != nil {
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

func (r *linux_amd64_supervisor_runner01) Status() error {
	available, resData, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exist. Can't return status.")
	}

	var projectConfig types.Project
	err = viper.UnmarshalKey(runner01projectName, &projectConfig)
	if err != nil {
		return err
	}
	log.Info("Project configuration")
	util.PrettyPrintKVStruct(projectConfig)

	log.Info("Resource information")
	util.PrettyPrintKVStruct(resData)

	status, _ := exec.Command("supervisorctl", "status").Output()

	var supervisorStatus = make(map[string]interface{})

	statusLines := strings.Split(string(status), "\n")
	var anyStatusLine = false
	for _, v := range statusLines {
		if match, err := regexp.MatchString(runner01gatewayProgramName+"_"+r.InstanceId, v); err == nil && match {
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

func (r *linux_amd64_supervisor_runner01) Logs(lines int) error {
	available, _, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exists. Can't tail logs")
	}
	// Check for resource
	fileSubscriptions := make(map[string]string)
	var runner01logRootDir = "/var/log/supervisor/"
	err = filepath.Walk(runner01logRootDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			for _, v := range []string{runner01gatewaySupervisorConfFile + "_" + r.InstanceId + "-stdout.*",
				runner01gatewaySupervisorConfFile + "_" + r.InstanceId + "-stderr.*"} {
				r, err := regexp.MatchString(v, f.Name())
				if err == nil && r {
					fileSubscriptions[v[:len(v)-2]] = runner01logRootDir + f.Name()
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

type runner01resource struct {
	Runner, Version, StartTime                                               string
	GatewayProgram, GatewayUser, GatewayRunDir, GatewayExecutablePath        string
	DiscoveryAddr, PubsubAddr, BootstrapAddr, KeystorePath, KeystorePassPath string
}

func (r *linux_amd64_supervisor_runner01) fetchResourceInformation(fileLocation string) (bool, runner01resource, error) {
	if _, err := os.Stat(fileLocation); os.IsNotExist(err) {
		return false, runner01resource{}, err
	}

	file, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return false, runner01resource{}, err
	}

	var resData = runner01resource{}
	err = json.Unmarshal([]byte(file), &resData)

	return true, resData, err
}

func (r *linux_amd64_supervisor_runner01) writeResourceToFile(resData runner01resource, fileLocation string) error {
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
