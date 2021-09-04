package gateway_polygonbor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"
	"time"

	"github.com/marlinprotocol/ctl2/modules/util"
	"github.com/marlinprotocol/ctl2/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type linux_amd64_supervisor_runner02_runnerdata struct {
	Gateway          string
	GatewayChecksum  string
	MevProxy         string
	MevProxyChecksum string
}

type linux_amd64_supervisor_runner02 struct {
	Version      string
	Storage      string
	InstanceId   string
	RunnerData   linux_amd64_supervisor_runner02_runnerdata
	SkipChecksum bool
}

const (
	runner02gatewayName                = "gateway_polygonbor_linux-amd64"
	runner02gatewayProgramName         = "gateway_polygonbor"
	runner02mevproxyName               = "mevproxy_linux-amd64"
	runner02mevproxyProgramName        = "mevproxy_polygon"
	runner02supervisorConfFiles        = "/etc/supervisor/conf.d"
	runner02gatewaySupervisorConfFile  = "gateway_polygonbor"
	runner02mevproxySupervisorConfFile = "mevproxy_polygonbor"
	runner02projectName                = "gateway_polygonbor"
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

	err = util.DownloadExecutable("gateway", r.Version, r.RunnerData.Gateway, r.SkipChecksum, r.RunnerData.GatewayChecksum, dirPath+"/"+runner02gatewayName)
	if err != nil {
		return err
	}

	err = util.DownloadExecutable("mevproxy", r.Version, r.RunnerData.MevProxy, r.SkipChecksum, r.RunnerData.MevProxyChecksum, dirPath+"/"+runner02mevproxyName)
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

	substitutions := runner02resource{
		"linux-amd64.supervisor.runner02", r.Version, time.Now().Format(time.RFC822Z),
		runner02gatewayProgramName + "_" + r.InstanceId, currentUser.Username, currentUser.HomeDir, r.Storage + "/" + r.Version + "/" + runner02gatewayName,
		"", "", "", "", "", "", "",
		runner02mevproxyProgramName + "_" + r.InstanceId, currentUser.Username, currentUser.HomeDir, r.Storage + "/" + r.Version + "/" + runner02mevproxyName,
		"", "",
	}

	for k, v := range runtimeArgs {
		if k != "GatewayProgram" &&
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
		command={{.GatewayExecutablePath}} --discovery-addr {{.DiscoveryAddr}} --pubsub-addr {{.PubsubAddr}} {{if .BootstrapAddr}} --beacon-addr {{.BootstrapAddr}}{{end}} {{if .SpamcheckAddr}} --spamcheck-addr {{.SpamcheckAddr}}{{end}} {{if .KeystorePath}} --keystore-path {{.KeystorePath}}{{end}} {{if .KeystorePassPath}} --keystore-pass-path {{.KeystorePassPath}} {{end}} --contracts {{.Contracts}}
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
	if err := gt.Execute(gFile, substitutions); err != nil {
		panic(err)
	}
	gFile.Close()

	mpt := template.Must(template.New("mevproxy-template").Parse(util.TrimSpacesEveryLine(`
		[program:{{.MevProxyProgram}}]
		process_name={{.MevProxyProgram}}
		user={{.MevProxyUser}}
		directory={{.MevProxyRunDir}}
		command={{.MevProxyExecutablePath}} -listenAddr {{.MevProxyListenAddr}} -rpcAddr {{.MevProxyBundleAddr}}
		priority=100
		numprocs=1
		numprocs_start=1
		autostart=true
		autorestart=true
		stdout_logfile=/var/log/supervisor/{{.MevProxyProgram}}-stdout.log
		stderr_logfile=/var/log/supervisor/{{.MevProxyProgram}}-stderr.log
	`)))
	mpFile, err := os.Create(runner02supervisorConfFiles + "/" + runner02mevproxySupervisorConfFile + "_" + r.InstanceId + ".conf")
	if err != nil {
		return err
	}
	if err := mpt.Execute(mpFile, substitutions); err != nil {
		panic(err)
	}
	mpFile.Close()

	err = util.SupervisorStart([]string{substitutions.GatewayProgram, substitutions.MevProxyProgram})
	if err != nil {
		return err
	}
	util.SupervisorStatusBestEffort([]string{runner02gatewayProgramName, runner02mevproxyProgramName}, r.InstanceId)
	return r.writeResourceToFile(substitutions, GetResourceFileLocation(r.Storage, r.InstanceId))
}

func (r *linux_amd64_supervisor_runner02) Restart() error {
	available, resData, err := r.fetchResourceInformation(GetResourceFileLocation(r.Storage, r.InstanceId))
	if err != nil {
		return err
	}
	if !available {
		return errors.New("resource by id " + r.InstanceId + " doesn't exist. Can't return status.")
	}

	util.SupervisorRestartProgramBestEffort("gateway", resData.GatewayProgram)
	util.SupervisorRestartProgramBestEffort("mevpoxy", resData.MevProxyProgram)

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

	errs := util.SupervisorStop([]string{resData.MevProxyProgram, resData.GatewayProgram})
	if len(errs) != 0 {
		return errors.New(fmt.Sprintf("Error while stopping programs %v", errs))
	}
	return nil
}

func (r *linux_amd64_supervisor_runner02) PostRun() error {
	var gatewayConfig = runner02supervisorConfFiles + "/" + runner02gatewaySupervisorConfFile + "_" + r.InstanceId + ".conf"

	if _, err := os.Stat(gatewayConfig); !os.IsNotExist(err) {
		if err := os.Remove(gatewayConfig); err != nil {
			return err
		}
	}

	err := util.SupervisorRereadUpdate()
	if err != nil {
		return err
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

	log.Info("Resource information")
	util.PrettyPrintKVStruct(resData)

	log.Info("Process Status")
	util.SupervisorStatusBestEffort([]string{resData.GatewayProgram, resData.MevProxyProgram}, r.InstanceId)

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

	return util.LogTailer([]string{runner02gatewaySupervisorConfFile, runner02mevproxySupervisorConfFile},
		r.InstanceId, lines)
}

type runner02resource struct {
	Runner, Version, StartTime                                                                         string
	GatewayProgram, GatewayUser, GatewayRunDir, GatewayExecutablePath                                  string
	DiscoveryAddr, PubsubAddr, BootstrapAddr, KeystorePath, KeystorePassPath, SpamcheckAddr, Contracts string
	MevProxyProgram, MevProxyUser, MevProxyRunDir, MevProxyExecutablePath                              string
	MevProxyListenAddr, MevProxyBundleAddr                                                             string
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
