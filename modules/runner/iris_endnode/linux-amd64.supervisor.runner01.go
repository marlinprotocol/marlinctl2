package iris_endnode

import (
	log "github.com/sirupsen/logrus"
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
	log.Info("Prerun sanity for runner01")
	return nil
}

func (r *linux_amd64_supervisor_runner01) Download() error {
	return nil
}

func (r *linux_amd64_supervisor_runner01) Prepare() error {
	log.Info("Prepare for runner01")
	return nil
}

func (r *linux_amd64_supervisor_runner01) Create() error {
	log.Info("Create for runner01")
	return nil
}

func (r *linux_amd64_supervisor_runner01) Destroy() error {
	return nil
}

func (r *linux_amd64_supervisor_runner01) PostRun() error {
	return nil
}

func (r *linux_amd64_supervisor_runner01) Status() error {
	return nil
}

func (r *linux_amd64_supervisor_runner01) Logs(tailLogs bool, prevLines int) error {
	return nil
}
