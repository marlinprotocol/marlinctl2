package gateway_iris

import (
	"errors"

	"github.com/marlinprotocol/ctl2/modules/runner"
)

func GetRunnerInstance(runnerId string, version string, storage string, runnerData interface{}, skipRunnerData bool, skipChecksum bool, instanceId string) (runner.Runner, error) {
	switch runnerId {
	case "linux-amd64.supervisor.runner01":
		if skipRunnerData {
			return &linux_amd64_supervisor_runner01{
				Version:      version,
				Storage:      storage,
				SkipChecksum: skipChecksum,
				InstanceId:   instanceId,
			}, nil
		}
		runnerDataMap := runnerData.(map[string]interface{})

		gateway, ok1 := runnerDataMap["gateway"]
		gatewayChecksum, ok2 := runnerDataMap["gateway_checksum"]
		bridge, ok3 := runnerDataMap["bridge"]
		bridgeChecksum, ok4 := runnerDataMap["bridge_checksum"]

		if !ok1 || !ok2 || !ok3 || !ok4 {
			return &linux_amd64_supervisor_runner01{}, errors.New("Incomplete / wrong runner data for version: " + version)
		}

		return &linux_amd64_supervisor_runner01{
			Version: version,
			Storage: storage,
			RunnerData: linux_amd64_supervisor_runner01_runnerdata{
				Gateway:         gateway.(string),
				GatewayChecksum: gatewayChecksum.(string),
				Bridge:          bridge.(string),
				BridgeChecksum:  bridgeChecksum.(string),
			},
			SkipChecksum: skipChecksum,
			InstanceId:   instanceId,
		}, nil
	default:
		return &linux_amd64_supervisor_runner01{}, errors.New("Unknown runnerId: " + runnerId)
	}
}

func GetResourceFileLocation(storage string, instanceId string) string {
	return storage + "/common/project_gateway_iris_instance" + instanceId + ".resource"
}
