package cp

import (
	"errors"

	"github.com/marlinprotocol/ctl2/modules/runner"
)

func GetRunnerInstance(runnerId string, version string, storage string, runnerData interface{}, skipRunnerData bool, skipChecksum bool, instanceId string) (runner.Runner, error) {
	switch runnerId {
	case "linux-amd64.supervisor.runner02" :
		if skipRunnerData {
			return &linux_amd64_supervisor_runner02{
				Version:      version,
				Storage:      storage,
				SkipChecksum: skipChecksum,
				InstanceId:   instanceId,
			}, nil
		}
		runnerDataMap := runnerData.(map[string]interface{})

		cp, ok1 := runnerDataMap["cp"]
		cpChecksum, ok2 := runnerDataMap["cp_checksum"]
		
		if !ok1 || !ok2 {
			return &linux_amd64_supervisor_runner02{}, errors.New("Incomplete / wrong runner data for version: " + version)
		}

		return &linux_amd64_supervisor_runner02{
			Version: version,
			Storage: storage,
			RunnerData: linux_amd64_supervisor_runner02_runnerdata {
				Cp:         cp.(string),
				CpChecksum: cpChecksum.(string),
			},
			SkipChecksum: skipChecksum,
			InstanceId:   instanceId,
		}, nil	
	default:
		return &linux_amd64_supervisor_runner02{}, errors.New("Unknown runnerId: " + runnerId)
	}
}

func GetResourceFileLocation(storage string, instanceId string) string {
	return storage + "/common/cp" + instanceId + ".resource"
}