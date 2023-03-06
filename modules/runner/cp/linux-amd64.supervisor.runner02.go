package cp

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
	Cp string
	CpChecksum string
}

type linux_amd64_supervisor_runner02 struct {
	Version      string
	Storage      string
	InstanceId   string
	RunnerData   linux_amd64_supervisor_runner02_runnerdata
	SkipChecksum bool
}

type runner02resource struct {
	Runner, Version, StartTime string
	CpProgram, CpUser, CpRunDir, CpPath, AwsProfile, KeyName, Rpc, Regions string
}

const (
	runner02cpName               = "control-plane"
	runner02cpProgramName        = "cp"
	runner02defaultUser              = "root"
	runner02supervisorConfFiles      = "/etc/supervisor/conf.d"
	runner02cpSupervisorConfFile = "cp"
	runner02logRootDir               = "/var/log/supervisor"
	runner02oldLogRootDir            = "/var/log/old_logs"
	runner02projectName              = "control-plane"
)

func (r *linux_amd64_supervisor_runner02) PreRunSanity() error {
	if !util.IsSupervisorAvailable() {
		return errors.New("system does not support supervisor")
	}
	if !util.IsSupervisorInRunningState() {
		return errors.New("system does not have supervisor in running state")
	}
	return nil
}

func (r *linux_amd64_supervisor_runner02) Download() error {
	var dirPath = r.Storage + "/" + r.Version
	err := util.CreateDirPathIfNotExists(dirPath)
	if err != nil {
		return err
	}

	var cpLocation = dirPath + "/" + runner02cpName

	if _, err := os.Stat(cpLocation); os.IsNotExist(err) {
		log.Info("Fetching control-plane from upstream for version ", r.Version)
		util.DownloadFile(cpLocation, r.RunnerData.Cp)
	}
	if !r.SkipChecksum {
		err := util.VerifyChecksum(cpLocation, r.RunnerData.CpChecksum)
		if err != nil {
			return errors.New("Error while verifying cp checksum: " + err.Error())
		} else {
			log.Debug("Successully verified cp's integrity")
		}
	}

	err = os.Chmod(cpLocation, 0755)
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
	currentUser, err := util.GetUser()
	if err != nil {
		return err
	}

	substitutions := runner02resource {
		"linux-amd64.supervisor.runner02", r.Version, time.Now().Format(time.RFC822Z),
		runner02cpProgramName + "_" + r.InstanceId, currentUser.Username, currentUser.HomeDir, r.Storage + "/" + r.Version + "/" + runner02cpName, "default", "marlin", "", "ap-south-1",
	}

	for k, v := range runtimeArgs {
		if k != "CpProgram" &&
			reflect.ValueOf(&substitutions).Elem().FieldByName(k).CanSet() {
			reflect.ValueOf(&substitutions).Elem().FieldByName(k).SetString(v)
		}
	}

	log.Info("Running configuration")
	util.PrettyPrintKVStruct(substitutions)

	gt := template.Must(template.New("cp-template").Parse(util.TrimSpacesEveryLine(`
		[program:{{.CpProgram}}]
		process_name={{.CpProgram}}
		user={{.CpUser}}
		directory={{.CpRunDir}}
		command={{.CpExecutablePath}} --profile {{.AwsProfile}} --key-name {{.KeyName}} --rpc {{.Rpc}} --regions {{.Regions}}
		priority=100
		numprocs=1
		numprocs_start=1
		autostart=true
		autorestart=true
		stdout_logfile=/var/log/supervisor/{{.CpProgram}}-stdout.log
		stderr_logfile=/var/log/supervisor/{{.CpProgram}}-stderr.log
	`)))
	gFile, err := os.Create(runner02supervisorConfFiles + "/" + runner02cpSupervisorConfFile + r.InstanceId + ".conf")
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

	_, err = exec.Command("supervisorctl", "start", substitutions.CpProgram).Output()
	if err != nil {
		return errors.New("Error while starting proxy: " + err.Error())
	}
	log.Debug("Trigerred proxy run")

	log.Info("Waiting 2 seconds to poll for status")
	time.Sleep(2 * time.Second)

	status, err := exec.Command("supervisorctl", "status").Output()
	if err != nil {
		log.Warning("Error while reading supervisor status: " + err.Error())
	}
	var supervisorStatus = make(map[string]interface{})

	statusLines := strings.Split(string(status), "\n")
	var anyStatusLine = false
	for _, v := range statusLines {
		if match, err := regexp.MatchString(substitutions.CpProgram, v); err == nil && match {
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

	_, err1 := exec.Command("supervisorctl", "restart", resData.CpProgram).Output()

	if err1 == nil {
		log.Info("Triggered restart")
	} else {
		log.Warning("Triggered restart, however supervisor did return some errors. ", err1.Error())
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

	_, err = exec.Command("supervisorctl", "stop", resData.CpProgram).Output()
	if err != nil {
		return errors.New("Error while stopping cp: " + err.Error())
	}
	log.Debug("Trigerred beacon stop")

	log.Info("Waiting 5 seconds for SIGTERM to take effect")
	time.Sleep(5 * time.Second)

	return nil
}

func (r *linux_amd64_supervisor_runner02) PostRun() error {
	var proxyConfig = runner02supervisorConfFiles + "/" + runner02cpSupervisorConfFile + r.InstanceId + ".conf"

	if _, err := os.Stat(proxyConfig); !os.IsNotExist(err) {
		if err := os.Remove(proxyConfig); err != nil {
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

	log.Info("All relevant processes stopped, resources deleted, supervisor configs removed, logs marked as old")
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
		if match, err := regexp.MatchString(runner02cpProgramName+"_"+r.InstanceId, v); err == nil && match {
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
	available, resData, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exists. Can't tail logs")
	}
	// Check for resource
	fileSubscriptions := make(map[string]string)
	err = filepath.Walk(runner02logRootDir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			for _, v := range []string{resData.CpProgram + "-stdout.*",
				resData.CpProgram + "-stderr.*"} {
				r, err := regexp.MatchString(v, f.Name())
				if err == nil && r {
					fileSubscriptions[v[:len(v)-2]] = runner02logRootDir + "/" + f.Name()
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

